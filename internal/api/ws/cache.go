package ws

import (
	"strconv"
	"strings"
	"sync"
)

type RoomCache struct {
	m  map[int]map[*Connection]struct{}
	rw *sync.RWMutex
}

func NewRoomCache(startSize uint) *RoomCache {
	return &RoomCache{m: make(map[int]map[*Connection]struct{}, startSize), rw: new(sync.RWMutex)}
}

func (r *RoomCache) GetConnectionsFromRoom(roomID int) ([]*Connection, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	room, ok := r.m[roomID]
	if !ok {
		return nil, ErrCacheRoomNotExist
	}

	connections := make([]*Connection, 0, len(room))

	for conn, _ := range room {
		connections = append(connections, conn)
	}

	return connections, nil
}

// Set add new Connection in cache, return true if create new room.
func (r *RoomCache) Set(roomID int, conn *Connection) (newRoom bool) {
	r.rw.Lock()
	if _, ok := r.m[roomID]; !ok {
		newRoom = true
		r.m[roomID] = make(map[*Connection]struct{})
	}

	r.m[roomID][conn] = struct{}{}

	r.rw.Unlock()

	return
}

func (r *RoomCache) DelRoom(roomID int) error {
	r.rw.RLock()
	if _, ok := r.m[roomID]; !ok {
		r.rw.RUnlock()
		return ErrCacheRoomNotExist
	}
	r.rw.RUnlock()

	r.rw.Lock()
	delete(r.m, roomID)
	r.rw.Unlock()

	return nil
}

func (r *RoomCache) DelConnection(roomID int, conn *Connection) error {
	r.rw.RLock()
	if _, ok := r.m[roomID]; !ok {
		r.rw.RUnlock()
		return ErrCacheRoomNotExist
	}

	if _, ok := r.m[roomID][conn]; !ok {
		r.rw.RUnlock()
		return ErrCacheConnectNotExist
	}
	r.rw.RUnlock()

	r.rw.Lock()
	delete(r.m[roomID], conn)
	r.rw.Unlock()

	return nil
}

// DelConnectionAndRoomIfZero deletes the connection, and if it was the last one, deletes the room.
func (r *RoomCache) DelConnectionAndRoomIfZero(roomID int, conn *Connection) error {
	if err := r.DelConnection(roomID, conn); err != nil {
		return err
	}

	if len(r.m[roomID]) == 0 {
		if err := r.DelRoom(roomID); err != nil {
			return err
		}
	}

	return nil
}

func (r *RoomCache) Statistic() string {
	r.rw.RLock()
	defer r.rw.RUnlock()

	builder := strings.Builder{}
	builder.WriteString("List of rooms with the number of clients:\n")
	for room, connections := range r.m {
		builder.WriteString("room №: ")
		builder.WriteString(strconv.Itoa(room))
		builder.WriteString(", with active clients: ")
		builder.WriteString(strconv.Itoa(len(connections)))
		builder.WriteString("\n")
	}

	return builder.String()
}
