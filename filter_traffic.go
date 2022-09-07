package filter_traffic

type (
	Filter interface {
		LetThrough(country string) bool
	}

	FilterTrafficConfig[T comparable] struct {
		Enabled bool
	}

	FilterTraffic[T comparable, TFilter PerValueFilter[T]] struct {
		enabled      bool
		globalFilter GlobalFilter[T]
		filter       TFilter
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

	if r.globalFilter.Reset(key) {
		r.globalFilter.Increment(key)
		return true
	}

	limit := r.globalFilter.GetLimit(key)
	old := r.globalFilter.Increment(key)

	if old < limit {
		return true
	}

	if !r.filter.HasKey(key) {
		return false
	}

	if r.filter.Reset(key) {
		r.filter.Increment(key)
		return true
	}

	limit = r.filter.GetLimit(key)
	old = r.filter.Increment(key)

	return old < limit
}
