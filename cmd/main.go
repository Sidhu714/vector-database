package main

import (
	"fmt"
	"vector/internal/quantize"
)

func main() {
	vector := []float32{
		-1.0,
		-0.75,
		-0.5,
		-0.25,
		0.0,
		0.25,
		0.5,
		0.75,
		1.0,
	}

	quantized :=  quantize.Quantize(vector)

	fmt.Println("Original:")
	fmt.Println(vector)

	fmt.Println("Quantized:")
	fmt.Println(quantized.Values)

	dequantized := quantize.Dequantize(quantized)

	fmt.Println("DeQuantized:")
	fmt.Println(dequantized)
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
