package main

import (
	"errors"
	"fmt"
	"math"
)

type Vector struct {
	Id     int
	Values []float32
}

func EuclideanDistance(v1, v2 []float32) (float32,error)  {

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

func main() {
	v1 := []float32{7.0, 2.0, 3.0}
	v2 := []float32{6.0, 1.0}

	distance, err := EuclideanDistance(v1, v2)

	if err != nil {
		fmt.Println("Error while calculating")
		return
	}

	fmt.Printf("The Eucledian distance is %f", distance)
}
