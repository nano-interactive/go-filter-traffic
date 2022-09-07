package filter_traffic

import "testing"

func BenchmarkRepository_FilterTraffic_WithMaps(b *testing.B) {
	globalFilter := GlobalFilter[string]{
		Limit: 50,
		Counter: &Counter{
			ResetNumber: 100,
		},
	}

	perValueFilter := NewPerValueFilterMap(100, map[string]uint64{
		"UK": 10,
		"DE": 50,
	})

	r := New(FilterTrafficConfig[string]{Enabled: true}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.LetThrough("DE")
			r.LetThrough("UK")
			r.LetThrough("RS")
		}
	})
}

type TestFilter struct {
	UK        uint64
	DE        uint64
	uKCounter *Counter
	dECounter *Counter
}

func (t TestFilter) HasKey(key string) bool {
	switch key {
	case "UK":
		return true
	case "DE":
		return true
	default:
		return false
	}
}

func (t TestFilter) GetCounter(key string) *Counter {
	switch key {
	case "UK":
		return t.uKCounter
	case "DE":
		return t.dECounter
	default:
		return nil
	}
}

func (t TestFilter) GetLimit(key string) uint64 {
	switch key {
	case "UK":
		return t.UK
	case "DE":
		return t.DE
	default:
		return 0
	}
}

func BenchmarkRepository_FilterTraffic_WithStruct(b *testing.B) {
	globalFilter := GlobalFilter[string]{
		Limit: 50,
		Counter: &Counter{
			ResetNumber: 100,
		},
	}

	perValueFilter := TestFilter{
		UK: 10,
		DE: 50,
		uKCounter: &Counter{ResetNumber: 100},
		dECounter: &Counter{ResetNumber: 100},
	}

	r := New(FilterTrafficConfig[string]{Enabled: true}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.LetThrough("DE")
			r.LetThrough("UK")
			r.LetThrough("RS")
		}
	})
}
