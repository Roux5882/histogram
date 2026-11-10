package histogram

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestNewSingleBin(t *testing.T) {
	h := New(1, 0, 10)
	if len(h.Bins) != 1 {
		t.Fatalf("expected 1 bin, got %d", len(h.Bins))
	}
	if !almostEqual(h.Bins[0].LowerBound, 0) || !almostEqual(h.Bins[0].UpperBound, 10) {
		t.Errorf("bin bounds = [%v, %v), want [0, 10)", h.Bins[0].LowerBound, h.Bins[0].UpperBound)
	}
}

func TestNewMultipleBins(t *testing.T) {
	h := New(5, 0, 10)
	if len(h.Bins) != 5 {
		t.Fatalf("expected 5 bins, got %d", len(h.Bins))
	}
	for i, b := range h.Bins {
		wantLower := float64(i) * 2
		wantUpper := float64(i+1) * 2
		if !almostEqual(b.LowerBound, wantLower) {
			t.Errorf("bin %d lower = %v, want %v", i, b.LowerBound, wantLower)
		}
		if !almostEqual(b.UpperBound, wantUpper) {
			t.Errorf("bin %d upper = %v, want %v", i, b.UpperBound, wantUpper)
		}
	}
}

func TestNewDegenerateRange(t *testing.T) {
	h := New(3, 5, 5)
	if len(h.Bins) != 1 {
		t.Fatalf("expected 1 bin, got %d", len(h.Bins))
	}
	if !almostEqual(h.Bins[0].LowerBound, 5) || !almostEqual(h.Bins[0].UpperBound, 5) {
		t.Errorf("bin bounds = [%v, %v), want [5, 5)", h.Bins[0].LowerBound, h.Bins[0].UpperBound)
	}
}

func TestNewPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for bins < 1")
		}
	}()
	New(0, 0, 1)
}

func TestAdd(t *testing.T) {
	h := New(2, 0, 10)
	h.Add(1)
	h.Add(5)
	h.Add(9)
	h.Add(10)
	counts := h.Counts()
	if counts[0] != 1 {
		t.Errorf("bin 0 count = %d, want 1", counts[0])
	}
	if counts[1] != 3 {
		t.Errorf("bin 1 count = %d, want 3", counts[1])
	}
	if h.Total != 4 {
		t.Errorf("total = %d, want 4", h.Total)
	}
}

func TestAddOutOfRange(t *testing.T) {
	h := New(2, 0, 10)
	h.Add(-1)
	h.Add(11)
	if h.Total != 0 {
		t.Errorf("total = %d, want 0", h.Total)
	}
}

func TestAddBoundary(t *testing.T) {
	h := New(2, 0, 10)
	h.Add(0)
	h.Add(5)
	h.Add(10)
	counts := h.Counts()
	if counts[0] != 1 {
		t.Errorf("bin 0 count = %d, want 1", counts[0])
	}
	if counts[1] != 2 {
		t.Errorf("bin 1 count = %d, want 2", counts[1])
	}
}

func TestPercentages(t *testing.T) {
	h := New(2, 0, 10)
	h.Add(1)
	h.Add(2)
	h.Add(7)
	pcts := h.Percentages()
	want0 := 2.0 / 3.0 * 100
	want1 := 1.0 / 3.0 * 100
	if !almostEqual(pcts[0], want0) {
		t.Errorf("pct[0] = %v, want %v", pcts[0], want0)
	}
	if !almostEqual(pcts[1], want1) {
		t.Errorf("pct[1] = %v, want %v", pcts[1], want1)
	}
}

func TestPercentagesEmpty(t *testing.T) {
	h := New(3, 0, 3)
	pcts := h.Percentages()
	for i, p := range pcts {
		if p != 0 {
			t.Errorf("pct[%d] = %v, want 0", i, p)
		}
	}
}

func TestString(t *testing.T) {
	h := New(2, 0, 10)
	h.Add(1)
	s := h.String()
	if s == "" {
		t.Errorf("String() returned empty string")
	}
	// Basic check that it contains expected substrings.
	if !contains(s, "Total: 1") {
		t.Errorf("String() = %q, missing Total: 1", s)
	}
	if !contains(s, "[0, 5): 1") {
		t.Errorf("String() = %q, missing [0, 5): 1", s)
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestEqual(t *testing.T) {
	h1 := New(2, 0, 10)
	h1.Add(1)
	h1.Add(6)
	h2 := New(2, 0, 10)
	h2.Add(1)
	h2.Add(6)
	if !h1.Equal(h2) {
		t.Errorf("histograms should be equal")
	}
	h3 := New(2, 0, 10)
	h3.Add(2)
	if h1.Equal(h3) {
		t.Errorf("histograms should not be equal")
	}
}
