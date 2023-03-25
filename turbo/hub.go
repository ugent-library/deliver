package turbo

import (
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

var ErrInvalidStreamNames = errors.New("invalid stream names")

type client struct {
	ws      *websocket.Conn
	streams []string
	msgs    chan []byte
}

type Hub struct {
	config    Config
	upgrader  websocket.Upgrader
	clients   map[*client]struct{}
	streams   map[string]map[*client]struct{}
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

// TODO context, error handling
func (h *Hub) Broadcast(msgs ...StreamMessage) {
	if len(msgs) == 0 {
		return
	}
	b, err := Encode(msgs)
	if err != nil {
		return
	}
	for c := range h.clients {
		c.msgs <- b
	}
}

// TODO context, error handling
func (h *Hub) Send(stream string, msgs ...StreamMessage) {
	if len(msgs) == 0 {
		return
	}

	b, err := Encode(msgs)
	if err != nil {
		log.Print(err)
		return
	}

	h.streamsMu.RLock()
	defer h.streamsMu.RUnlock()

	if clients, ok := h.streams[stream]; ok {
		for c := range clients {
			c.msgs <- b
		}
	}
}

// TODO error handler
// TODO middleware
func (h *Hub) Handle(w http.ResponseWriter, r *http.Request, cryptedStreams string) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	streams, err := h.DecryptStreamNames(cryptedStreams)
	if err != nil {
		return err
	}

	c := &client{
		ws:      conn,
		streams: streams,
		msgs:    make(chan []byte),
	}

	h.addClient(c, streams)

	go wsWrite(h, c)
	go wsRead(h, c)

	return nil
}

func (h *Hub) addClient(c *client, streams []string) {
	h.clientsMu.Lock()
	h.clients[c] = struct{}{}
	h.clientsMu.Unlock()

	h.streamsMu.Lock()
	for _, stream := range streams {
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

// TODO logging
func wsWrite(h *Hub, c *client) {
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

// TODO logging
func wsRead(h *Hub, c *client) {
	defer func() {
		h.removeClient(c)
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
