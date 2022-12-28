package ws

type CacheInterface interface {
	Set(roomID int, conn *Connection) (newRoom bool)
	DelRoom(roomID int) error
	DelConnection(roomID int, conn *Connection) error
	GetConnectionsFromRoom(roomID int) ([]*Connection, error)
	DelConnectionAndRoomIfZero(roomID int, conn *Connection) error
}
