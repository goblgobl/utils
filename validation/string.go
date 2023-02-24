package validation

import (
	"fmt"
	"regexp"

	"src.goblgobl.com/utils"
)

type StringTransform func(value string) string
type StringFuncValidator[T any] func(value string, ctx *Context[T]) any

type StringValidator[T any] struct {
	invalidLength  *Invalid
	invalidChoice  *Invalid
	invalidPattern *Invalid
	pattern        *regexp.Regexp
	fn             StringFuncValidator[T]
	tx             StringTransform
	dflt           any
	choices        []string
	minLength      int
	maxLength      int
	required       bool
	nullable       bool
}

func String[T any]() *StringValidator[T] {
	return new(StringValidator[T])
}

func (v *StringValidator[T]) Validate(raw any, ctx *Context[T]) any {
	if raw == nil {
		if dflt := v.dflt; dflt != nil {
			return dflt
		}
		if v.nullable {
			return nil
		}
		if v.required {
			ctx.InvalidField(Required)
		}
		return nil
	}

	value, ok := raw.(string)
	if !ok {
		ctx.InvalidField(TypeString)
		return nil
	}

	if tx := v.tx; tx != nil {
		value = tx(value)
	}

	if min := v.minLength; len(value) < min {
		ctx.InvalidField(v.invalidLength)
		return value
	}
	if max := v.maxLength; max > 0 && len(value) > max {
		ctx.InvalidField(v.invalidLength)
		return value
	}

	if pattern := v.pattern; pattern != nil && !pattern.MatchString(value) {
		ctx.InvalidField(v.invalidPattern)
		return value
	}

	if choices := v.choices; choices != nil {
		validChoice := false
		for _, valid := range choices {
			if value == valid {
				validChoice = true
				break
			}
		}
		if !validChoice {
			ctx.InvalidField(v.invalidChoice)
			return value
		}
	}

	if fn := v.fn; fn != nil {
		return fn(value, ctx)
	}
	return value
}

func (v *StringValidator[T]) Required() *StringValidator[T] {
	v.required = true
	return v
}

// Meant to be used in conjunction with Clone(). Maybe on create, the field is
// required, but on update, it isn't.
func (v *StringValidator[T]) NotRequired() *StringValidator[T] {
	v.required = false
	return v
}

func (v *StringValidator[T]) Nullable() *StringValidator[T] {
	v.nullable = true
	return v
}

func (v *StringValidator[T]) Default(dflt any) *StringValidator[T] {
	v.dflt = dflt
	return v
}

func (v *StringValidator[T]) Min(min int) *StringValidator[T] {
	v.minLength = min
	v.invalidLength = InvalidStringLength(min, v.maxLength)
	return v
}

func (v *StringValidator[T]) Max(max int) *StringValidator[T] {
	v.maxLength = max
	v.invalidLength = InvalidStringLength(v.minLength, max)
	return v
}

func (v *StringValidator[T]) Length(min int, max int) *StringValidator[T] {
	v.minLength = min
	v.maxLength = max
	v.invalidLength = InvalidStringLength(min, max)
	return v
}

func (v *StringValidator[T]) Pattern(pattern string, message ...string) *StringValidator[T] {
	v.pattern = regexp.MustCompile(pattern)
	v.invalidPattern = InvalidStringPattern(message...)
	return v
}

func (v *StringValidator[T]) Choice(choices ...string) *StringValidator[T] {
	v.choices = choices
	v.invalidChoice = InvalidStringChoice(choices)
	return v
}

func (v *StringValidator[T]) Func(fn StringFuncValidator[T]) *StringValidator[T] {
	v.fn = fn
	return v
}

func (v *StringValidator[T]) Transform(tx StringTransform) *StringValidator[T] {
	v.tx = tx
	return v
}

func (v *StringValidator[T]) Clone() *StringValidator[T] {
	return &StringValidator[T]{
		dflt:           v.dflt,
		required:       v.required,
		nullable:       v.nullable,
		minLength:      v.minLength,
		maxLength:      v.maxLength,
		pattern:        v.pattern,
		choices:        v.choices,
		tx:             v.tx,
		fn:             v.fn,
		invalidLength:  v.invalidLength,
		invalidChoice:  v.invalidChoice,
		invalidPattern: v.invalidPattern,
	}
}

func InvalidStringLength(min int, max int) *Invalid {
	hasMin := min != 0
	hasMax := max != 0

	if !hasMin && !hasMax {
		return nil
	}

	if hasMin && hasMax {
		return &Invalid{
			Code:  utils.VAL_STRING_LEN,
			Error: fmt.Sprintf("must be between %d and %d characters", min, max),
			Data:  RangeData(min, max),
		}
	}

	if hasMin {
		c := "characters"
		if min == 1 {
			c = "character"
		}
		return &Invalid{
			Code:  utils.VAL_STRING_LEN,
			Error: fmt.Sprintf("must be atleast %d %s", min, c),
			Data:  MinData(min),
		}
	}

	c := "characters"
	if max == 1 {
		c = "character"
	}
	return &Invalid{
		Code:  utils.VAL_STRING_LEN,
		Error: fmt.Sprintf("must be no more than %d %s", max, c),
		Data:  MaxData(max),
	}
}

func InvalidStringChoice(choices []string) *Invalid {
	return &Invalid{
		Code:  utils.VAL_STRING_CHOICE,
		Error: "is not a valid choice",
		Data:  ChoiceData(choices),
	}
}
