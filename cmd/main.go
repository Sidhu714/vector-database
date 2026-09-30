package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"sort"
	"vector/internal/quantize"
	"vector/internal/search"
	"vector/internal/vector"
)

type QuantizedVector struct {
	Id        int
	Quantized quantize.QuantizedVector
}

func main() {
	source := rand.NewPCG(42, 999)
	r := rand.New(source)

	datasetSize := 10000
	queryCount := 100
	dimensions := 128
	k := 10

	// --------------------------------
	// Generate original dataset
	// --------------------------------

	nonQuantizedDataset := vector.GenerateRandomVectors(
		r,
		datasetSize,
		dimensions,
	)

	// --------------------------------
	// Calculate Float32 dataset size
	// --------------------------------

	var float32Bytes int

	for _, v := range nonQuantizedDataset {
		float32Bytes += len(v.Values) * 4
	}

	fmt.Printf(
		"Float32 dataset memory: %.2f KB\n",
		float64(float32Bytes)/1024,
	)

	// --------------------------------
	// Create quantized dataset
	// --------------------------------

	quantizedVector := make(
		[]QuantizedVector,
		0,
		len(nonQuantizedDataset),
	)

	for _, v := range nonQuantizedDataset {

		quantizedValue := quantize.Quantize(
			v.Values,
		)

		quantizedVector = append(
			quantizedVector,
			QuantizedVector{
				Id:        v.Id,
				Quantized: quantizedValue,
			},
		)
	}

	// --------------------------------
	// Calculate Quantized dataset size
	// --------------------------------

	var quantizedBytes int

	for _, v := range quantizedVector {

		// uint8 values
		quantizedBytes += len(v.Quantized.Values)

		// MinValue + Scale
		quantizedBytes += 4
		quantizedBytes += 4
	}

	fmt.Printf(
		"Quantized dataset memory: %.2f KB\n",
		float64(quantizedBytes)/1024,
	)

	// --------------------------------
	// Memory reduction
	// --------------------------------

	reduction := (1 -
		float64(quantizedBytes)/float64(float32Bytes)) * 100

	fmt.Printf(
		"Memory reduction: %.2f%%\n",
		reduction,
	)

	// --------------------------------
	// Run queries
	// --------------------------------

	totalRecall := float32(0)

	for q := 0; q < queryCount; q++ {

		query := vector.GenerateRandomVectors(
			r,
			1,
			dimensions,
		)

		// --------------------------------
		// Ground truth search
		// --------------------------------

		nonQuantizedSearch, err :=
			search.BruteForceSearch(
				query[0].Values,
				nonQuantizedDataset,
				k,
			)

		if err != nil {
			log.Fatalf(
				"Brute force search failed for query %d: %v",
				q+1,
				err,
			)
		}

		// --------------------------------
		// Quantized search
		// --------------------------------

		quantizedSearch :=
			make([]vector.ScoredVector, 0, len(quantizedVector))

		for _, qv := range quantizedVector {

			dequantized :=
				quantize.Dequantize(
					qv.Quantized,
				)

			distance, err :=
				vector.EuclideanDistance(
					query[0].Values,
					dequantized,
				)

			if err != nil {
				log.Fatalf(
					"Quantized search failed for query %d: %v",
					q+1,
					err,
				)
			}

			quantizedSearch = append(
				quantizedSearch,
				vector.ScoredVector{
					Vec: vector.Vector{
						Id:     qv.Id,
						Values: dequantized,
					},
					Distance: distance,
				},
			)
		}

		sort.Slice(
			quantizedSearch,
			func(i, j int) bool {
				return quantizedSearch[i].Distance <
					quantizedSearch[j].Distance
			},
		)

		quantizedSearch =
			quantizedSearch[:k]

		// --------------------------------
		// Recall@K
		// --------------------------------

		recall := search.RecallAtK(
			nonQuantizedSearch,
			quantizedSearch,
			k,
		)

		totalRecall += recall
	}

	// --------------------------------
	// Average Recall
	// --------------------------------

	averageRecall :=
		totalRecall / float32(queryCount)

	fmt.Println()
	fmt.Printf(
		"Dataset: %d vectors\n",
		datasetSize,
	)

	fmt.Printf(
		"Queries: %d\n",
		queryCount,
	)

	fmt.Printf(
		"Dimensions: %d\n",
		dimensions,
	)

	fmt.Printf(
		"K: %d\n",
		k,
	)

	fmt.Printf(
		"Average Recall@%d: %.2f%%\n",
		k,
		averageRecall*100,
	)
}

// func main() {

// 	source := rand.NewPCG(42, 999)
// 	r := rand.New(source)

// 	// Configuration
// 	datasetSize := 1000
// 	dimensions := 128
// 	numClusters := 20
// 	k := 10
// 	numQueries := 20

// 	// Generate dataset
// 	datasets := GenerateRandomVectors(
// 		r,
// 		datasetSize,
// 		dimensions,
// 	)

// 	// Generate queries
// 	queries := GenerateRandomVectors(
// 		r,
// 		numQueries,
// 		dimensions,
// 	)

// 	// Build IVF index
// 	fmt.Println("Building IVF index...")

// 	clusters := KMeans(
// 		r,
// 		datasets,
// 		numClusters,
// 		3,
// 	)

// 	fmt.Println("IVF index built.")

// 	// =========================================
// 	// BRUTE FORCE BENCHMARK
// 	// =========================================

// 	var bruteForceTotalTime time.Duration

// 	// Store ground truth for each query
// 	groundTruths := make([][]ScoredVector, len(queries))

// 	for i, query := range queries {

// 		start := time.Now()

// 		groundTruth, err := BruteForceSearch(
// 			query.Values,
// 			datasets,
// 			k,
// 		)

// 		elapsed := time.Since(start)

// 		if err != nil {
// 			fmt.Println("Error during brute force search:", err)
// 			return
// 		}

// 		bruteForceTotalTime += elapsed
// 		groundTruths[i] = groundTruth
// 	}

// 	bruteForceAverageTime :=
// 		bruteForceTotalTime / time.Duration(numQueries)

// 	// =========================================
// 	// PRINT RESULTS
// 	// =========================================

// 	fmt.Println()
// 	fmt.Println("==============================================")
// 	fmt.Println("              VECTOR SEARCH BENCHMARK")
// 	fmt.Println("==============================================")

// 	fmt.Printf(
// 		"Brute Force | Recall@%d = 1.00 (100%%) | Avg Time = %v\n",
// 		k,
// 		bruteForceAverageTime,
// 	)

// 	fmt.Println("----------------------------------------------")

// 	// =========================================
// 	// IVF BENCHMARK
// 	// =========================================

// 	nprobes := []int{1, 2, 5, 10, 20}

// 	for _, nprobe := range nprobes {

// 		var totalRecall float32
// 		var totalIVFTime time.Duration

// 		for i, query := range queries {

// 			start := time.Now()

// 			approx, err := IVFSearch(
// 				query.Values,
// 				clusters,
// 				nprobe,
// 				k,
// 			)

// 			elapsed := time.Since(start)

// 			if err != nil {
// 				fmt.Println("Error during IVF search:", err)
// 				return
// 			}

// 			totalIVFTime += elapsed

// 			// Compare IVF result with brute-force ground truth
// 			recall := RecallAtK(
// 				groundTruths[i],
// 				approx,
// 			)

// 			totalRecall += recall
// 		}

// 		averageRecall :=
// 			totalRecall / float32(numQueries)

// 		averageIVFTime :=
// 			totalIVFTime / time.Duration(numQueries)

// 		fmt.Printf(
// 			"IVF nprobe=%-2d | Recall@%d = %.2f (%3.0f%%) | Avg Time = %v\n",
// 			nprobe,
// 			k,
// 			averageRecall,
// 			averageRecall*100,
// 			averageIVFTime,
// 		)
// 	}

// 	fmt.Println("==============================================")
// }
