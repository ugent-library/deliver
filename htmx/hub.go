package htmx

import (
	"context"
	"errors"
	"net/http"
	"strings"
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
	channels []string
	msgs     chan string
	tooSlow  chan bool
}

type Hub struct {
	config     Config
	clients    map[*client]struct{}
	channels   map[string]map[*client]struct{}
	clientsMu  sync.RWMutex
	channelsMu sync.RWMutex
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
		config:   config,
		clients:  make(map[*client]struct{}),
		channels: make(map[string]map[*client]struct{}),
	}
}

func (h *Hub) EncryptChannelNames(names []string) (string, error) {
	return EncryptChannelNames(h.config.Secret, names)
}

func (h *Hub) DecryptChannelNames(names string) ([]string, error) {
	return DecryptChannelNames(h.config.Secret, names)
}

func (h *Hub) Broadcast(msgs ...string) error {
	if len(msgs) == 0 {
		return nil
	}

	msg := strings.Join(msgs, "")

	for c := range h.clients {
		select {
		case c.msgs <- msg:
		default:
			c.tooSlow <- true
		}
	}

	return nil
}

func (h *Hub) Send(channel string, msgs ...string) error {
	if len(msgs) == 0 {
		return nil
	}

	msg := strings.Join(msgs, "")

	h.channelsMu.RLock()
	defer h.channelsMu.RUnlock()

	if clients, ok := h.channels[channel]; ok {
		for c := range clients {
			select {
			case c.msgs <- msg:
			default:
				c.tooSlow <- true
			}
		}
	}

	return nil
}

// TODO middleware with error handler
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, cryptedChannels string) error {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	// Close with error status unless already closed normally.
	defer ws.Close(websocket.StatusInternalError, InternalServerErrorText)

	channels, err := h.DecryptChannelNames(cryptedChannels)
	if err != nil {
		return err
	}

	err = h.connectWebSocket(r.Context(), ws, channels)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if cs := websocket.CloseStatus(err); cs == websocket.StatusNormalClosure || cs == websocket.StatusGoingAway {
		return nil
	}
	return err
}

func (h *Hub) connectWebSocket(ctx context.Context, ws *websocket.Conn, channels []string) error {
	c := &client{
		channels: channels,
		msgs:     make(chan string, h.config.MessageBuffer),
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
// func (h *Hub) HandleSSE(w http.ResponseWriter, r *http.Request, cryptedChannels string) error {
// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		return ErrStreamingNotSupported
// 	}

// 	channels, err := h.DecryptChannelNames(cryptedChannels)
// 	if err != nil {
// 		return err
// 	}

// 	c := &client{
// 		channels: channels,
// 		msgs:     make(chan string, h.config.MessageBuffer),
// 	}

// 	h.addClient(c)
// 	defer h.removeClient(c)

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	seq := 0

// 	for {
// 		select {
// 		case <-r.Context().Done():
// 			return r.Context().Err()
// 		case <-c.tooSlow:
// 			return nil
// 		case msg := <-c.msgs:
// 			err = writeMessage(w, seq, "message", msg)
// 			if err != nil {
// 				return err
// 			}
// 			flusher.Flush()
// 			seq++
// 		}
// 	}
// }

// func writeMessage(w io.Writer, id int, event, msg string) error {
// 	_, err := fmt.Fprintf(w, "event: %s\nid: %d\n", event, id)
// 	if err != nil {
// 		return err
// 	}

// 	scanner := bufio.NewScanner(strings.NewReader(msg))
// 	for scanner.Scan() {
// 		_, err = fmt.Fprintf(w, "data: %s\n", scanner.Text())
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if err = scanner.Err(); err != nil {
// 		return err
// 	}

// 	_, err = fmt.Fprint(w, "\n")
// 	return err
// }

func (h *Hub) addClient(c *client) {
	h.clientsMu.Lock()
	h.clients[c] = struct{}{}
	h.clientsMu.Unlock()

	h.channelsMu.Lock()
	for _, stream := range c.channels {
		if clients, ok := h.channels[stream]; ok {
			clients[c] = struct{}{}
		} else {
			h.channels[stream] = map[*client]struct{}{c: {}}
		}
	}
	h.channelsMu.Unlock()
}

func (h *Hub) removeClient(c *client) {
	h.channelsMu.Lock()
	for _, stream := range c.channels {
		if clients, ok := h.channels[stream]; ok {
			delete(clients, c)
			if len(clients) == 0 {
				delete(h.channels, stream)
			}
		}
	}
	h.channelsMu.Unlock()

	h.clientsMu.Lock()
	delete(h.clients, c)
	h.clientsMu.Unlock()
}

func writeWithTimeout(ctx context.Context, timeout time.Duration, ws *websocket.Conn, msg string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return ws.Write(ctx, websocket.MessageText, []byte(msg))
}
