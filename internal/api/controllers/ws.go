package controllers

import (
	"github.com/Dsmit05/chat/internal/api/ws"
	"github.com/Dsmit05/chat/internal/config"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"log"
)

type WS struct {
	cfg *config.App
}

func NewWS(cfg *config.App) *WS {
	return &WS{cfg: cfg}
}

var fastUpgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (w *WS) Chat(ctx *fasthttp.RequestCtx, proc *ws.Processor) {
	err := fastUpgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		id := ctx.UserValue("id").(int)

		cliConn := ws.NewConnection(proc, conn, id, w.cfg.WS)
		cliConn.Proc.Register <- cliConn

		go cliConn.Write()
		cliConn.Read()
	})

	if err != nil {
		log.Println(err)
	}
}
