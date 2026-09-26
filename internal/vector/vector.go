package main

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

type Vector struct {
	Id     int
	Values []float32
}

type ScoredVector struct {
	Vec      Vector
	Distance float32
}

type Cluster struct {
	Centroid []float32
	Members  []Vector
}

// var idCounter uint64

// func nextID() uint64 {
// 	return atomic.AddUint64(&idCounter, 1)
// }

func EuclideanDistance(v1, v2 []float32) (float32, error) {

	if len(v1) != len(v2) {

		return 0, errors.New("vectors must have the same dimensions")
	}

	preDistance := 0.0

	for i := 0; i < len(v1); i++ {
		diff := v1[i] - v2[i]

		preDistance += float64(diff) * float64(diff)
	}

	distance := math.Sqrt(preDistance)

	return float32(distance), nil
}

// TODO: handle zero-vector case
func CosineSimilarity(v1, v2 []float32) (float32, error) {

	if len(v1) != len(v2) {
		return 0, errors.New("vectors must have the same dimensions")
	}

	dotProduct := 0.0
	var magV1 float64
	var magV2 float64

	for i := 0; i < len(v1); i++ {

		mul := v1[i] * v2[i]

		dotProduct += float64(mul)

		magV1 += float64((v1[i] * v1[i]))
		magV2 += float64((v2[i] * v2[i]))

	}

	cs := dotProduct / (math.Sqrt(magV1) * math.Sqrt(magV2))

	return float32(cs), nil
}

func GenerateRandomVectors(r *rand.Rand, count int, dimensions int) []Vector {

	vectors := []Vector{}
	idCounter := 0

	for i := 0; i < count; i++ {
		subarray := []float32{}
		for j := 0; j < dimensions; j++ {

			randomValue := (r.Float32() * 2) - 1

			subarray = append(subarray, randomValue)

		}
		idCounter++
		vectors = append(vectors, Vector{
			Id:     idCounter,
			Values: subarray,
		})

	}

	return vectors
}

func BruteForceSearch(query []float32, dataset []Vector, k int) ([]ScoredVector, error) {

	distances := []ScoredVector{}

	if k > len(dataset) {
		return []ScoredVector{}, errors.New("K shouldn't be greater than the dataset")
	}

	for i := 0; i < len(dataset); i++ {

		distance, err := EuclideanDistance(query, dataset[i].Values)

		if err != nil {
			return []ScoredVector{}, errors.New("There was a error while calcualting the distance")
		}

		distances = append(distances, ScoredVector{
			Vec:      dataset[i],
			Distance: distance,
		})
	}

	slices.SortFunc(distances, func(a, b ScoredVector) int {
		return cmp.Compare(a.Distance, b.Distance)
	})

	topKDistance := distances[:k]

	return topKDistance, nil

}

func KMeans(r *rand.Rand, dataset []Vector, k int, maxIteration int) []Cluster {

	cluster := []Cluster{}

	for i := 0; i < k; i++ {
		randIndex := r.IntN(len(dataset))
		cluster = append(cluster, Cluster{
			Centroid: dataset[randIndex].Values,
		})
	}

	for iteration := 0; iteration < maxIteration; iteration++ {

		for i := range cluster {
			cluster[i].Members = nil
		}

		for i := 0; i < len(dataset); i++ {
			smallDistance := float32(math.MaxFloat32)
			closestClusterIndex := 0
			for j := 0; j < len(cluster); j++ {

				distance, err := EuclideanDistance(dataset[i].Values, cluster[j].Centroid)

				if err != nil {
					fmt.Println("Error calculating distance")
					return []Cluster{}
				}

				if distance < smallDistance {
					smallDistance = distance
					closestClusterIndex = j
				}

			}
			cluster[closestClusterIndex].Members = append(cluster[closestClusterIndex].Members, dataset[i])

		}

		for i := range cluster {

			if len(cluster[i].Members) == 0 {
				continue // keep old centroid, skip recompute
			}
			
			newCentroids := RecomputeCentroid(cluster[i].Members)

			cluster[i].Centroid = newCentroids
		}

	}

	return cluster

}

func RecomputeCentroid(members []Vector) []float32 {
	dimensions := len(members[0].Values)
	newCentroid := make([]float32, dimensions)

	for i := 0; i < dimensions; i++ {

		var sum float32

		for _, vector := range members {

			sum += vector.Values[i]
		}

		newCentroid[i] = sum / float32(len(members))
	}

	return newCentroid
}

func main() {

	source := rand.NewPCG(42, 999)
	r := rand.New(source)

	datasets := GenerateRandomVectors(r, 10, 2)
	// query := GenerateRandomVectors(r, 1, 128)

	// start := time.Now()
	// topK, err := BruteForceSearch(query[0].Values, datasets, 100)
	// elapsed := time.Since(start)

	// fmt.Println("N =", 100000, "took", elapsed)

	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }

	// for i := 0; i < len(topK); i++ {
	// 	fmt.Printf("The id is %d and the distance is %f\n", topK[i].Vec.Id, topK[i].Distance)
	// }

	cluster := KMeans(r, datasets, 8, 2)

	for i := 0; i < len(cluster); i++ {
		fmt.Printf("The cluster %d: %+v\n", i, cluster[i])
	}

	// vectors := []Vector{
	// 	{
	// 		Id:     8,
	// 		Values: []float32{0.08737445, -0.29944146},
	// 	},
	// }

	// reCompute := RecomputeCentroid(vectors)

	// fmt.Println(reCompute)

}
