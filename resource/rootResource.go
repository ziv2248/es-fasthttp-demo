package resource

import (
	"time"

	"github.com/valyala/fasthttp"
)

type RootResource struct {
}

func (r *RootResource) Ping(ctx *fasthttp.RequestCtx) {
	ctx.Success("text/plain", []byte("pong"))
	time.Sleep(time.Duration(5) * time.Second)
}
