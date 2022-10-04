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

	if r.disableAfterLimit {
		if r.limitReached.Load() {
			return true
		}

		if r.counter.CompareAndSwap(r.limit, r.limit) {
			r.limitReached.Store(true)
			return true
		}
	}

	old := r.counter.Add(1)

	return uint16(old%100) > r.filterPercentage
}
