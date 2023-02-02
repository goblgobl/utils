package http

import (
	"errors"
	"testing"

	"src.goblgobl.com/tests/assert"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/typed"
	"src.goblgobl.com/utils/validation"

	"github.com/valyala/fasthttp"
)

type TestResponse struct {
	status int
	body   string
	json   typed.Typed
	log    map[string]string
}

func Test_Ok_NoBody(t *testing.T) {
	res := read(Ok(nil))
	assert.Equal(t, res.status, 200)
	assert.Equal(t, len(res.body), 0)
	assert.Equal(t, res.log["res"], "0")
	assert.Equal(t, res.log["status"], "200")
}

func Test_Ok_Body(t *testing.T) {
	res := read(Ok(map[string]any{"over": 9000}))
	assert.Equal(t, res.status, 200)
	assert.Equal(t, res.body, `{"over":9000}`)
	assert.Equal(t, res.log["res"], "13")
	assert.Equal(t, res.log["status"], "200")
}

func Test_Ok_InvalidBody(t *testing.T) {
	res := read(Ok(make(chan bool)))

	errorId := res.log["eid"]
	assert.Equal(t, len(errorId), 36)
	assert.Equal(t, res.log["_code"], "2002")
	assert.Equal(t, res.log["res"], "95")
	assert.Equal(t, res.log["status"], "500")

	assert.Equal(t, res.status, 500)
	assert.Equal(t, res.json.Int("code"), 2002)
	assert.Equal(t, res.json.String("error"), "internal server error")
	assert.Equal(t, res.json.String("error_id"), errorId)

}

func Test_StaticNotFound(t *testing.T) {
	res := read(StaticNotFound(1023))
	assert.Equal(t, res.status, 404)
	assert.Equal(t, res.body, `{"code":1023,"error":"not found"}`)
	assert.Equal(t, res.log["_code"], "1023")
	assert.Equal(t, res.log["res"], "33")
	assert.Equal(t, res.log["status"], "404")
}

func Test_StaticError(t *testing.T) {
	res := read(StaticError(511, 1002, "oops"))
	assert.Equal(t, res.status, 511)
	assert.Equal(t, res.body, `{"code":1002,"error":"oops"}`)
	assert.Equal(t, res.log["_code"], "1002")
	assert.Equal(t, res.log["res"], "28")
	assert.Equal(t, res.log["status"], "511")
}

func Test_ServerError_NotFullError(t *testing.T) {
	res := read(ServerError(errors.New("an_error1"), false))
	assert.Equal(t, res.status, 500)
	assert.Equal(t, res.json.Int("code"), 2001)
	assert.Equal(t, res.json.String("error"), "internal server error")

	errorId := res.json.String("error_id")
	assert.Equal(t, len(errorId), 36)

	assert.Equal(t, res.log["_err"], "an_error1")
	assert.Equal(t, res.log["_code"], "2001")
	assert.Equal(t, res.log["res"], "95")
	assert.Equal(t, res.log["status"], "500")
	assert.Equal(t, res.log["eid"], errorId)
}

func Test_ServerError_FullError(t *testing.T) {
	res := read(ServerError(errors.New("an_error1"), true))
	assert.Equal(t, res.status, 500)
	assert.Equal(t, res.json.Int("code"), 2001)
	assert.Equal(t, res.json.String("error"), "an_error1")

	errorId := res.json.String("error_id")
	assert.Equal(t, len(errorId), 36)

	assert.Equal(t, res.log["_err"], "an_error1")
	assert.Equal(t, res.log["_code"], "2001")
	assert.Equal(t, res.log["res"], "83")
	assert.Equal(t, res.log["status"], "500")
	assert.Equal(t, res.log["eid"], errorId)
}

func Test_Validation(t *testing.T) {
	rules := validation.Object().
		Field("field1", validation.String().Required()).
		Field("field2", validation.Int().Min(10))

	result := validation.NewResult(5)
	rules.Validate(map[string]any{"over": 9000, "field2": 3}, result)

	res := read(Validation(result))
	assert.Equal(t, res.status, 400)
	assert.Equal(t, res.json.Int("code"), 2004)
	assert.Equal(t, res.json.String("error"), "invalid data")

	invalid := res.json.Objects("invalid")
	assert.Equal(t, len(invalid), 2)
	assert.Equal(t, invalid[0].Int("code"), 1001)
	assert.Equal(t, invalid[0].String("field"), "field1")
	assert.Equal(t, invalid[0].String("error"), "required")
	assert.Nil(t, invalid[0].Object("data"))

	assert.Equal(t, invalid[1].Int("code"), 1006)
	assert.Equal(t, invalid[1].String("field"), "field2")
	assert.Equal(t, invalid[1].String("error"), "must be greater or equal to 10")
	assert.Equal(t, invalid[1].Object("data").Int("min"), 10)

	assert.Equal(t, res.log["_code"], "2004")
	assert.Equal(t, res.log["res"], "200")
	assert.Equal(t, res.log["status"], "400")
}

func read(res Response) TestResponse {
	conn := &fasthttp.RequestCtx{}
	logger := res.Write(conn, log.Request("test"))
	defer logger.Release()

	body := conn.Response.Body()
	var json typed.Typed
	if len(body) > 0 {
		json = typed.Must(body)
	}

	return TestResponse{
		json:   json,
		body:   string(body),
		status: conn.Response.StatusCode(),
		log:    log.KvParse(string(logger.Bytes())),
	}
}
