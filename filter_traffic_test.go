package filter_traffic

import "testing"

func BenchmarkRepository_FilterTraffic(b *testing.B) {
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

	r := New[string](FilterTrafficConfig[string]{Enabled: true}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.LetThrough("DE")
			r.LetThrough("UK")
			r.LetThrough("RS")
		}
	})
}
