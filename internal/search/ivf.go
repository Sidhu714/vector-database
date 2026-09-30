package search

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"vector/internal/vector"
)


type Cluster struct {
	Centroid []float32
	Members  []vector.Vector
}


type ClusterData struct {
	Index    int
	Distance float32
}


func KMeans(r *rand.Rand, dataset []vector.Vector, k int, maxIteration int) []Cluster {

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

				distance, err := vector.EuclideanDistance(dataset[i].Values, cluster[j].Centroid)

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



func IVFSearch(query []float32, clusters []Cluster, nprobe int, k int) ([]vector.ScoredVector, error) {

	// k used here is how many final results to return (what the user actually asked for) !!
	// nprobe controls how many clusters to look inside (an index-time/search-time speed knob)

	NCentroid := make([]ClusterData, 0, len(clusters))
	var memberCoimbed []vector.Vector

	if nprobe > len(clusters) {
		return []vector.ScoredVector{}, errors.New("nprobe cannot be greater than the number of clusters")
	}

	for i := 0; i < len(clusters); i++ {
		distance, err := vector.EuclideanDistance(query, clusters[i].Centroid)

		if err != nil {
			return []vector.ScoredVector{}, err
		}

		NCentroid = append(NCentroid, ClusterData{
			Index:    i,
			Distance: distance,
		})

	}

	sort.Slice(NCentroid, func(a, b int) bool {
		return NCentroid[a].Distance < NCentroid[b].Distance
	})

	NCentroid = NCentroid[:nprobe]

	for i := 0; i < len(NCentroid); i++ {

		clusterIndex := NCentroid[i].Index

		memberCoimbed = append(memberCoimbed, clusters[clusterIndex].Members...)
	}

	results, err := BruteForceSearch(query, memberCoimbed, k)

	if err != nil {
		return []vector.ScoredVector{}, err
	}

	return results, nil

}

func RecomputeCentroid(members []vector.Vector) []float32 {
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