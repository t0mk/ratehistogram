package ratehistogram

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/paulbellamy/ratecounter"
	"gopkg.in/yaml.v2"
)

// RateHistogram has bins and edges. bins are counters.
// For three bins 0-1,1-4,4-5, the edges arrays is {1,4,5}.
// sort.S
type RateHistogram struct {
	bins        []*ratecounter.RateCounter
	edges       []float64
	highestEdge float64
	secs        int
}

// NamedArray is a helper type for e.g. serialization
type NamedArray struct {
	Name string
	Bins []int64
}

// HMap is just a map of histograms, for convenience
type HMap map[string]*RateHistogram

// Conf is prescription for rate histogram
type Conf struct {
	Edges []float64 `yaml: "edges"`
	Secs  int       `yaml: "secs"`
}

// NewHMapFromYAML creates histomap from YAML
func NewHMapFromYAML(data []byte) (HMap, error) {
	conf := make(map[string]Conf)
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return NewHMap(conf)
}

// NewHMap creates map of rate histograms
func NewHMap(conf map[string]Conf) (HMap, error) {
	ret := make(HMap)
	for k, v := range conf {

		nh, err := NewRateHistogram(v)
		if err != nil {
			return nil, err
		}
		ret[k] = nh
	}
	return ret, nil
}

// Observe returns names of histos plus counts
func (hmap *HMap) Observe() map[string][]int64 {
	ret := make(map[string][]int64)
	for k, v := range *hmap {
		ret[k] = v.Observe()
	}
	return ret

}

// GetSnapshot returns snapshot for a particular hits
func (hmap *HMap) GetSnapshot(hname string) *NamedArray {
	return &NamedArray{Name: hname, Bins: (*hmap)[hname].Observe()}

}

// GetNames returns the names of recorded histograms
func (hmap *HMap) GetNames() []string {
	ret := make([]string, 0)
	for k := range *hmap {
		ret = append(ret, k)
	}
	return ret
}

// NewRateHistogram creates RateHistogram and inits counters
func NewRateHistogram(conf Conf) (*RateHistogram, error) {
	edges := conf.Edges
	interval := time.Second * time.Duration(conf.Secs)
	bins := make([]*ratecounter.RateCounter, len(edges))
	for i := range edges {
		bins[i] = ratecounter.NewRateCounter(interval)
	}
	highestEdge := edges[len(edges)-1]
	if !sort.Float64sAreSorted(edges) {
		return nil, errors.New("edges array must be sorted")
	}
	return &RateHistogram{bins: bins, edges: edges, highestEdge: highestEdge, secs: conf.Secs}, nil
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
