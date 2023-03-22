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

const (
	AppendAction  StreamAction = "append"
	PrependAction StreamAction = "prepend"
	ReplaceAction StreamAction = "replace"
	UpdateAction  StreamAction = "update"
	RemoveAction  StreamAction = "remove"
	BeforeAction  StreamAction = "before"
	AfterAction   StreamAction = "after"

	ContentType = "text/vnd.turbo-stream.html"
)

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

type Client struct {
	hub     *Hub
	ws      *websocket.Conn
	streams []string
	msgs    chan []byte
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
func (c *Client) Send(streams ...StreamMessage) {
	if len(streams) == 0 {
		return
	}
	msgs, err := serializeStreamMessages(context.TODO(), streams)
	if err != nil {
		return
	}
	c.msgs <- msgs
}

type Hub struct {
	config    Config
	upgrader  websocket.Upgrader
	clients   map[*Client]struct{}
	streams   map[string]map[*Client]struct{}
	clientsMu sync.RWMutex
	streamsMu sync.RWMutex
}

type Config struct {
	// Secret should be a random 256 bit key
	Secret []byte
	// Time allowed to write a message to the peer.
	WriteWait time.Duration
	// Time allowed to read the next pong message from the peer.
	PongWait time.Duration
	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod time.Duration
	// Maximum message size allowed from peer.
	MaxMessageSize int64
}

func NewHub(config Config) *Hub {
	if config.WriteWait == 0 {
		config.WriteWait = 10 * time.Second
	}
	if config.PongWait == 0 {
		config.PongWait = 60 * time.Second
	}
	if config.PingPeriod == 0 {
		config.PingPeriod = (config.PongWait * 9) / 10
	}
	if config.MaxMessageSize == 0 {
		config.MaxMessageSize = 512
	}

	return &Hub{
		config: config,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients: make(map[*Client]struct{}),
		streams: make(map[string]map[*Client]struct{}),
	}
}

// see https://github.com/gtank/cryptopasta/blob/master/encrypt.go
// and https://www.alexedwards.net/blog/working-with-cookies-in-go#encrypted-cookies
func (h *Hub) EncryptStreamNames(names []string) (string, error) {
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

func (h *Hub) DecryptStreamNames(names string) ([]string, error) {
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
func (h *Hub) Broadcast(streams ...StreamMessage) {
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
func (h *Hub) Send(k string, msgs ...StreamMessage) {
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

func (h *Hub) disconnect(c *Client) {
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
func (h *Hub) Handle(w http.ResponseWriter, r *http.Request, cryptedStreams string) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	streams, err := h.DecryptStreamNames(cryptedStreams)
	if err != nil {
		return err
	}

	c := &Client{
		hub:     h,
		ws:      conn,
		streams: streams,
		msgs:    make(chan []byte, 64),
	}

	h.clientsMu.Lock()
	h.clients[c] = struct{}{}
	h.clientsMu.Unlock()

	for _, stream := range streams {
		h.addClientToStream(stream, c)
	}

	go wsWrite(h, c)
	go wsRead(h, c)

	return nil
}

func (h *Hub) addClientToStream(k string, c *Client) {
	h.streamsMu.Lock()
	if clients, ok := h.streams[k]; ok {
		clients[c] = struct{}{}
	} else {
		h.streams[k] = map[*Client]struct{}{c: {}}
	}
	h.streamsMu.Unlock()
}

func (h *Hub) removeClientFromStream(stream string, c *Client) {
	h.streamsMu.Lock()
	if clients, ok := h.streams[stream]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.streams, stream)
		}
	}
	h.streamsMu.Unlock()
}

// TODO logging
func wsWrite(h *Hub, c *Client) {
	pingTicker := time.NewTicker(h.config.PingPeriod)

	defer func() {
		pingTicker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case <-pingTicker.C:
			c.ws.SetWriteDeadline(time.Now().Add(h.config.WriteWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msg, ok := <-c.msgs:
			c.ws.SetWriteDeadline(time.Now().Add(h.config.WriteWait))

			if !ok {
				// The hub closed the channel.
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func wsRead(h *Hub, c *Client) {
	defer func() {
		h.disconnect(c)
		c.ws.Close()
	}()

	c.ws.SetReadLimit(h.config.MaxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(h.config.PongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(h.config.PongWait))
		return nil
	})

	for {
		if _, _, err := c.ws.ReadMessage(); err != nil {
			break
		}
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
