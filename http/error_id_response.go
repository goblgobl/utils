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
				Int("code", utils.RES_SERVER_ERROR).
				Int("status", 500).
				Finalize()

	serializationErrorLogData = log.NewField().
					Int("code", utils.RES_SERIALIZATION_ERROR).
					Int("status", 500).
					Finalize()
)

type ErrorIdResponse struct {
	err     error
	ErrorId string
	Body    []byte
	LogData log.Field
}

func (r ErrorIdResponse) Write(conn *fasthttp.RequestCtx, logger log.Logger) log.Logger {
	conn.SetStatusCode(500)
	conn.Response.Header.SetBytesK([]byte("Error-Id"), r.ErrorId)
	conn.SetBody(r.Body)
	return logger.
		Err(r.err).
		Field(r.LogData).
		String("eid", r.ErrorId).
		Int("res", len(r.Body))
}

func ServerError(err error) Response {
	errorId := uuid.String()

	data := struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		ErrorId string `json:"error_id"`
	}{
		ErrorId: errorId,
		Code:    utils.RES_SERVER_ERROR,
		Error:   "internal server error",
	}
	body, _ := json.Marshal(data)

	return ErrorIdResponse{
		err:     err,
		Body:    body,
		ErrorId: errorId,
		LogData: serverErrorLogData,
	}
}

func SerializationError(err error) Response {
	errorId := uuid.String()

	data := struct {
		Code    int    `json:"code"`
		Error   string `json:"error"`
		ErrorId string `json:"error_id"`
	}{
		ErrorId: errorId,
		Code:    utils.RES_SERIALIZATION_ERROR,
		Error:   "internal server error",
	}
	body, _ := json.Marshal(data)

	return ErrorIdResponse{
		err:     err,
		Body:    body,
		ErrorId: errorId,
		LogData: serializationErrorLogData,
	}
}
