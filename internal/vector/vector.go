package vector

import (
	"errors"
	"math"
	"math/rand/v2"
)

type Vector struct {
	Id     int
	Values []float32
}

type ScoredVector struct {
	Vec      Vector
	Distance float32
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

