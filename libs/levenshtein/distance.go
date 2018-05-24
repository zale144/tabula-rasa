package levenshtein

import (
	"math"
	"strings"
)

func Distance(s string, t string) int {
	// degenerate cases
	s = strings.ToLower(s)
	t = strings.ToLower(t)
	if s == t {
		return 0
	}
	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}

	// create two work vectors of integer distances
	v0 := make([]int, len(t)+1)
	v1 := make([]int, len(t)+1)

	// initialize v0 (the previous row of distances)
	// this row is A[0][i]: edit distance for an empty s
	// the distance is just the number of characters to delete from t
	for i := 0; i < len(v0); i++ {
		v0[i] = i
	}

	for i := 0; i < len(s); i++ {
		// calculate v1 (current row distances) from the previous row v0

		// first element of v1 is A[i+1][0]
		//   edit distance is delete (i+1) chars from s to match empty t
		v1[0] = i + 1

		// use formula to fill in the rest of the row
		for j := 0; j < len(t); j++ {
			var cost int
			if s[i] == t[j] {
				cost = 0
			} else {
				cost = 1
			}
			v1[j+1] = int(math.Min(float64(v1[j]+1), math.Min(float64(v0[j+1]+1), float64(v0[j]+cost))))
		}

		// copy v1 (current row) to v0 (previous row) for next iteration
		for j := 0; j < len(v0); j++ {
			v0[j] = v1[j]
		}
	}

	return v1[len(t)]
}
