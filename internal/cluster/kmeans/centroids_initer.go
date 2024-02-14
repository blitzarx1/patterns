package kmeans

import (
	"math"
	"math/rand"
)

type CentroidsIniterType int

const (
	RandomCentroidsIniter CentroidsIniterType = iota
	PlusPlusCentroidsIniter
)

func (t CentroidsIniterType) String() string {
	switch t {
	case RandomCentroidsIniter:
		return "random"
	case PlusPlusCentroidsIniter:
		return "plusplus"
	}
	return "unknown"
}

type centroidsIniter func(data []float64, k int) []float64

func getCentroidsIniter(t CentroidsIniterType) centroidsIniter {
	switch t {
	case RandomCentroidsIniter:
		return initializeCentroidsRandom
	case PlusPlusCentroidsIniter:
		return initializeCentroidsKMeansPlusPlus
	}
	return nil
}

// initializeCentroidsRandom selects k unique random points from the data as the initial centroids.
func initializeCentroidsRandom(data []float64, k int) []float64 {
	centroids := make([]float64, k)
	perm := rand.Perm(len(data))
	for i := 0; i < k; i++ {
		centroids[i] = data[perm[i]]
	}
	return centroids
}

// initializeCentroidsKMeansPlusPlus selects k unique centroids using the k-means++ algorithm.
func initializeCentroidsKMeansPlusPlus(data []float64, k int) []float64 {
	if len(data) == 0 || k <= 0 {
		return nil // handle edge cases
	}

	centroids := make([]float64, 0, k)
	// randomly select the first centroid from the data points.
	firstCentroidIndex := rand.Intn(len(data))
	centroids = append(centroids, data[firstCentroidIndex])

	// repeat until we have k centroids
	for len(centroids) < k {
		distances := make([]float64, len(data))
		totalDistance := 0.0

		// for each data point, compute the distance to the nearest centroid
		for i, point := range data {
			minDist := math.Inf(1)
			for _, centroid := range centroids {
				dist := math.Abs(point - centroid)
				if dist < minDist {
					minDist = dist
				}
			}
			distances[i] = minDist * minDist // square the distance to increase probability for farther points
			totalDistance += distances[i]
		}

		// select the next centroid
		randomPoint := rand.Float64() * totalDistance
		for i, d := range distances {
			randomPoint -= d
			if randomPoint <= 0 {
				centroids = append(centroids, data[i])
				break
			}
		}
	}

	return centroids
}
