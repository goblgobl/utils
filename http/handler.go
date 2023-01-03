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
}

func Handler[T Env](routeName string, loadEnv func(ctx *fasthttp.RequestCtx) (T, Response, error), next func(ctx *fasthttp.RequestCtx, env T) (Response, error)) func(ctx *fasthttp.RequestCtx) {
	return func(conn *fasthttp.RequestCtx) {
		start := time.Now()

		var haveEnv bool
		var logger log.Logger
		env, res, err := loadEnv(conn)

		header := &conn.Response.Header
		header.SetContentTypeBytes([]byte("application/json"))

		if res == nil && err == nil {
			haveEnv = true
			defer env.Release()
			header.SetBytesK([]byte("RequestId"), env.RequestId())
			res, err = next(conn, env)
		}

		if err == nil {
			if haveEnv {
				logger = env.Request(routeName)
			} else {
				logger = log.Request(routeName)
			}
		} else {
			if haveEnv {
				logger = env.Error("handler").Err(err).String("route", routeName)
			} else {
				logger = log.Error("handler").Err(err).String("route", routeName)
			}
			res = ServerError()
		}

		res.Write(conn)
		res.EnhanceLog(logger).
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
			logger = log.Info("req")
		} else {
			res = ServerError()
			logger = log.Error("handler").Err(err)
		}

		res.Write(conn)
		res.EnhanceLog(logger).
			String("route", routeName).
			Int64("ms", time.Now().Sub(start).Milliseconds()).
			Log()
	}
}
