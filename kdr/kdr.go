package kdr

// Tri-state for either keeping a value as is, deleting the value, or replacing it

type KdrAction int

const (
	KDR_ACTION_KEEP = iota
	KDR_ACTION_DELETE
	KDR_ACTION_REPLACE
)

type Value[T any] struct {
	Replacement T
	Action      KdrAction
}

func (v Value[T]) IsKeep() bool {
	return v.Action == KDR_ACTION_KEEP
}

func (v Value[T]) IsDelete() bool {
	return v.Action == KDR_ACTION_DELETE
}

func (v Value[T]) IsReplace() bool {
	return v.Action == KDR_ACTION_REPLACE
}

func Keep[T any]() Value[T] {
	return Value[T]{Action: KDR_ACTION_KEEP}
}

func Delete[T any]() Value[T] {
	return Value[T]{Action: KDR_ACTION_DELETE}
}

func Replace[T any](replacement T) Value[T] {
	return Value[T]{
		Action:      KDR_ACTION_REPLACE,
		Replacement: replacement,
	}
}
