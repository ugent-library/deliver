package turbo

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512

	Append  Action = "append"
	Prepend Action = "prepend"
	Replace Action = "replace"
	Update  Action = "update"
	Remove  Action = "remove"
	Before  Action = "before"
	After   Action = "after"
)

type Action string

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	msgs chan []byte
}

type Stream interface {
	Action() Action
	Target() string
	TargetSelector() string
	Template() []byte
}

func (c *Client) Send(streams ...Stream) {
	if len(streams) == 0 {
		return
	}
	c.msgs <- serializeStreams(streams)
}

type Responder interface {
	Respond(Client, []byte)
}

type Hub struct {
	responder Responder
	upgrader  websocket.Upgrader
	clients   map[*Client]struct{}
	mu        sync.RWMutex
}

type Config struct {
	Responder Responder
}

func NewHub(config Config) *Hub {
	return &Hub{
		responder: config.Responder,
		clients:   make(map[*Client]struct{}),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (h *Hub) Broadcast(streams ...Stream) {
	if len(streams) == 0 {
		return
	}
	msg := serializeStreams(streams)
	for c := range h.clients {
		c.msgs <- msg
	}
}

func (h *Hub) disconnect(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &Client{
		hub:  h,
		conn: conn,
		msgs: make(chan []byte, 64),
	}

	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	go writePump(h, c)
	readPump(h, c)
}

func writePump(h *Hub, c *Client) {
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

func readPump(h *Hub, c *Client) {
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
	}
}

func serializeStreams(streams []Stream) []byte {
	b := bytes.Buffer{}
	for _, s := range streams {
		b.WriteString(`<turbo-stream action="`)
		b.Write([]byte(s.Action()))
		b.WriteString(`" `)
		if s.Target() != "" {
			b.WriteString(`target="`)
			b.WriteString(s.Target())
		} else {
			b.WriteString(`targets="`)
			b.WriteString(s.TargetSelector())
		}
		b.WriteString(`"><template>`)
		b.Write(s.Template())
		b.WriteString(`</template></turbo-stream>`)
	}
	return b.Bytes()
}
