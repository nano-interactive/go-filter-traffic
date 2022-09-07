package filter_traffic

import "sync/atomic"

type (
	PerValueFilter[T any] interface {
		HasKey(T) bool
		GetLimit(T) uint64
		GetCounter(T) *Counter
	}

	Counter struct {
		ResetNumber uint64
		counter     atomic.Uint64
	}

	GlobalFilter[T any] struct {
		Limit   uint64
		Counter *Counter
	}

	PerValueFilterMap[T comparable] struct {
		limits  map[T]uint64
		counter map[T]*Counter
	}
)

var (
	_ PerValueFilter[string] = PerValueFilterMap[string]{}
	_ PerValueFilter[string] = GlobalFilter[string]{}
)


func (p GlobalFilter[T]) GetLimit(T) uint64 {
	return p.Limit
}


func (p GlobalFilter[T]) HasKey(T) bool {
	return true
}

func (p GlobalFilter[T]) GetCounter(key T) *Counter {
	return p.Counter
}

func NewPerValueFilterMap[T comparable](max uint64, limits map[T]uint64) PerValueFilterMap[T] {
	counter := make(map[T]*Counter, len(limits))

	for key := range limits {
		counter[key] = &Counter{ResetNumber: max}
	}

	return PerValueFilterMap[T]{
		limits,
		counter,
	}
}

func (p PerValueFilterMap[T]) GetLimit(key T) uint64 {
	return p.limits[key]
}

func (p PerValueFilterMap[T]) HasKey(key T) bool {
	_, ok := p.counter[key]
	return ok
}

func (p PerValueFilterMap[T]) GetCounter(key T) *Counter {
	return p.counter[key]
}
