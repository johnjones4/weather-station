package core

func Pointer[T any](in T) *T {
	return &in
}
