package validation

import (
	"fmt"

	"src.goblgobl.com/utils"
)

var (
	globalPool *Pool
)

func Required() Invalid {
	return Invalid{
		Code:  utils.VAL_REQUIRED,
		Error: "required",
	}
}

func InvalidStringType() Invalid {
	return Invalid{
		Code:  utils.VAL_STRING_TYPE,
		Error: "must be a string",
	}
}

func InvalidStringLength(min int, max int) Invalid {
	return Invalid{
		Code:  utils.VAL_STRING_LEN,
		Data:  Range(min, max),
		Error: fmt.Sprintf("must be between %d and %d characters", min, max),
	}
}

func InvalidStringPattern(errorMessage ...string) Invalid {
	err := "is not valid"
	if errorMessage != nil {
		err = errorMessage[0]
	}

	return Invalid{
		Error: err,
		Code:  utils.VAL_STRING_PATTERN,
	}
}

func InvalidStringChoice(choices []string) Invalid {
	return Invalid{
		Code:  utils.VAL_STRING_CHOICE,
		Error: "is not a valid choice",
		Data:  Choice(choices),
	}
}

func InvalidIntType() Invalid {
	return Invalid{
		Code:  utils.VAL_INT_TYPE,
		Error: "must be a number",
	}
}

func InvalidIntMin(min int) Invalid {
	return Invalid{
		Code:  utils.VAL_INT_MIN,
		Error: fmt.Sprintf("must be greater or equal to %d", min),
		Data:  Min(min),
	}
}

func InvalidIntMax(max int) Invalid {
	return Invalid{
		Code:  utils.VAL_INT_MAX,
		Error: fmt.Sprintf("must be less than or equal to %d", max),
		Data:  Max(max),
	}
}

func InvalidIntRange(min int, max int) Invalid {
	return Invalid{
		Code:  utils.VAL_INT_RANGE,
		Error: fmt.Sprintf("must be between %d and %d", min, max),
		Data:  Range(min, max),
	}
}

func InvalidBoolType() Invalid {
	return Invalid{
		Code:  utils.VAL_BOOL_TYPE,
		Error: "must be true or false",
	}
}

func InvalidUUIDType() Invalid {
	return Invalid{
		Code:  utils.VAL_UUID_TYPE,
		Error: "must be a uuid",
	}
}

func InvalidArrayType() Invalid {
	return Invalid{
		Code:  utils.VAL_ARRAY_TYPE,
		Error: "must be a list of values",
	}
}

func InvalidArrayMinLength(min int) Invalid {
	return Invalid{
		Code:  utils.VAL_ARRAY_MIN_LENGTH,
		Error: fmt.Sprintf("must have at least %d values", min),
		Data:  Min(min),
	}
}

func InvalidArrayMaxLength(max int) Invalid {
	return Invalid{
		Code:  utils.VAL_ARRAY_MAX_LENGTH,
		Error: fmt.Sprintf("must have no more than %d values", max),
		Data:  Max(max),
	}
}

func InvalidArrayRangeLength(min int, max int) Invalid {
	return Invalid{
		Code:  utils.VAL_ARRAY_RANGE_LENGTH,
		Error: fmt.Sprintf("must have between %d and %d values", min, max),
		Data:  Range(min, max),
	}
}

// Not used by this library (yet), but is used by callers of this library.
// We put it here because (a) we expect to use this eventually here and (b)
// we want this to be consistent across all users of this library
func InvalidArrayStringType() Invalid {
	return Invalid{
		Code:  utils.VAL_ARRAY_TYPE_STRING,
		Error: fmt.Sprintf("must be a list of strings"),
	}
}

func InvalidFloatType() Invalid {
	return Invalid{
		Code:  utils.VAL_FLOAT_TYPE,
		Error: "must be a float",
	}
}

func InvalidFloatMin(min float64) Invalid {
	return Invalid{
		Code:  utils.VAL_FLOAT_MIN,
		Error: fmt.Sprintf("must be greater or equal to %f", min),
		Data:  Min(min),
	}
}

func InvalidFloatMax(max float64) Invalid {
	return Invalid{
		Code:  utils.VAL_FLOAT_MAX,
		Error: fmt.Sprintf("must be less than or equal to %f", max),
		Data:  Max(max),
	}
}

func InvalidFloatRange(min float64, max float64) Invalid {
	return Invalid{
		Code:  utils.VAL_FLOAT_RANGE,
		Error: fmt.Sprintf("must be between %f and %f", min, max),
		Data:  Range(min, max),
	}
}

func Checkout() *Result {
	return globalPool.Checkout()
}

type Invalid struct {
	Code  uint32 `json:"code"`
	Error string `json:"error"`
	Data  any    `json:"data"` //TODO: add omitempty https://github.com/goccy/go-json/issues/391
}

type InvalidField struct {
	Invalid
	Field string `json:"field"`
}

type DataRange struct {
	Min any `json:"min"`
	Max any `json:"max"`
}

func Range(min any, max any) any {
	return DataRange{
		Min: min,
		Max: max,
	}
}

type DataMin struct {
	Min any `json:"min"`
}

func Min(min any) any {
	return DataMin{
		Min: min,
	}
}

type DataMax struct {
	Max any `json:"max"`
}

func Max(max any) any {
	return DataMax{
		Max: max,
	}
}

type DataChoice[T any] struct {
	Valid []T `json:"valid"`
}

func Choice[T any](valid []T) any {
	return DataChoice[T]{
		Valid: valid,
	}
}
