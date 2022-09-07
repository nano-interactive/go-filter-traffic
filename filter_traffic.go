package filter_traffic

type (
	Filter interface {
		LetThrough(country string) bool
	}

	FilterTrafficConfig[T comparable] struct {
		Enabled bool
	}

	FilterTraffic[T comparable, TFilter PerValueFilter[T]] struct {
		globalFilter GlobalFilter[T]
		filter       TFilter
		enabled      bool
	}
)

var _ Filter = FilterTraffic[string, GlobalFilter[string]]{}

func New[T comparable, TFilter PerValueFilter[T]](config FilterTrafficConfig[T], globalFilter GlobalFilter[T], other TFilter) FilterTraffic[T, TFilter] {
	return FilterTraffic[T, TFilter]{
		enabled:      config.Enabled,
		globalFilter: globalFilter,
		filter:       other,
	}
}

func (r FilterTraffic[T, TFilter]) LetThrough(key T) bool {
	if !r.enabled {
		return true
	}

	counter := r.globalFilter.Counter

	if counter.counter.CompareAndSwap(100, 0) {
		counter.counter.Add(1)
		return true
	}

	limit := r.globalFilter.Limit
	old := r.globalFilter.Counter.counter.Add(1)

	if old < limit {
		return true
	}

	if !r.filter.HasKey(key) {
		return false
	}

	counter = r.filter.GetCounter(key)

	if counter.counter.CompareAndSwap(100, 0) {
		counter.counter.Add(1)
		return true
	}

	limit = r.filter.GetLimit(key)
	old = counter.counter.Add(1)

	return old < limit
}
