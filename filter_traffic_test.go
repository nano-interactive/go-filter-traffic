package filter_traffic

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterTraffic_SkipEverything_GlobalTraffic_Is_Zero(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	globalFilter := GlobalFilter[string]{
		Limit: 0,
		Counter: &Counter{
			ResetNumber: 100,
		},
	}

	perValueFilter := PerValueFilterMap[string]{}

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	assert.Equal(false, r.Do("Test"))
}

func TestFilterTraffic_PerValueFiller_Is_Zero(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	globalFilter := GlobalFilter[string]{
		Limit: 2,
		Counter: &Counter{
			ResetNumber: 100,
		},
	}

	perValueFilter := NewPerValueFilterMap(100, map[string]uint64{"DE": 0})

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	assert.Equal(false, r.Do("Test")) // Fails as 'Test' is not found
	assert.Equal(false, r.Do("DE"))   // Fails as 'DE' is found but limit is 0
}

func TestFilterTraffic_GlobalReset(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	counter := &Counter{
		ResetNumber: 100,
		counter:     atomic.Uint64{},
	}

	counter.counter.Store(100)

	globalFilter := GlobalFilter[string]{
		Limit:   2,
		Counter: counter,
	}

	perValueFilter := NewPerValueFilterMap(100, map[string]uint64{"DE": 100})

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	assert.Equal(false, r.Do("Test")) // Fails as 'Test' is not found
	assert.Equal(true, r.Do("DE"))    // Passes as 'DE' is found and limit is 100
}

func TestFilterTraffic_PerValue_Reset(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	counter := &Counter{
		ResetNumber: 100,
		counter:     atomic.Uint64{},
	}

	globalFilter := GlobalFilter[string]{
		Limit:   10,
		Counter: counter,
	}

	perValueFilter := NewPerValueFilterMap(100, map[string]uint64{"DE": 1})
	perValueFilter.counter["DE"].counter.Store(100)

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	assert.Equal(false, r.Do("Test")) // Fails as 'Test' is not found
	assert.Equal(true, r.Do("DE"))    // Passes as 'DE' is found and limit is 1 and counter is reset
}

func BenchmarkRepository_FilterTraffic_Disabled(b *testing.B) {
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

	r := New(FilterTrafficConfig{Enabled: false}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Do("DE")
			r.Do("UK")
			r.Do("RS")
		}
	})
}

func BenchmarkRepository_FilterTraffic_WithMaps_AllFound(b *testing.B) {
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

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Do("DE")
			r.Do("UK")
		}
	})
}

func BenchmarkRepository_FilterTraffic_WithMaps_AllNotFoundFound(b *testing.B) {
	globalFilter := GlobalFilter[string]{
		Limit: 50,
		Counter: &Counter{
			ResetNumber: 100,
		},
	}

	perValueFilter := NewPerValueFilterMap(100, map[string]uint64{
		"UK": 10,
	})

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Do("DE")
			r.Do("RS")
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
		UK:        10,
		DE:        50,
		uKCounter: &Counter{ResetNumber: 100},
		dECounter: &Counter{ResetNumber: 100},
	}

	r := New(FilterTrafficConfig{Enabled: true}, globalFilter, perValueFilter)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Do("DE")
			r.Do("UK")
			r.Do("RS")
		}
	})
}
