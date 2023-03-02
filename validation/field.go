package validation

import (
	"strings"
)

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

func (f *Field) nest(parent *Field) *Field {
	path := append(parent.Path, f.Name)

	// If this was an array field (where the last path value was an empty (placeholder))
	// then it should remain so. This is necessary because of the ForceField function
	// which is forces a specific field path without being aware of being forced
	// into an array field.
	if f.Path[len(f.Path)-1] == "" {
		path = append(path, "")
	}

	return &Field{
		// the name doesnt' change, the name is always the outer part
		Name: f.Name,
		Path: path,
		Flat: strings.Join(path, "."),
	}
}

func BuildField(flat string) *Field {
	var name string
	path := strings.Split(flat, ".")
	for i, part := range path {
		if part == "#" {
			path[i] = ""
		} else {
			name = path[i] // the name is the last non-array placeholder
		}
	}

	// 1 Field can be used in 2 ways. It can be used to validate the array itself
	// (like too many entries, or not an array type). And, it can be used to validate
	// data WITHIN the array. When validating the array itself, we'll use field.Flat
	// because that's like validating any other object field, e.g user.favorites.
	// This is static.
	// When we're validating a value inside, we'll use field.Path, e.g. user.favorites.#
	// This is dynamic.
	// So, for the flat name, we never want the trailing placeholder. In fact,
	// placeholders within the flat name NEVER makes sense, because flat == static.
	if path[len(path)-1] == "" {
		flat = strings.Join(path[:len(path)-1], ".")
	}

	return &Field{
		Name: name,
		Path: path,
		Flat: flat,
	}
}

func SimpleField(name string) *Field {
	return &Field{Flat: name}
}
