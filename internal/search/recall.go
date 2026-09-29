package search


func RecallAtK(groundTruth []ScoredVector, approx []ScoredVector) float32 {

	idMap := make(map[int]bool)
	count := 0

	if len(groundTruth) == 0 {
		return 0.0
	}

	for i := 0; i < len(groundTruth); i++ {

		id := groundTruth[i].Vec.Id

		idMap[id] = true

	}

	for i := 0; i < len(approx); i++ {

		id := approx[i].Vec.Id

		if idMap[id] {
			count++
		}
	}

	return float32(count) / float32(len(groundTruth))

}