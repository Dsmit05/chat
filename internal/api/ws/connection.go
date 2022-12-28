package ws

import (
	"github.com/Dsmit05/chat/internal/config"
	"github.com/Dsmit05/chat/internal/models"
	"log"
	"time"

	"github.com/fasthttp/websocket"
)

// Connection provides methods for working with a websocket.
type Connection struct {
	Proc   *Processor
	Conn   *websocket.Conn
	Send   chan []byte
	RoomID int
	cfg    config.WS
}

func NewConnection(proc *Processor, conn *websocket.Conn, roomID int, cfg config.WS) *Connection {
	return &Connection{Proc: proc, Conn: conn, Send: make(chan []byte, 512), RoomID: roomID, cfg: cfg}
}

// Read incoming data from client.
func (c *Connection) Read() {
	defer func() {
		c.Proc.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(c.cfg.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))
	c.Conn.SetPongHandler(
		func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))
			return nil
		})

	for {
		var msg models.Msg

		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("error: %v", err)
			break
		}

		c.Proc.Broadcast <- &Room{
			Id:        c.RoomID,
			Broadcast: []byte(msg.Text),
		}
	}
}

// Write sending messages for each client.
func (c *Connection) Write() {
	ticker := time.NewTicker(c.cfg.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			for i := 0; i < len(c.Send); i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send) // add all the messages in the channel to the client.
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
