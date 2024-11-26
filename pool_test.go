package main

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

func TestStringSlicePool(t *testing.T) {
	p := newStringSlicePool()

	for range 100 {
		s := p.get()

		if len(s) > 0 {
			t.Errorf("get() returned non-empty slice: %q", s)
		}

		count := rand.IntN(100)

		for i := range count {
			s = append(s, strconv.Itoa(i))
		}

		switch rand.IntN(3) {
		case 1:
			s = s[:0]
		case 2:
			s = s[:count/2]
		}

		p.put(s)

		if idx := slices.IndexFunc(s[:count], func(value string) bool {
			return value != ""
		}); idx != -1 {
			t.Errorf("Slice was not cleared: %q", s[:count])
		}
	}
}
