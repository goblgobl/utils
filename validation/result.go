package validation

import (
	"strconv"
	"strings"
)

type Result struct {
	len          uint64
	errors       []any
	pool         *Pool
	arrayIndexes []int
	arrayCount   int
}

func NewResult(maxErrors uint16) *Result {
	return &Result{
		arrayCount:   -1,
		errors:       make([]any, maxErrors),
		arrayIndexes: make([]int, 10),
	}
}

func (r Result) Errors() []any {
	return r.errors[:r.len]
}

func (r Result) IsValid() bool {
	return r.len == 0
}

func (r Result) Len() uint64 {
	return r.len
}
func (r *Result) AddInvalidField(field Field, invalid Invalid) {
	fieldName := field.Flat

	if r.arrayCount != -1 {
		// We're inside of an array, we need to create field name dynamically
		// TODO: optimize this code
		var w strings.Builder

		// Over allocate a little so that we likely won't have to allocate + copy.
		w.Grow(len(field.Name) + 20)

		indexIndex := 0
		indexes := r.arrayIndexes
		for _, part := range field.Path {
			w.WriteByte('.')
			if part == "" {
				w.WriteString(strconv.Itoa(indexes[indexIndex]))
				indexIndex += 1
			} else {
				w.WriteString(part)
			}
		}

		// [1:] to strip out the leading .
		fieldName = w.String()[1:]
	}

	r.Add(InvalidField{
		Field:   fieldName,
		Invalid: invalid,
	})
}

func (r *Result) Add(error any) {
	l := r.len
	errors := r.errors
	if l < uint64(len(errors)) {
		errors[l] = error
		r.len = l + 1
	}
}

func (r *Result) BeginArray() {
	r.arrayCount += 1
}

func (r *Result) ArrayIndex(i int) {
	r.arrayIndexes[r.arrayCount] = i
}

func (r *Result) EndArray() {
	r.arrayCount -= 1
}

func (r *Result) Release() {
	if pool := r.pool; pool != nil {
		r.len = 0
		r.arrayCount = -1
		pool.list <- r
	}
}
