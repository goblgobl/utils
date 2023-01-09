package http

/*
Our HTTP responses must do two things:
1 - Write themselves to the fasthttp.RequestCtx
2 - Optionally (but probably) enhance the logged req line

Note that RequestCtx has various ways of writing the body (e.g.
using a []byte directly, or maybe an io.Reader). It's up to
each response to figure out how it's going to interact with it.
*/

import (
	"src.goblgobl.com/utils/log"

	"github.com/valyala/fasthttp"
)

type Response interface {
	// return logger for chaining
	Write(conn *fasthttp.RequestCtx, logger log.Logger) log.Logger
}
