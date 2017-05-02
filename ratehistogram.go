package ratehistogram

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/paulbellamy/ratecounter"
)

// RateHistogram has bins and edges. bins are counters.
// For three bins 0-1,1-4,4-5, the edges arrays is {1,4,5}.
// sort.S
type RateHistogram struct {
	bins        []*ratecounter.RateCounter
	edges       []float64
	highestEdge float64
}

// NewRateHistogram creates RateHistogram and inits counters
func NewRateHistogram(edges []float64, interval time.Duration) (*RateHistogram, error) {
	bins := make([]*ratecounter.RateCounter, len(edges))
	for i := range edges {
		bins[i] = ratecounter.NewRateCounter(interval)
	}
	highestEdge := edges[len(edges)-1]
	if !sort.Float64sAreSorted(edges) {
		return nil, errors.New("edges array must be sorted")
	}
	return &RateHistogram{bins: bins, edges: edges, highestEdge: highestEdge}, nil
}

// Record records value to proper bucket
func (h *RateHistogram) Record(v float64) {
	if v > h.highestEdge {
		return
	}
	h.bins[sort.SearchFloat64s(h.edges, v)].Incr(1)
}

// ObservePretty returns map of upper_edge => bin_count
func (h *RateHistogram) ObservePretty() map[string]int64 {
	ret := make(map[string]int64)
	for i, b := range h.bins {
		ret[fmt.Sprintf("%.2f", h.edges[i])] = b.Rate()
	}
	return ret
}

// Observe returns only current bin counts
func (h *RateHistogram) Observe() []int64 {
	ret := make([]int64, len(h.bins))
	for i, b := range h.bins {
		ret[i] = b.Rate()
	}
	return ret

}
