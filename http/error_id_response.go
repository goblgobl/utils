package http

import (
	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/json"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/uuid"
)

var (
	serverErrorLogData = log.NewField().
				Int("_code", utils.RES_SERVER_ERROR).
				Int("status", 500).
				Finalize()

	serializationErrorLogData = log.NewField().
					Int("_code", utils.RES_SERIALIZATION_ERROR).
					Int("status", 500).
					Finalize()
)

type ErrorIdResponse struct {
	Err     error
	ErrorId string
	Body    []byte
	LogData log.Field
}

func NewErrorIdResponse(err error, errorId string, body []byte, logData log.Field) ErrorIdResponse {
	return ErrorIdResponse{
		Err:     err,
		ErrorId: errorId,
		Body:    body,
		LogData: logData,
	}
}

func (r ErrorIdResponse) Write(conn *fasthttp.RequestCtx, logger log.Logger) log.Logger {
	conn.SetStatusCode(500)
	conn.Response.Header.SetBytesK([]byte("Error-Id"), r.ErrorId)
	conn.SetBody(r.Body)
	return logger.
		Err(r.Err).
		Field(r.LogData).
		String("eid", r.ErrorId).
		Int("res", len(r.Body))
}

func ServerError(err error, fullError bool) Response {
	errorId := uuid.String()

	errorMessage := "internal server error"
	if fullError {
		errorMessage = err.Error()
	}

	data := struct {
		Error   string `json:"error"`
		ErrorId string `json:"error_id"`
		Code    int    `json:"code"`
	}{
		ErrorId: errorId,
		Error:   errorMessage,
		Code:    utils.RES_SERVER_ERROR,
	}
	body, _ := json.Marshal(data)
	return NewErrorIdResponse(err, errorId, body, serverErrorLogData)
}

func SerializationError(err error) Response {
	errorId := uuid.String()

	data := struct {
		Error   string `json:"error"`
		ErrorId string `json:"error_id"`
		Code    int    `json:"code"`
	}{
		ErrorId: errorId,
		Error:   "internal server error",
		Code:    utils.RES_SERIALIZATION_ERROR,
	}
	body, _ := json.Marshal(data)

	return NewErrorIdResponse(err, errorId, body, serializationErrorLogData)
}
