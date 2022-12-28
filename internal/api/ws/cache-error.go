package ws

import "errors"

var (
	ErrCacheRoomNotExist    = errors.New("room does not exist")
	ErrCacheConnectNotExist = errors.New("connection does not exist")
)
