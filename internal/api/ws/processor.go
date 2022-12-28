package ws

import (
	"fmt"
	"log"
)

// Processor contains a set of rooms and active connections
// Performs the role of an intermediate layer for processing various events(Register, Unregister, Broadcast).
type Processor struct {
	clients    CacheInterface
	Broadcast  chan *Room
	Register   chan *Connection
	Unregister chan *Connection
}

func NewProcessor(roomCache CacheInterface) *Processor {
	return &Processor{
		Register:   make(chan *Connection),
		Unregister: make(chan *Connection),
		clients:    roomCache,

		// the channel size should be greater than the number of goroutines of the function Run()
		Broadcast: make(chan *Room, 15),
	}
}

func (p *Processor) Run() {
	for {
		select {
		case conn := <-p.Register:
			if p.clients.Set(conn.RoomID, conn) {
				log.Println("new room has been created")
			}

			p.Broadcast <- &Room{
				Id:        conn.RoomID,
				Broadcast: []byte(fmt.Sprintf("A new user is in the chat %s", getRemoteAddr(conn))),
			}

		case conn := <-p.Unregister:
			close(conn.Send)
			err := p.clients.DelConnectionAndRoomIfZero(conn.RoomID, conn)
			if err != nil {
				log.Printf("DelConnectionAndRoomIfZero error: %v", err)
				continue
			}

			msg := fmt.Sprintf("user %s out of the chat", getRemoteAddr(conn))

			p.Broadcast <- &Room{
				Id:        conn.RoomID,
				Broadcast: []byte(msg),
			}

		case room := <-p.Broadcast:
			connections, err := p.clients.GetConnectionsFromRoom(room.Id)
			if err != nil {
				log.Printf("GetConnectionsFromRoom error: %v", err)
				continue
			}

			for _, conn := range connections {
				select {
				case conn.Send <- room.Broadcast:
				default:
					close(conn.Send)
					err = p.clients.DelConnection(conn.RoomID, conn)
					if err != nil {
						log.Printf("DelConnection error: %v", err)
						continue
					}
				}
			}
		}
	}
}

func getRemoteAddr(conn *Connection) (remoteAddr string) {
	defer func() { // bug in fasthttp, panic in github.com/valyala/fasthttp.(*hijackConn).RemoteAddr(0x1?)
		if r := recover(); r != nil {
			err := fmt.Errorf("recover(): %v", r)
			log.Printf("Panic: %v", err)
		}
	}()

	if conn == nil {
		return
	}

	wsConn := conn.Conn
	if wsConn == nil {
		return
	}

	addr := wsConn.RemoteAddr()
	if addr == nil {
		return
	}

	remoteAddr = addr.String()

	return
}
