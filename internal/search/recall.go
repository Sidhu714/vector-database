package search

import "vector/internal/vector"

func RecallAtK(
	groundTruth []vector.ScoredVector,
	approx []vector.ScoredVector,
	k int,
) float32 {

	if k <= 0 {
		return 0
	}

	if len(groundTruth) < k {
		k = len(groundTruth)
	}

	if len(approx) < k {
		k = len(approx)
	}

	idMap := make(map[int]bool)

	for i := 0; i < k; i++ {
		idMap[groundTruth[i].Vec.Id] = true
	}

	count := 0

	for i := 0; i < k; i++ {
		if idMap[approx[i].Vec.Id] {
			count++
		}
	}

	return float32(count) / float32(k)
}
