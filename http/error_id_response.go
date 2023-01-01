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
	ErrorId string
	Body    []byte
	LogData log.Field
}

func (r ErrorIdResponse) Write(conn *fasthttp.RequestCtx) {
	conn.SetStatusCode(500)
	conn.Response.Header.SetBytesK([]byte("Error-Id"), r.ErrorId)
	conn.SetBody(r.Body)
}

func (r ErrorIdResponse) EnhanceLog(logger log.Logger) log.Logger {
	logger.Field(r.LogData).String("eid", r.ErrorId).Int("res", len(r.Body))
	return logger
}

func ServerError() Response {
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
		Body:    body,
		ErrorId: errorId,
		LogData: serverErrorLogData,
	}
}

func SerializationError() Response {
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
		Body:    body,
		ErrorId: errorId,
		LogData: serializationErrorLogData,
	}
}
