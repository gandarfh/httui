package topointer

func New[T any](v T) *T {
	return &v
}
