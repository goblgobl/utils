package log

import "fmt"

// An error that's designed to be logged in a more structured manner
type StructuredError struct {
	Code int            `json:"code"`
	Err  error          `json:"err"`
	Data map[string]any `json:"data"`
}

func (e StructuredError) Error() string {
	return fmt.Sprintf("code: %d - %s", e.Code, e.Err.Error())
}

func (e *StructuredError) Int(key string, value int) *StructuredError {
	e.ensureMap()
	e.Data[key] = value
	return e
}

func (e *StructuredError) String(key string, value string) *StructuredError {
	e.ensureMap()
	e.Data[key] = value
	return e
}

func (e *StructuredError) ensureMap() {
	if e.Data == nil {
		e.Data = make(map[string]any, 1)
	}
}

func Err(code int, err error) *StructuredError {
	return &StructuredError{
		Err:  err,
		Code: code,
	}
}

func Errf(code int, format string, args ...any) *StructuredError {
	return Err(code, fmt.Errorf(format, args...))
}
