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
	validationLogData = log.NewField().
				Int("code", utils.RES_VALIDATION).
				Int("status", 400).
				Finalize()

	OkLogData = log.NewField().
			Int("status", 200).
			Finalize()
)

type ValidationProvider interface {
	Errors() []any
}

// body isn't known until runtime, but we know the status
// and code and can put those in logData
type DynamicResponse struct {
	Status  int
	Body    []byte
	LogData log.Field
}

func (r DynamicResponse) Write(conn *fasthttp.RequestCtx) {
	conn.SetStatusCode(r.Status)
	conn.SetBody(r.Body)
}

func (r DynamicResponse) EnhanceLog(logger log.Logger) log.Logger {
	logger.Field(r.LogData).Int("res", len(r.Body))
	return logger
}

func Validation(validator ValidationProvider) DynamicResponse {
	data := struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		Invalid []any  `json:"invalid"`
	}{
		Code:    utils.RES_VALIDATION,
		Error:   "invalid data",
		Invalid: validator.Errors(),
	}
	body, _ := json.Marshal(data)

	return DynamicResponse{
		Body:    body,
		Status:  400,
		LogData: validationLogData,
	}
}

func Ok(data any) Response {
	var body []byte
	if data != nil {
		var err error
		if body, err = json.Marshal(data); err != nil {
			se := SerializationError()
			logger := log.Error("res_ok_json").Err(err)
			se.EnhanceLog(logger).Log()
			return se
		}
	}
	return OkBytes(body)
}

func OkBytes(body []byte) DynamicResponse {
	return DynamicResponse{
		Status:  200,
		Body:    body,
		LogData: OkLogData,
	}
}
