package http

import (
	"time"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/log"
)

type Env interface {
	Release()
	RequestId() string
	Request(string) log.Logger
	Error(string) log.Logger
	ServerError(err error, conn *fasthttp.RequestCtx) Response
}

func Handler[T Env](routeName string, loadEnv func(ctx *fasthttp.RequestCtx) (T, Response, error), next func(ctx *fasthttp.RequestCtx, env T) (Response, error)) func(ctx *fasthttp.RequestCtx) {
	return func(conn *fasthttp.RequestCtx) {
		start := time.Now()

		var logger log.Logger
		env, res, err := loadEnv(conn)

		header := &conn.Response.Header
		header.SetContentTypeBytes([]byte("application/json"))

		if err != nil {
			res = ServerError(err, false)
		}
		if res == nil {
			// we can only be here if loadEnv didn't return a response or an error
			// (which means it should have returned an env)
			defer env.Release()
			logger = env.Request(routeName)
			header.SetBytesK([]byte("RequestId"), env.RequestId())
			res, err = next(conn, env)
			if err != nil {
				res = env.ServerError(err, conn)
			}
		} else {
			logger = log.Request(routeName)
		}

		res.Write(conn, logger).
			Int64("ms", time.Now().Sub(start).Milliseconds()).
			Log()
	}
}

func NoEnvHandler(routeName string, next func(ctx *fasthttp.RequestCtx) (Response, error)) func(ctx *fasthttp.RequestCtx) {
	return func(conn *fasthttp.RequestCtx) {
		start := time.Now()
		var logger log.Logger

		header := &conn.Response.Header
		header.SetContentTypeBytes([]byte("application/json"))

		res, err := next(conn)

		if err == nil {
			logger = log.Request(routeName)
		} else {
			res = ServerError(err, false)
			logger = log.Error("handler").String("route", routeName)
		}

		res.Write(conn, logger).
			Int64("ms", time.Now().Sub(start).Milliseconds()).
			Log()
	}
}
