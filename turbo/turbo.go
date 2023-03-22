package turbo

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type StreamAction string

// TODO move timeouts to config
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512

	AppendAction  StreamAction = "append"
	PrependAction StreamAction = "prepend"
	ReplaceAction StreamAction = "replace"
	UpdateAction  StreamAction = "update"
	RemoveAction  StreamAction = "remove"
	BeforeAction  StreamAction = "before"
	AfterAction   StreamAction = "after"
)

const ContentType = "text/vnd.turbo-stream.html"

var bufPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

var ErrInvalid = errors.New("invalid stream names")

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
	hub     *Hub[T]
	conn    *websocket.Conn
	streams []string
	msgs    chan []byte
	Data    T
}

type Renderer interface {
	Render(context.Context, io.Writer) error
}

type StreamMessage struct {
	Action         StreamAction
	Target         string
	TargetSelector string
	Template       string
	Renderer       Renderer
}

func (s StreamMessage) Render(r Renderer) StreamMessage {
	s.Renderer = r
	return s
}

func Append(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   AppendAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func AppendMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         AppendAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Prepend(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   PrependAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func PrependMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         PrependAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Replace(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   ReplaceAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func ReplaceMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         ReplaceAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Update(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   UpdateAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func UpdateMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         UpdateAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Remove(target string) StreamMessage {
	return StreamMessage{
		Action: RemoveAction,
		Target: target,
	}
}

func RemoveMatch(target string) StreamMessage {
	return StreamMessage{
		Action:         RemoveAction,
		TargetSelector: target,
	}
}

func Before(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   BeforeAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func BeforeMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         BeforeAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func After(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   AfterAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func AfterMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         AfterAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

// TODO context, error handling
func (c *Client[T]) Send(streams ...StreamMessage) {
	if len(streams) == 0 {
		return
	}
	msgs, err := serializeStreamMessages(context.TODO(), streams)
	if err != nil {
		return
	}
	c.msgs <- msgs
}

// func (c *Client[T]) Write(b []byte) (int, error) {
// 	c.msgs <- b
// 	return len(b), nil
// }

type Responder[T any] interface {
	Respond(*Client[T], []byte)
}

type Hub[T any] struct {
	config    Config[T]
	upgrader  websocket.Upgrader
	clients   map[*Client[T]]struct{}
	streams   map[string]map[*Client[T]]struct{}
	clientsMu sync.RWMutex
	streamsMu sync.RWMutex
}

type Config[T any] struct {
	// Secret should be a random 256 bit key
	Secret    []byte
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
		streams: make(map[string]map[*Client[T]]struct{}),
	}
}

// see https://github.com/gtank/cryptopasta/blob/master/encrypt.go
// and https://www.alexedwards.net/blog/working-with-cookies-in-go#encrypted-cookies
func (h *Hub[T]) EncryptStreamNames(names []string) (string, error) {
	msg := strings.Join(names, ",")

	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(h.config.Secret)
	if err != nil {
		return "", err
	}

	// Wrap the cipher block in Galois Counter Mode.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a unique nonce containing 12 random bytes.
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	// Encrypt the data using aesGCM.Seal(). By passing the nonce as the first
	// parameter, the encrypted message will be appended to the nonce so
	// that the encrypted message will be in the format
	// "{nonce}{encrypted message}".
	cryptedMsg := gcm.Seal(nonce, nonce, []byte(msg), nil)

	return base64.URLEncoding.EncodeToString(cryptedMsg), nil
}

func (h *Hub[T]) DecryptStreamNames(names string) ([]string, error) {
	cryptedMsg, err := base64.URLEncoding.DecodeString(names)
	if err != nil {
		return nil, err
	}

	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(h.config.Secret)
	if err != nil {
		return nil, err
	}

	// Wrap the cipher block in Galois Counter Mode.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	// Avoid potential 'index out of range' panic in the next step.
	if len(cryptedMsg) < nonceSize {
		return nil, ErrInvalid
	}

	// Split cryptedMsg in nonce and encrypted message and use gcm.Open() to
	// decrypt and authenticate the data.
	msg, err := gcm.Open(nil, cryptedMsg[:nonceSize], cryptedMsg[nonceSize:], nil)
	if err != nil {
		return nil, errors.Join(err, ErrInvalid)
	}

	return strings.Split(string(msg), ","), nil
}

// TODO context, error handling
func (h *Hub[T]) Broadcast(streams ...StreamMessage) {
	if len(streams) == 0 {
		return
	}
	msg, err := serializeStreamMessages(context.TODO(), streams)
	if err != nil {
		return
	}
	for c := range h.clients {
		c.msgs <- msg
	}
}

// TODO context, error handling
func (h *Hub[T]) Send(k string, msgs ...StreamMessage) {
	if len(msgs) == 0 {
		return
	}

	msg, err := serializeStreamMessages(context.TODO(), msgs)
	if err != nil {
		log.Print(err)
		return
	}

	h.streamsMu.RLock()
	defer h.streamsMu.RUnlock()

	log.Printf("clients: %+v", h.streams)

	if clients, ok := h.streams[k]; ok {
		for c := range clients {
			log.Printf("send msg: %+v %+v", c, msg)
			c.msgs <- msg
		}
	}
}

func (h *Hub[T]) disconnect(c *Client[T]) {
	h.streamsMu.Lock()
	for _, k := range c.streams {
		h.removeClientFromStream(k, c)
	}
	h.streamsMu.Unlock()
	h.clientsMu.Lock()
	delete(h.clients, c)
	h.clientsMu.Unlock()
}

// TODO error handler
func (h *Hub[T]) Handle(w http.ResponseWriter, r *http.Request, cryptedStreams string) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	streams, err := h.DecryptStreamNames(cryptedStreams)
	if err != nil {
		return err
	}

	c := &Client[T]{
		hub:     h,
		conn:    conn,
		streams: streams,
		msgs:    make(chan []byte, 64),
	}

	h.clientsMu.Lock()
	h.clients[c] = struct{}{}
	h.clientsMu.Unlock()

	for _, stream := range streams {
		h.addClientToStream(stream, c)
	}

	log.Printf("streams: %+v", h.streams)

	go writePump(h, c)
	go readPump(h, c)

	return nil
}

func (h *Hub[T]) addClientToStream(k string, c *Client[T]) {
	h.streamsMu.Lock()
	if clients, ok := h.streams[k]; ok {
		clients[c] = struct{}{}
	} else {
		h.streams[k] = map[*Client[T]]struct{}{c: {}}
	}
	h.streamsMu.Unlock()
}

func (h *Hub[T]) removeClientFromStream(stream string, c *Client[T]) {
	h.streamsMu.Lock()
	if clients, ok := h.streams[stream]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.streams, stream)
		}
	}
	h.streamsMu.Unlock()
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

func serializeStreamMessages(ctx context.Context, streams []StreamMessage) ([]byte, error) {
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

func Render(w http.ResponseWriter, r *http.Request, code int, streams ...StreamMessage) error {
	if hdr := w.Header(); hdr.Get("Content-Type") == "" {
		hdr.Set("Content-Type", ContentType)
	}
	w.WriteHeader(code)
	b, err := serializeStreamMessages(r.Context(), streams)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
