package validation

import (
	"strconv"
	"strings"

	"src.goblgobl.com/utils/typed"
)

type Context[T any] struct {
	pool *Pool[T]

	// The first currently being validated
	Field *Field

	// the full user input
	Input typed.Typed

	// The current object being looked at (this is different than Input for arrays
	// and nested objects)
	Object typed.Typed

	// app-specific data that we want to make available to validation callbacks
	Env T

	// stack of array indexes as we go deeper and deeper into nesting
	arrayIndexes []int

	errors []any

	// errors is a pre-allocated array of len(maxErrors).
	// errLen keeps track of how many erros we've already added so that we can
	// use errors[errLen] = invalid when adding an error
	// errLen is reset to 0 when this context is released (and thus can be reused)
	errLen int

	// how deeply nested we are (only cares about arrays)
	depth int
}

func NewContext[T any](maxErrors uint16) *Context[T] {
	return &Context[T]{
		depth:        -1,
		arrayIndexes: make([]int, 10),
		errors:       make([]any, maxErrors),
	}
}

func (c *Context[T]) IsValid() bool {
	return c.errLen == 0
}

func (c *Context[T]) Errors() []any {
	return c.errors[:c.errLen]
}

// Directly executes a validator with the given value and for the given field.
// This is generally called when validation of fields is dynamic, often in a
// a Func callbac(like Object().Func) where the fields to validate are dependent
// on other parts of the data
func (c *Context[T]) Validate(field *Field, value any, validator Validator[T]) (any, bool) {
	c.Field = field
	l := c.errLen
	value = validator.Validate(value, c)
	return value, c.errLen == l
}

// Arrays can be deeply nested. We keep a stack of values indicating what array
// index we're currently at, at each level.
// StartArray adds a new depth to the stack
func (r *Context[T]) StartArray() {
	r.depth += 1
}

// Sets the index array for the current array
func (r *Context[T]) ArrayIndex(i int) {
	r.arrayIndexes[r.depth] = i
}

// Removes an array from the stack
func (r *Context[T]) EndArray() {
	r.depth -= 1
}

func (c *Context[T]) InvalidField(invalid *Invalid) {
	c.InvalidWithField(invalid, c.Field)
}

func (c *Context[T]) InvalidWithField(invalid *Invalid, field *Field) {
	depth := c.depth
	if depth == -1 {
		// we're not in an array, so the field name isn't dynamic
		c.addInvalid(InvalidField{Field: field.Flat, Invalid: invalid})
		return
	}

	// we're in an array (possibly deeply nested), to the field name
	// depends on the array index (thus it's dynamic)
	// TODO: optimize this code
	var w strings.Builder

	// Over allocate a little so that we likely won't have to allocate + copy.
	w.Grow(len(field.Name) + 20)

	indexIndex := 0
	indexes := c.arrayIndexes
	for _, part := range field.Path {
		if part == "" {
			if indexIndex > depth {
				// The Path contains the entire label, say ["entries", "", "names", ""]
				// But sometimes the error is on the array itself, say "entries.3.names".
				// So we only want to generate this label based on the current depth
				// of the context
				break
			}
			w.WriteByte('.')
			index := indexes[indexIndex]
			w.WriteString(strconv.Itoa(index))
			indexIndex += 1
		} else {
			w.WriteByte('.')
			w.WriteString(part)
		}
	}

	// [1:] to strip out the leading .
	c.addInvalid(InvalidField{Field: w.String()[1:], Invalid: invalid})
}

func (c *Context[T]) Invalid(invalid *Invalid) {
	c.addInvalid(invalid)
}

func (c *Context[T]) addInvalid(error any) {
	l := c.errLen
	errors := c.errors

	// add up to MAX allowed errors
	if l < len(errors) {
		errors[l] = error
		c.errLen = l + 1
	}
}

func (c *Context[T]) Release() {
	if pool := c.pool; pool != nil {
		var noEnv T
		c.Env = noEnv
		c.depth = -1
		c.errLen = 0
		c.pool.list <- c
	}
}
