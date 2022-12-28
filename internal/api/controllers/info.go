package controllers

import (
	"github.com/valyala/fasthttp"
)

type statistic interface {
	Statistic() string
}

type Info struct {
	stat statistic
}

func NewInfo(stat statistic) *Info {
	return &Info{stat: stat}
}

func (f *Info) Statistic(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	ctx.SuccessString("text/plain; charset=utf8", f.stat.Statistic())
}
