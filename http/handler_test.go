package http

import (
	"errors"
	"testing"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/tests"
	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/typed"
)

func Test_Handler_EnvLoader_Error(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, Response, error) {
		return nil, nil, errors.New("env load fail")
	}

	conn := &fasthttp.RequestCtx{}
	logged := tests.CaptureLog(func() {
		Handler("", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
			assert.Fail(t, "next should not be called")
			return nil, nil
		})(conn)
	})

	reqLog := log.KvParse(logged)
	assert.Equal(t, reqLog["l"], "error")
	assert.Equal(t, reqLog["code"], "2001")

	errorId := reqLog["eid"]
	assert.Equal(t, len(errorId), 36)

	res := conn.Response
	assert.Equal(t, res.StatusCode(), 500)
	assert.Equal(t, string(res.Header.Peek("Error-Id")), errorId)
	assert.Equal(t, typed.Must(res.Body()).String("error_id"), errorId)
}

func Test_Handler_EnvLoader_Response(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, Response, error) {
		return nil, StaticError(61, 60, ""), nil
	}

	conn := &fasthttp.RequestCtx{}
	logged := tests.CaptureLog(func() {
		Handler("", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
			assert.Fail(t, "next should not be called")
			return nil, nil
		})(conn)
	})

	reqLog := log.KvParse(logged)
	assert.Equal(t, reqLog["l"], "req")
	assert.Equal(t, reqLog["code"], "60")

	res := conn.Response
	assert.Equal(t, res.StatusCode(), 61)
}

func Test_Handler_CallsHandlerWithEnv(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, Response, error) {
		return testEnv(200), nil, nil
	}

	conn := &fasthttp.RequestCtx{}
	Handler("", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
		assert.Equal(t, env.id, 200)
		return StaticError(2, 2, ""), nil
	})(conn)
	assert.Equal(t, conn.Response.StatusCode(), 2)
}

func Test_Handler_LogsResponse(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, Response, error) {
		return testEnv(201), nil, nil
	}

	conn := &fasthttp.RequestCtx{}
	logged := tests.CaptureLog(func() {
		Handler("test-route", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
			return StaticNotFound(9001), nil
		})(conn)
	})

	reqLog := log.KvParse(logged)
	assert.Equal(t, reqLog["l"], "req")
	assert.Equal(t, reqLog["status"], "404")
	assert.Equal(t, reqLog["res"], "33")
	assert.Equal(t, reqLog["code"], "9001")
	assert.Equal(t, reqLog["c"], "test-route")
}

func Test_Handler_LogsError(t *testing.T) {
	testLoader := func(conn *fasthttp.RequestCtx) (*TestEnv, Response, error) {
		return testEnv(202), nil, nil
	}

	conn := &fasthttp.RequestCtx{}
	logged := tests.CaptureLog(func() {
		Handler("test2", testLoader, func(conn *fasthttp.RequestCtx, env *TestEnv) (Response, error) {
			return nil, errors.New("Not Over 9000!")
		})(conn)
	})

	res := conn.Response
	assert.Equal(t, res.StatusCode(), 500)

	errorId := res.Header.Peek("Error-Id")
	assert.Equal(t, len(errorId), 36)

	reqLog := log.KvParse(logged)
	assert.Equal(t, reqLog["l"], "error")
	assert.Equal(t, reqLog["status"], "500")
	assert.Equal(t, reqLog["route"], "test2")
	assert.Equal(t, reqLog["res"], "95")
	assert.Equal(t, reqLog["code"], "2001")
	assert.Equal(t, reqLog["c"], "handler")
	assert.Equal(t, reqLog["err"], `"Not Over 9000!"`)
	assert.Equal(t, reqLog["eid"], string(errorId))
}

func Test_NoEnvHandler_LogsResponse(t *testing.T) {
	conn := &fasthttp.RequestCtx{}
	logged := tests.CaptureLog(func() {
		NoEnvHandler("test-route", func(conn *fasthttp.RequestCtx) (Response, error) {
			return StaticNotFound(9001), nil
		})(conn)
	})

	reqLog := log.KvParse(logged)
	assert.Equal(t, reqLog["l"], "req")
	assert.Equal(t, reqLog["status"], "404")
	assert.Equal(t, reqLog["c"], "test-route")
	assert.Equal(t, reqLog["res"], "33")
	assert.Equal(t, reqLog["code"], "9001")
}

func Test_NoEnvHandler_LogsError(t *testing.T) {
	conn := &fasthttp.RequestCtx{}
	logged := tests.CaptureLog(func() {
		NoEnvHandler("test2", func(conn *fasthttp.RequestCtx) (Response, error) {
			return nil, errors.New("Not Over 9000!")
		})(conn)
	})

	res := conn.Response
	assert.Equal(t, res.StatusCode(), 500)

	errorId := res.Header.Peek("Error-Id")
	assert.Equal(t, len(errorId), 36)

	reqLog := log.KvParse(logged)
	assert.Equal(t, reqLog["l"], "error")
	assert.Equal(t, reqLog["status"], "500")
	assert.Equal(t, reqLog["route"], "test2")
	assert.Equal(t, reqLog["res"], "95")
	assert.Equal(t, reqLog["code"], "2001")
	assert.Equal(t, reqLog["c"], "handler")
	assert.Equal(t, reqLog["err"], `"Not Over 9000!"`)
	assert.Equal(t, reqLog["eid"], string(errorId))
}

type TestEnv struct {
	id       int
	released bool
	logger   log.Logger
}

func testEnv(id int) *TestEnv {
	return &TestEnv{
		id:     id,
		logger: log.NewKvLogger(1024, nil, log.INFO, true),
	}
}

func (e *TestEnv) Release() {
	e.released = true
	e.logger.Release()
}

func (e TestEnv) RequestId() string {
	return ""
}

func (e TestEnv) Request(route string) log.Logger {
	return e.logger.Request(route)
}

func (e TestEnv) Error(ctx string) log.Logger {
	return e.logger.Error(ctx)
}

func assertCode(t *testing.T, conn *fasthttp.RequestCtx, expected int) {
	t.Helper()
	res := conn.Response
	body := res.Body()
	json, _ := typed.Json(body)
	assert.Equal(t, json.Int("code"), expected)
}
