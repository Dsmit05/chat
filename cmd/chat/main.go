package main

import (
	"github.com/Dsmit05/chat/internal/api"
	"github.com/Dsmit05/chat/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := &config.App{
		WS: config.WS{
			MaxMessageSize: 1024,
			PongWait:       90 * time.Second,
			PingPeriod:     60 * time.Second,
			WriteWait:      10 * time.Second,
			Workers:        4,
		},
		Server: config.Server{
			Host: "",
			Port: 8080,
		},
	}

	s := api.NewApiServer(cfg)
	go s.Start()
	defer s.Stop()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
