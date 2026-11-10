# histogram

A small Go library for bucketing numeric values into configurable bins and reporting counts and percentages per bin.

## Usage

```go
package main

import (
	"fmt"
	"histogram"
)

func main() {
	h := histogram.New(4, 0, 100)
	h.Add(10)
	h.Add(25)
	h.Add(25)
	h.Add(80)
	h.Add(100)

	fmt.Println(h.Counts())      // [1 2 0 2]
	fmt.Println(h.Percentages()) // [20 40 0 40]
	fmt.Println(h)
}
```

## Why this exists

When analysing data, you often need a quick way to see how values are distributed. This library provides a minimal, dependency-free histogram that works with a fixed number of equal-width bins. The trade-off is simplicity: it does not support adaptive binning, weights, or non-uniform bins. If you need those, you should look elsewhere.

## Edge cases

- Values outside the range `[min, max]` are ignored.
- The upper bound of the last bin is inclusive; all other bins are half-open `[lower, upper)`.
- If `min == max`, a single bin covering that value is created.
- `New` panics if `bins < 1` or `min > max`.
