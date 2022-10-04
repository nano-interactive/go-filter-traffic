package filter_traffic

type (
	CounterWithResetFilterConfig struct {
		Enabled bool
	}

	CounterWithResetFilter[T comparable, TFilter PerValueFilter[T]] struct {
		globalFilter GlobalFilter[T]
		filter       TFilter
		enabled      bool
	}
)

var _ Filter[string] = CounterWithResetFilter[string, GlobalFilter[string]]{}

func NewCounterFilter[T comparable, TFilter PerValueFilter[T]](config PercentageFilterConfig, globalFilter GlobalFilter[T], other TFilter) CounterWithResetFilter[T, TFilter] {
	return CounterWithResetFilter[T, TFilter]{
		enabled:      config.Enabled,
		globalFilter: globalFilter,
		filter:       other,
	}
}

func (r CounterWithResetFilter[T, TFilter]) Do(key T) bool {
	if !r.enabled {
		return true
	}

	counter := r.globalFilter.Counter

	if counter.counter.CompareAndSwap(counter.ResetNumber, 1) {
		// This goto is used to avoid code duplication
		// as we need to check if other filter rules should apply
		// maybe to reset it or to check if it is still valid param
		// passed in the argument ket<T>
		goto checkOtherFilter
	}

	if counter.counter.Add(1) > r.globalFilter.Limit {
		return false
	}

checkOtherFilter:
	counter = r.filter.GetCounter(key)

	if counter == nil {
		return false
	}

	if counter.counter.CompareAndSwap(counter.ResetNumber, 1) {
		return true
	}

	old := counter.counter.Add(1)

	return old < r.filter.GetLimit(key)
}
