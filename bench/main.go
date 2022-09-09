package main

import traffic "github.com/nano-interactive/go-filter-traffic"

func main() {
	globalFilter := traffic.GlobalFilter[string]{
		Limit: 50,
		Counter: &traffic.Counter{
			ResetNumber: 100,
		},
	}

	perValueFilter := traffic.NewPerValueFilterMap(100, map[string]uint64{
		"UK": 10,
		"DE": 50,
	})

	r := traffic.New(traffic.FilterTrafficConfig[string]{Enabled: true}, globalFilter, perValueFilter)

	for i := 0; i < 1000_0000; i++ {
		r.Do("DE")
		r.Do("UK")
		r.Do("RS")
	}
}
