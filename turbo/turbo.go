package turbo

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Action string

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512

	AppendAction  Action = "append"
	PrependAction Action = "prepend"
	ReplaceAction Action = "replace"
	UpdateAction  Action = "update"
	RemoveAction  Action = "remove"
	BeforeAction  Action = "before"
	AfterAction   Action = "after"
)

const ContentType = "text/vnd.turbo-stream.html"

func Request(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Accept"), ContentType)
}

func FrameRequestID(r *http.Request) string {
	return r.Header.Get("Turbo-Frame")
}

func FrameRequest(r *http.Request) bool {
	return FrameRequestID(r) != ""
}

// TODO give each client a unique id
// add exclude to Send to avoid jitter etc
type Client[T any] struct {
	hub       *Hub[T]
	conn      *websocket.Conn
	indexKeys []string
	msgs      chan []byte
	Data      T
}

type Renderer interface {
	Render(context.Context, io.Writer) error
}

type Stream struct {
	Action         Action
	Target         string
	TargetSelector string
	Template       string
	Renderer       Renderer
}

func (s Stream) Render(r Renderer) Stream {
	s.Renderer = r
	return s
}

func Append(target string, tmpls ...string) Stream {
	return Stream{
		Action:   AppendAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func AppendMatch(target string, tmpls ...string) Stream {
	return Stream{
		Action:         AppendAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Prepend(target string, tmpls ...string) Stream {
	return Stream{
		Action:   PrependAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func PrependMatch(target string, tmpls ...string) Stream {
	return Stream{
		Action:         PrependAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Replace(target string, tmpls ...string) Stream {
	return Stream{
		Action:   ReplaceAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func ReplaceMatch(target string, tmpls ...string) Stream {
	return Stream{
		Action:         ReplaceAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Update(target string, tmpls ...string) Stream {
	return Stream{
		Action:   UpdateAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func UpdateMatch(target string, tmpls ...string) Stream {
	return Stream{
		Action:         UpdateAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Remove(target string) Stream {
	return Stream{
		Action: RemoveAction,
		Target: target,
	}
}

func RemoveMatch(target string) Stream {
	return Stream{
		Action:         RemoveAction,
		TargetSelector: target,
	}
}

func Before(target string, tmpls ...string) Stream {
	return Stream{
		Action:   BeforeAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func BeforeMatch(target string, tmpls ...string) Stream {
	return Stream{
		Action:         BeforeAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func After(target string, tmpls ...string) Stream {
	return Stream{
		Action:   AfterAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func AfterMatch(target string, tmpls ...string) Stream {
	return Stream{
		Action:         AfterAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

// TODO context, error handling
func (c *Client[T]) Send(streams ...Stream) {
	if len(streams) == 0 {
		return
	}
	msgs, err := serializeStreams(context.TODO(), streams)
	if err != nil {
		return
	}
	c.msgs <- msgs
}

// func (c *Client[T]) Write(b []byte) (int, error) {
// 	c.msgs <- b
// 	return len(b), nil
// }

func (c *Client[T]) Join(keys ...string) {
	for _, k := range keys {
		c.hub.addClientToIndex(k, c)
		knownKey := false
		for _, ik := range c.indexKeys {
			if ik == k {
				knownKey = true
				break
			}
		}
		if !knownKey {
			// TODO mutex or channel
			c.indexKeys = append(c.indexKeys, k)
		}
	}
}

func (c *Client[T]) Leave(keys ...string) {
	for _, k := range keys {
		c.hub.removeClientFromIndex(k, c)
		// TODO mutex or channel
		for i, ik := range c.indexKeys {
			if ik == k {
				c.indexKeys = append(c.indexKeys[:i], c.indexKeys[i+1:]...)
				break
			}
		}
	}
}

func (c *Client[T]) LeaveAll() {
	// TODO mutex or channel
	for _, k := range c.indexKeys {
		c.hub.removeClientFromIndex(k, c)
	}
	c.indexKeys = nil
}

type Responder[T any] interface {
	Respond(*Client[T], []byte)
}

type Hub[T any] struct {
	config    Config[T]
	upgrader  websocket.Upgrader
	clients   map[*Client[T]]struct{}
	index     map[string]map[*Client[T]]struct{}
	clientsMu sync.RWMutex
	indexMu   sync.RWMutex
}

type Config[T any] struct {
	Responder Responder[T]
}

func NewHub[T any](config Config[T]) *Hub[T] {
	return &Hub[T]{
		config: config,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients: make(map[*Client[T]]struct{}),
		index:   make(map[string]map[*Client[T]]struct{}),
	}
}

// TODO context, error handling
func (h *Hub[T]) Broadcast(streams ...Stream) {
	if len(streams) == 0 {
		return
	}
	msg, err := serializeStreams(context.TODO(), streams)
	if err != nil {
		return
	}
	for c := range h.clients {
		c.msgs <- msg
	}
}

// TODO context, error handling
func (h *Hub[T]) Send(k string, streams ...Stream) {
	if len(streams) == 0 {
		return
	}

	msg, err := serializeStreams(context.TODO(), streams)
	if err != nil {
		return
	}

	h.indexMu.RLock()
	defer h.indexMu.RUnlock()

	if clients, ok := h.index[k]; ok {
		for c := range clients {
			c.msgs <- msg
		}
	}
}

func (h *Hub[T]) disconnect(c *Client[T]) {
	c.LeaveAll()
	h.clientsMu.Lock()
	delete(h.clients, c)
	h.clientsMu.Unlock()
}

func (h *Hub[T]) Handle(w http.ResponseWriter, r *http.Request, visitors ...func(*Client[T])) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &Client[T]{
		hub:  h,
		conn: conn,
		msgs: make(chan []byte, 64),
	}

	for _, fn := range visitors {
		fn(c)
	}

	h.clientsMu.Lock()
	h.clients[c] = struct{}{}
	h.clientsMu.Unlock()

	go writePump(h, c)
	go readPump(h, c)
}

func (h *Hub[T]) addClientToIndex(k string, c *Client[T]) {
	h.indexMu.Lock()
	if clients, ok := h.index[k]; ok {
		clients[c] = struct{}{}
	} else {
		h.index[k] = map[*Client[T]]struct{}{c: {}}
	}
	h.indexMu.Unlock()
}

func (h *Hub[T]) removeClientFromIndex(k string, c *Client[T]) {
	h.indexMu.Lock()
	if clients, ok := h.index[k]; ok {
		delete(clients, c)
	}
	h.indexMu.Unlock()
}

func writePump[T any](h *Hub[T], c *Client[T]) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msg, ok := <-c.msgs:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print(err)
				return
			}
			if _, err := w.Write(msg); err != nil {
				log.Print(err)
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func readPump[T any](h *Hub[T], c *Client[T]) {
	defer func() {
		h.disconnect(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("message: %s", msg)
		h.config.Responder.Respond(c, msg)
	}
}

func serializeStreams(ctx context.Context, streams []Stream) ([]byte, error) {
	b := bufPool.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		bufPool.Put(b)
	}()

	for _, s := range streams {
		b.WriteString(`<turbo-stream action="`)
		b.WriteString(string(s.Action))
		b.WriteString(`" `)
		if s.Target != "" {
			b.WriteString(`target="`)
			b.WriteString(s.Target)
		} else {
			b.WriteString(`targets="`)
			b.WriteString(s.TargetSelector)
		}
		b.WriteString(`">`)
		if s.Action != RemoveAction {
			b.WriteString(`<template>`)
			if s.Renderer != nil {
				s.Renderer.Render(ctx, b)
			} else {
				b.WriteString(s.Template)
			}
			b.WriteString(`</template>`)
		}
		b.WriteString(`</turbo-stream>`)
	}
	return b.Bytes(), nil
}

func Render(w http.ResponseWriter, r *http.Request, code int, streams ...Stream) error {
	if hdr := w.Header(); hdr.Get("Content-Type") == "" {
		hdr.Set("Content-Type", ContentType)
	}
	w.WriteHeader(code)
	b, err := serializeStreams(r.Context(), streams)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

var bufPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}
