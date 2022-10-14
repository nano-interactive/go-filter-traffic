package filter_traffic

type Filter[T any] interface {
	Do(T) bool
}
