package filter_traffic

import "sync/atomic"

type (
	PercentageFilterConfig struct {
		Enabled           bool
		DisableAfterLimit bool
		Limit             uint64
		FilterPercentage  uint16
	}

	PercentageFilter[T any] struct {
		counter           atomic.Uint64
		limit             uint64
		limitReached      atomic.Bool
		disableAfterLimit bool
		enabled           bool
		filterPercentage  uint16
	}
)

var _ Filter[any] = &PercentageFilter[any]{}

func NewPercentageFilter[T any](cfg PercentageFilterConfig) *PercentageFilter[T] {
	if cfg.FilterPercentage > 100 {
		cfg.FilterPercentage = 100
	}

	return &PercentageFilter[T]{
		enabled:           cfg.Enabled,
		counter:           atomic.Uint64{},
		limitReached:      atomic.Bool{},
		limit:             cfg.Limit,
		disableAfterLimit: cfg.DisableAfterLimit,
		filterPercentage:  100 - cfg.FilterPercentage,
	}
}

func (r *PercentageFilter[T]) Do(key T) bool {
	if !r.enabled {
		return true
	}

	// Disable after limit is a feature that disables the filter after the limit is reached.
	// This is useful for cases where the limit is reached and the filter is not needed anymore.
	// e.g. database filter -> using filter to limit the number of requests to the database
	// after the certain number of requests have been cached, database will not be bothered anymore
	// as much as it was when the service started. This avoids the database from being overloaded.
	// In some cases 'disable after limit' feature is not useful so it can be optimized by the compiler
	// if the value is known at the compile time -> referred as 'Dead Code Elimination (DCE)'
	if r.disableAfterLimit {
		if r.limitReached.Load() {
			return true
		}

		// TODO: Check if compare and swap is faster than load with the same compare
		// if r.counter.Load() >= r.limit {
		// 	r.limitReached.Store(true)
		// 	return true
		// }
		if r.counter.CompareAndSwap(r.limit, r.limit) {
			r.limitReached.Store(true)
			return true
		}
	}

	old := r.counter.Add(1)

	// TODO: Check if compareAndSwap with reset to 100 is faster than division on atomic values
	// wrap around algorithm
	return uint16(old%100) > r.filterPercentage
}
