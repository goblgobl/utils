package validation

import (
	"src.goblgobl.com/utils"
)

var (
	Required      = &Invalid{Code: utils.VAL_REQUIRED, Error: "required"}
	TypeInt       = &Invalid{Code: utils.VAL_INT_TYPE, Error: "must be an integer"}
	TypeUUID      = &Invalid{Code: utils.VAL_UUID_TYPE, Error: "must be a uuid"}
	TypeBool      = &Invalid{Code: utils.VAL_BOOL_TYPE, Error: "must be a boolean"}
	TypeArray     = &Invalid{Code: utils.VAL_ARRAY_TYPE, Error: "must be an array"}
	TypeString    = &Invalid{Code: utils.VAL_STRING_TYPE, Error: "must be a string"}
	TypeFloat     = &Invalid{Code: utils.VAL_FLOAT_TYPE, Error: "must be a number"}
	TypeObject    = &Invalid{Code: utils.VAL_OBJECT_TYPE, Error: "must be an object"}
	StringPattern = &Invalid{Code: utils.VAL_STRING_PATTERN, Error: "is not valid"}
)

type Invalid struct {
	Data  any    `json:"data"` //TODO: add omitempty https://github.com/goccy/go-json/issues/391
	Error string `json:"error"`
	Code  uint32 `json:"code"`
}

type InvalidField struct {
	*Invalid
	Field string `json:"field"`
}

func InvalidStringPattern(message ...string) *Invalid {
	if message == nil {
		return StringPattern
	}
	return &Invalid{
		Error: message[0],
		Code:  utils.VAL_STRING_PATTERN,
	}
}

// Used to create a data field of type: `data: {min: #}`
// Since this is a common "data" to have (e.g. min string length, min integer)
// it helps to have this type with the accompanying MinData func to ensure consistency
type minData struct {
	Min any `json:"min"`
}

func MinData(min any) minData {
	return minData{min}
}

// see minData types for description
type maxData struct {
	Max any `json:"max"`
}

func MaxData(max any) maxData {
	return maxData{max}
}

// see minData types for description
type valueData struct {
	Value any `json:"value"`
}

func ValueData(value any) valueData {
	return valueData{value}
}

// see minData types for description
type rangeData struct {
	Min any `json:"min"`
	Max any `json:"max"`
}

func RangeData(min any, max any) rangeData {
	return rangeData{min, max}
}

// see minData types for description
type choiceData[T any] struct {
	Valid []T `json:"valid"`
}

func ChoiceData[T any](valid []T) choiceData[T] {
	return choiceData[T]{Valid: valid}
}
