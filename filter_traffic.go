package filter_traffic

type (
	Filter interface {
		LetThrough(country string) bool
	}

	FilterTrafficConfig[T comparable] struct {
		Enabled bool
	}

	FilterTraffic[T comparable] struct {
		enabled bool
		filters []PerValueFilter[T]
	}
)

var _ Filter = FilterTraffic[string]{}

func New[T comparable](config FilterTrafficConfig[T], filters ...PerValueFilter[T]) FilterTraffic[T] {
	return FilterTraffic[T]{
		enabled: config.Enabled,
		filters: filters,
	}
}

func (r FilterTraffic[T]) LetThrough(key T) bool {
	if !r.enabled {
		return true
	}

	for _, filter := range r.filters {
		if !filter.HasKey(key) {
			continue
		}

		if filter.Reset(key) {
			filter.Increment(key)
			return true
		}

		limit := filter.GetLimit(key)
		old := filter.Increment(key)

		if old < limit {
			return true
		}
	}

	return false
}
