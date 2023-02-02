package http

/*
Generic errors that don't chnange. Things like internal server error.
*/

import (
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/json"
	"src.goblgobl.com/utils/log"

	"github.com/valyala/fasthttp"
)

var (
	InvalidJSON = StaticError(400, utils.RES_INVALID_JSON_PAYLOAD, "invalid json payload")
)

// We know the status/body/logData upfront (lets us optimize
// EnhanceLog)
type StaticResponse struct {
	status  int
	body    []byte
	logData log.Field
}

func (r StaticResponse) Write(conn *fasthttp.RequestCtx, logger log.Logger) log.Logger {
	conn.SetStatusCode(r.status)
	conn.SetBody(r.body)
	return logger.Field(r.logData)
}

func StaticError(status int, code int, error string) StaticResponse {
	data := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  code,
		Error: error,
	}
	body, err := json.Marshal(data)
	if err != nil {
		// static errors should only be called at startup
		panic(err)
	}

	logData := log.NewField().
		Int("_code", code).
		Int("status", status).
		Int("res", len(body)).
		Finalize()

	return StaticResponse{
		body:    body,
		status:  status,
		logData: logData,
	}
}

func StaticNotFound(code int) StaticResponse {
	return StaticError(404, code, "not found")
}
