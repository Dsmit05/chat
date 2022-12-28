package config

import (
	"fmt"
	"time"
)

type WS struct {
	MaxMessageSize int64
	PongWait       time.Duration // time to read the next message
	PingPeriod     time.Duration // client ping time
	WriteWait      time.Duration // time to write message client
	Workers        int
}

type Server struct {
	Host string
	Port int
}

type App struct {
	WS
	Server
}

func NewAppConfig() *App {
	return &App{}
}

func (a *App) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
