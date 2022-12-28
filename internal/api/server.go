package api

import (
	"context"
	"github.com/Dsmit05/chat/internal/api/controllers"
	"github.com/Dsmit05/chat/internal/api/ws"
	"github.com/Dsmit05/chat/internal/config"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	fast      fasthttp.Server
	processor *ws.Processor
	cfg       *config.App
	cache     *ws.RoomCache
}

func NewApiServer(cfg *config.App) *Server {
	cacheRoom := ws.NewRoomCache(10)
	processor := ws.NewProcessor(cacheRoom)

	for w := 0; w < cfg.Workers; w++ {
		go processor.Run()
	}

	server := &Server{
		fast:      fasthttp.Server{},
		processor: processor,
		cfg:       cfg,
		cache:     cacheRoom,
	}
	server.initHandlers()

	return server
}

func (s *Server) Start() {
	s.fast.ListenAndServe(s.cfg.GetServerAddr())
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stop := make(chan bool)
	go func() {
		s.fast.Shutdown()
		stop <- true
	}()

	select {
	case <-ctx.Done():
		log.Printf("server not gracefully stop, error: %v", ctx.Err())
	case <-stop:
		log.Println("server gracefully stop")
	}
}

func (s *Server) initHandlers() {
	wsControllers := controllers.NewWS(s.cfg)
	infoControllers := controllers.NewInfo(s.cache)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		p := strings.Split(string(ctx.Path()), "/")[1:]
		n := len(p)
		switch {
		case n == 1 && p[0] == "":
			infoControllers.Statistic(ctx)
		case n == 2 && p[0] == "rooms" && func() bool { // /rooms/{id}
			id, err := strconv.Atoi(p[1])
			if err != nil {
				return false
			}
			ctx.SetUserValue("id", id)

			return true
		}():
			wsControllers.Chat(ctx, s.processor)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)

		}
	}

	s.fast.Handler = requestHandler
}
