package histogram

import (
	"fmt"
	"math"
	"sort"
)

// Bin represents a single bin in a histogram.
type Bin struct {
	LowerBound float64
	UpperBound float64
	Count      int
}

// Histogram holds the bins and the total number of observations.
type Histogram struct {
	Bins  []Bin
	Total int
}

// New creates a histogram with the given number of bins covering [min, max].
// The bins are of equal width. If min == max, a single bin covering that value is created.
// It panics if bins < 1 or min > max.
func New(bins int, min, max float64) *Histogram {
	if bins < 1 {
		panic("histogram: number of bins must be at least 1")
	}
	if min > max {
		panic("histogram: min must be less than or equal to max")
	}
	if min == max {
		return &Histogram{
			Bins: []Bin{{LowerBound: min, UpperBound: max}},
		}
	}
	width := (max - min) / float64(bins)
	h := &Histogram{
		Bins: make([]Bin, bins),
	}
	for i := 0; i < bins; i++ {
		lower := min + float64(i)*width
		upper := min + float64(i+1)*width
		h.Bins[i] = Bin{LowerBound: lower, UpperBound: upper}
	}
	return h
}

// Add records a value into the histogram.
// Values outside the range [min, max] are ignored.
// The upper bound of the last bin is inclusive; all other bins are half-open [lower, upper).
func (h *Histogram) Add(value float64) {
	if len(h.Bins) == 0 {
		return
	}
	min := h.Bins[0].LowerBound
	max := h.Bins[len(h.Bins)-1].UpperBound
	if value < min || value > max {
		return
	}
	if value == max {
		h.Bins[len(h.Bins)-1].Count++
		h.Total++
		return
	}
	// Find the bin where lower <= value < upper
	idx := sort.Search(len(h.Bins), func(i int) bool {
		return value < h.Bins[i].UpperBound
	})
	if idx < len(h.Bins) {
		h.Bins[idx].Count++
		h.Total++
	}
}

// Counts returns the count for each bin.
func (h *Histogram) Counts() []int {
	counts := make([]int, len(h.Bins))
	for i, b := range h.Bins {
		counts[i] = b.Count
	}
	return counts
}

// Percentages returns the percentage of total observations that fall into each bin.
// If there are no observations, all percentages are 0.
func (h *Histogram) Percentages() []float64 {
	pcts := make([]float64, len(h.Bins))
	if h.Total == 0 {
		return pcts
	}
	for i, b := range h.Bins {
		pcts[i] = float64(b.Count) / float64(h.Total) * 100.0
	}
	return pcts
}

// String returns a simple textual representation of the histogram.
// It is intended for debugging; the format is not stable.
func (h *Histogram) String() string {
	if len(h.Bins) == 0 {
		return "Histogram{empty}"
	}
	var sb []byte
	sb = append(sb, fmt.Sprintf("Histogram{Total: %d, Bins: [", h.Total)...)
	for i, b := range h.Bins {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = append(sb, fmt.Sprintf("[%g, %g): %d", b.LowerBound, b.UpperBound, b.Count)...)
	}
	sb = append(sb, "]}"...)
	return string(sb)
}

// Equal compares two histograms for approximate equality.
// It is useful in tests where floating-point bounds may differ slightly.
func (h *Histogram) Equal(other *Histogram) bool {
	if h.Total != other.Total {
		return false
	}
	if len(h.Bins) != len(other.Bins) {
		return false
	}
	const eps = 1e-9
	for i := range h.Bins {
		if h.Bins[i].Count != other.Bins[i].Count {
			return false
		}
		if math.Abs(h.Bins[i].LowerBound-other.Bins[i].LowerBound) > eps {
			return false
		}
		if math.Abs(h.Bins[i].UpperBound-other.Bins[i].UpperBound) > eps {
			return false
		}
	}
	return true
}
