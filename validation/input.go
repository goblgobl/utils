package validation

/*
When we first create a validation rule, we have no field name. For
example:

	countValidation := Int().Required().Min(0)

Only once we add it to an object does it get a field name:

	Object().Field("count", countValidation)

This means that `countValidation` can be used with different field
names (either in the same object or in a different object):

	Object().
		Field("count", countValidation).
		Field("confirm_count", countValidation)

Adding a field to an object (e.g. Field(NAME, validationRule)) is
relatively expensive. We deep-clone the rule.

This logic also applies for nested objects:

	Object().Field("user", Object().
		Field("name", nameValidation).
		Field("age", ageValidation)
	)

After this is called, the code that represents the name validation
will have a field of "user.name". This is done in two steps,
first, the call to `Field("name", nameValidation)` will give it a
field of "name". Then when the "user" field is added all children
as passed the new parent, thus the field becomes "user.name".

This approach is inneficient as we end up creating a lot of copies of
rules and fields, but this should only be happening on startup.
*/

import (
	"strings"

	"github.com/valyala/fasthttp"
	"src.goblgobl.com/utils/typed"
)

type InputValidator interface {
	addField(name string) InputValidator
	validate(object typed.Typed, input typed.Typed, res *Result)
	argsToTyped(args *fasthttp.Args, dest typed.Typed)
}

/*
In normal cases, field names are simple. Once the validator is created
we have a static field name. That could be something like "name" or even
"user.name". However, when our validator is inside of an array, the
field name becomes dynamic as it includes the index: e.g. "users.4.name".
(of course, this could be deeply nested, e.g. "results.users.10.favorites.99.tag").

The Field struct maintains both the static fieldName and the individual
parts that can be used to create a dynamic indexed-containing field name.
When possible, i.e. when we're not within an array, the static name is used.
*/

type Field struct {
	// The name of the actual field that we need to look up. This is
	// always equal to the last element in Path
	Name string

	// strings.Join(Path, ".")
	Flat string

	// The full path of the field.
	// Could be a single value, like: ["name"],
	// Could be nested, like ["user", "name"],
	// Could contain placeholders for array indexex, like: ["users", "", name]
	Path []string
}

func NewField(name string) Field {
	return Field{
		Name: name,
		Flat: name,
		Path: []string{name},
	}
}

func NewIndexedField(name string) Field {
	return Field{
		Name: name,
		Flat: name,
		Path: []string{name, ""}, // empty space for the index placeholder
	}
}

func (f Field) add(name string) Field {
	flat := name
	var path []string

	if current := f.Path; current == nil {
		path = []string{name}
	} else {
		path = append([]string{name}, current...)
		name = f.Name // name is the first name we have, hence here we preserve it
		flat = strings.Join(path, ".")
	}

	return Field{
		Name: name,
		Path: path,
		Flat: flat,
	}
}
