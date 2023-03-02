package validation

import (
	"strconv"
	"strings"

	"src.goblgobl.com/utils/typed"
)

type Context[T any] struct {
	release func(*Context[T])

	// For errors within arrays, we need to create the field name dynamically
	// (e.g. users.#.name). We'll use this scrap space
	scrap []byte

	// The first currently being validated
	Field *Field

	// the full user input
	Input typed.Typed

	// app-specific data that we want to make available to validation callbacks
	Env T

	// stack of array indexes as we go deeper and deeper into nesting
	arrayIndexes []int

	objects []typed.Typed

	errors []any

	// errors is a pre-allocated array of len(maxErrors).
	// errLen keeps track of how many erros we've already added so that we can
	// use errors[errLen] = invalid when adding an error
	// errLen is reset to 0 when this context is released (and thus can be reused)
	errLen int

	// how deeply nested we are (only cares about arrays)
	arrayDepth  int
	objectDepth int
}

func NewContext[T any](maxErrors uint16) *Context[T] {
	return &Context[T]{
		arrayDepth:   -1,
		objectDepth:  -1,
		arrayIndexes: make([]int, 10),
		objects:      make([]typed.Typed, 10),
		errors:       make([]any, maxErrors),
	}
}

func (c *Context[T]) IsValid() bool {
	return c.errLen == 0
}

func (c *Context[T]) ErrorCount() int {
	return c.errLen
}

func (c *Context[T]) Errors() []any {
	return c.errors[:c.errLen]
}

func (c *Context[T]) Objects() []typed.Typed {
	d := c.objectDepth
	if d == -1 {
		return nil
	}
	return c.objects[:d+1]
}

func (c *Context[T]) CurrentObject() typed.Typed {
	d := c.objectDepth
	if d == -1 {
		return nil
	}
	return c.objects[d]
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
// StartArray adds a new arrayDepth to the stack
func (c *Context[T]) StartArray() {
	c.arrayDepth += 1
}

// Sets the index array for the current array
func (c *Context[T]) ArrayIndex(i int) {
	c.arrayIndexes[c.arrayDepth] = i
}

// Removes an array from the stack
func (c *Context[T]) EndArray() {
	c.arrayDepth -= 1
}

func (c *Context[T]) StartObject(o typed.Typed) {
	depth := c.objectDepth + 1
	c.objects[depth] = o
	c.objectDepth = depth
}

func (c *Context[T]) EndObject() {
	c.objectDepth -= 1
}

// This is an advanced and very hacky thing. It's generally a "good thing" that
// this context is stateful. But in some cases, a caller might need more flexibility.
// This came up because within the process of validating an array, we needed
// to validate some other data which would not inherit the array index. This
// call erases the array state of the context. It returns this state (which is
// just the arrayDepth) to the caller, so that the caller can pass it back to ResumeArray
// when it's time to go back to the normal valiation
func (c *Context[T]) SuspendArray() int {
	arrayDepth := c.arrayDepth
	c.arrayDepth = -1
	return arrayDepth
}

// See SuspendArray
func (c *Context[T]) ResumeArray(arrayDepth int) {
	c.arrayDepth = arrayDepth
}

func (c *Context[T]) InvalidField(invalid *Invalid) {
	c.InvalidWithField(invalid, c.Field)
}

func (c *Context[T]) InvalidWithField(invalid *Invalid, field *Field) {
	arrayDepth := c.arrayDepth
	if arrayDepth == -1 {
		// we're not in an array, so the field name isn't dynamic
		c.addInvalid(InvalidField{Field: field.Flat, Invalid: invalid})
		return
	}

	// We're in an array (possibly deeply nested), so the field name
	// depends on the array index (thus it's dynamic)
	var w strings.Builder

	// This memory has to live as long as the context, but within 1 context, we
	// might need multiple (one for each error to a dynamic field). I'm struggling
	// to find a good way to avoid this allocation.
	// Over-allocate a little so that we likely won't have to allocate + copy.
	w.Grow(len(field.Flat) + 15)

	indexIndex := 0
	indexes := c.arrayIndexes
	for _, part := range field.Path {
		if part == "" {
			if indexIndex > arrayDepth {
				// The Path contains the entire label, say ["entries", "", "names", ""]
				// But sometimes the error is on the array itself, say "entries.3.names".
				// So we only want to generate this label based on the current arrayDepth
				// of the context
				break
			}
			index := indexes[indexIndex]
			indexIndex += 1
			w.WriteByte('.')
			w.WriteString(strconv.Itoa(index))
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
	var noEnv T
	c.Env = noEnv
	c.errLen = 0
	c.arrayDepth = -1
	c.objectDepth = -1
	c.release(c)
}
