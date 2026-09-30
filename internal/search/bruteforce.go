package search


import (
	"slices"
	"errors"
	"cmp"
	"vector/internal/vector"
)



// BruteForceSearch compares the query against every vector in the dataset.
//
// Input:
//   - query   → vector we want to search for
//   - dataset → vectors to search through
//   - k       → number of nearest vectors to return
//
// Output:
//   - []ScoredVector → k closest vectors with their distances
//   - error          → returned if k is invalid or distance calculation fails
//
// Process:
//   1. Calculate distance between query and every vector.
//   2. Sort vectors by distance (smallest first).
//   3. Return the top k closest vectors.

func BruteForceSearch(query []float32, dataset []vector.Vector, k int) ([]vector.ScoredVector, error) {

	distanceVector := []vector.ScoredVector{}

	if k > len(dataset) {
		return []vector.ScoredVector{}, errors.New("not enough vectors in the selected clusters to return k results")
	}

	for i := 0; i < len(dataset); i++ {

		distance, err := vector.EuclideanDistance(query, dataset[i].Values)

		if err != nil {
			return []vector.ScoredVector{}, errors.New("There was a error while calcualting the distance")
		}

		distanceVector = append(distanceVector, vector.ScoredVector{
			Vec:      dataset[i],
			Distance: distance,
		})
	}

	slices.SortFunc(distanceVector, func(a, b vector.ScoredVector) int {
		return cmp.Compare(a.Distance, b.Distance)
	})

	topKDistance := distanceVector[:k]

	return topKDistance, nil

}
