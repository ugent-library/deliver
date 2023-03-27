package turbo

import (
	"bufio"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

	ErrInvalidStreamNames    = errors.New("invalid stream names")
	ErrStreamingNotSupported = errors.New("streaming not supported")
)

type client struct {
	streams []string
	msgs    chan []byte
	tooSlow chan bool
}

type Hub struct {
	config Config
	// upgrader  websocket.Upgrader
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
// TODO support SSE
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
		return nil, ErrInvalidStreamNames
	}

	// Split cryptedMsg in nonce and encrypted message and use gcm.Open() to
	// decrypt and authenticate the data.
	msg, err := gcm.Open(nil, cryptedMsg[:nonceSize], cryptedMsg[nonceSize:], nil)
	if err != nil {
		return nil, errors.Join(err, ErrInvalidStreamNames)
	}

	return strings.Split(string(msg), ","), nil
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

// TODO write timeout
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
			if err := r.Context().Err(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return err
			}
			w.WriteHeader(http.StatusOK)
			return nil
		case <-c.tooSlow:
			w.WriteHeader(http.StatusRequestTimeout)
			return nil
		case msg := <-c.msgs:
			err = writeMessage(w, seq, "message", msg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
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
	if err != nil {
		return err
	}

	return nil
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
