package turbo

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

const (
	DefaultWriteTimeout  = 5 * time.Second
	DefaultMessageBuffer = 16
)

var (
	ConnectionTooSlowText   = "Connection too slow"
	InternalServerErrorText = "Internal server error"

	ErrStreamingNotSupported = errors.New("streaming not supported")
)

type client struct {
	streams []string
	msgs    chan []byte
	tooSlow chan bool
}

type Hub struct {
	config    Config
	clients   map[*client]struct{}
	streams   map[string]map[*client]struct{}
	clientsMu sync.RWMutex
	streamsMu sync.RWMutex
}

type Config struct {
	// Secret should be a random 256 bit key.
	Secret []byte
	// WriteTimeout is the time allowed to write a message. Defaults to 5
	// seconds.
	WriteTimeout time.Duration
	// MessageBuffer is the maximum number of queued messages. Defaults to 16.
	MessageBuffer int
}

// TODO rate limiter
func NewHub(config Config) *Hub {
	if config.WriteTimeout == 0 {
		config.WriteTimeout = DefaultWriteTimeout
	}
	if config.MessageBuffer == 0 {
		config.MessageBuffer = DefaultMessageBuffer
	}

	return &Hub{
		config:  config,
		clients: make(map[*client]struct{}),
		streams: make(map[string]map[*client]struct{}),
	}
}

func (h *Hub) EncryptStreamNames(names []string) (string, error) {
	return EncryptStreamNames(h.config.Secret, names)
}

func (h *Hub) DecryptStreamNames(names string) ([]string, error) {
	return DecryptStreamNames(h.config.Secret, names)
}

func (h *Hub) Broadcast(msgs ...StreamMessage) error {
	if len(msgs) == 0 {
		return nil
	}

	b, err := Encode(msgs)
	if err != nil {
		return err

	}

	for c := range h.clients {
		select {
		case c.msgs <- b:
		default:
			c.tooSlow <- true
		}
	}

	return nil
}

func (h *Hub) Send(stream string, msgs ...StreamMessage) error {
	if len(msgs) == 0 {
		return nil
	}

	b, err := Encode(msgs)
	if err != nil {
		return err
	}

	h.streamsMu.RLock()
	defer h.streamsMu.RUnlock()

	if clients, ok := h.streams[stream]; ok {
		for c := range clients {
			select {
			case c.msgs <- b:
			default:
				c.tooSlow <- true
			}
		}
	}

	return nil
}

// TODO middleware with error handler
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, cryptedStreams string) error {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	// Close with error status unless already closed normally.
	defer ws.Close(websocket.StatusInternalError, InternalServerErrorText)

	streams, err := h.DecryptStreamNames(cryptedStreams)
	if err != nil {
		return err
	}

	err = h.connectWebSocket(r.Context(), ws, streams)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if cs := websocket.CloseStatus(err); cs == websocket.StatusNormalClosure || cs == websocket.StatusGoingAway {
		return nil
	}
	return err
}

func (h *Hub) connectWebSocket(ctx context.Context, ws *websocket.Conn, streams []string) error {
	c := &client{
		streams: streams,
		msgs:    make(chan []byte, h.config.MessageBuffer),
	}

	h.addClient(c)
	defer h.removeClient(c)

	ctx = ws.CloseRead(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.tooSlow:
			return ws.Close(websocket.StatusPolicyViolation, ConnectionTooSlowText)
		case msg := <-c.msgs:
			err := writeWithTimeout(ctx, h.config.WriteTimeout, ws, msg)
			if err != nil {
				return err
			}
		}
	}
}

// TODO write timeout, error handling
func (h *Hub) HandleSSE(w http.ResponseWriter, r *http.Request, cryptedStreams string) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return ErrStreamingNotSupported
	}

	streams, err := h.DecryptStreamNames(cryptedStreams)
	if err != nil {
		return err
	}

	c := &client{
		streams: streams,
		msgs:    make(chan []byte, h.config.MessageBuffer),
	}

	h.addClient(c)
	defer h.removeClient(c)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	seq := 0

	for {
		select {
		case <-r.Context().Done():
			return r.Context().Err()
		case <-c.tooSlow:
			return nil
		case msg := <-c.msgs:
			err = writeMessage(w, seq, "message", msg)
			if err != nil {
				return err
			}
			flusher.Flush()
			seq++
		}
	}
}

func writeMessage(w io.Writer, id int, event string, b []byte) error {
	_, err := fmt.Fprintf(w, "event: %s\nid: %d\n", event, id)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		_, err = fmt.Fprintf(w, "data: %s\n", scanner.Text())
		if err != nil {
			return err
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "\n")
	return err
}

func (h *Hub) addClient(c *client) {
	h.clientsMu.Lock()
	h.clients[c] = struct{}{}
	h.clientsMu.Unlock()

	h.streamsMu.Lock()
	for _, stream := range c.streams {
		if clients, ok := h.streams[stream]; ok {
			clients[c] = struct{}{}
		} else {
			h.streams[stream] = map[*client]struct{}{c: {}}
		}
	}
	h.streamsMu.Unlock()
}

func (h *Hub) removeClient(c *client) {
	h.streamsMu.Lock()
	for _, stream := range c.streams {
		if clients, ok := h.streams[stream]; ok {
			delete(clients, c)
			if len(clients) == 0 {
				delete(h.streams, stream)
			}
		}
	}
	h.streamsMu.Unlock()

	h.clientsMu.Lock()
	delete(h.clients, c)
	h.clientsMu.Unlock()
}

func writeWithTimeout(ctx context.Context, timeout time.Duration, ws *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return ws.Write(ctx, websocket.MessageText, msg)
}
