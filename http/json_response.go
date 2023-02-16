package http

/*
Most responses are "dynamic", which is to say, the exact nature
of the response is only known at runtime. Nevertheless, there are
some things we can prepare ahead of time, namely the status code
and part of the logged data (e.g. a validation response will have
a dynamic body (the list of validation errors), but will always
have a 400 status code and the logged data will always include
the validation error code, both of which we can prepare ahead of
time.
*/

import (
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/json"
	"src.goblgobl.com/utils/log"

	"github.com/valyala/fasthttp"
)

var (
	ValidationLogData = log.NewField().
				Int("_code", utils.RES_VALIDATION).
				Int("status", 400).
				Finalize()

	OKLogData = log.NewField().
			Int("status", 200).
			Finalize()

	CreatedLogData = log.NewField().
			Int("status", 201).
			Finalize()
)

type ValidationProvider interface {
	Errors() []any
}

type JSONResponse struct {
	Status  int
	Body    []byte
	LogData log.Field
}

func NewJSONResponse(data any, status int, logData log.Field) Response {
	var body []byte
	if data != nil {
		var err error
		if body, err = json.Marshal(data); err != nil {
			return SerializationError(err)
		}
	}

	return JSONResponse{
		Status:  status,
		Body:    body,
		LogData: logData,
	}
}

func (r JSONResponse) Write(conn *fasthttp.RequestCtx, logger log.Logger) log.Logger {
	conn.SetStatusCode(r.Status)
	conn.SetBody(r.Body)
	return logger.Field(r.LogData).Int("res", len(r.Body))
}

func Validation(validator ValidationProvider) Response {
	data := struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		Invalid []any  `json:"invalid"`
	}{
		Code:    utils.RES_VALIDATION,
		Error:   "invalid data",
		Invalid: validator.Errors(),
	}
	return NewJSONResponse(data, 400, ValidationLogData)
}

func OK(data any) Response {
	return NewJSONResponse(data, 200, OKLogData)
}

func OKBytes(body []byte) Response {
	return JSONResponse{
		Status:  200,
		Body:    body,
		LogData: OKLogData,
	}
}

func Created(data any) Response {
	return NewJSONResponse(data, 201, CreatedLogData)
}
