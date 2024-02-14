package kmeans

import (
	"math"

	"github.com/boson-research/patterns/internal/context"
)

func New(clustersNum int, maxIterations int, centroidsIniter CentroidsIniterType) *KMeans {
	return &KMeans{
		clustersNum:     clustersNum,
		maxIterations:   maxIterations,
		centroidsIniter: getCentroidsIniter(centroidsIniter),
	}
}

type KMeans struct {
	clustersNum     int
	maxIterations   int
	centroidsIniter centroidsIniter
}

func (k *KMeans) Cluster(ctx context.Context, data []float64) ([]float64, []int) {
	ctx, span := ctx.StartSpan("KMeans")
	defer span.End()

	ctx.Logger().Tracef("clustering %d points into %d clusters", len(data), k.clustersNum)

	centroids := k.centroidsIniter(data, k.clustersNum)
	for i := 0; i < k.maxIterations; i++ {
		labels := assignPointsToCentroids(data, centroids)
		newCentroids := updateCentroids(data, labels, k.clustersNum)
		if checkConvergence(centroids, newCentroids, 1e-5) {
			ctx.Logger().Tracef("converged after %d iterations", i+1)

			return newCentroids, labels
		}

		centroids = newCentroids
	}

	return centroids, assignPointsToCentroids(data, centroids)
}

// assignPointsToCentroids assigns each data point to the nearest centroid and returns the labels.
func assignPointsToCentroids(data []float64, centroids []float64) []int {
	labels := make([]int, len(data))
	for i, point := range data {
		minDist := math.MaxFloat64
		for j, centroid := range centroids {
			dist := math.Abs(point - centroid)
			if dist < minDist {
				minDist = dist
				labels[i] = j
			}
		}
	}
	return labels
}

// updateCentroids recalculates the centroids based on the assigned points.
func updateCentroids(data []float64, labels []int, k int) []float64 {
	sums := make([]float64, k)
	counts := make([]int, k)
	for i, label := range labels {
		sums[label] += data[i]
		counts[label]++
	}
	for i := range sums {
		if counts[i] > 0 { // avoid division by zero.
			sums[i] /= float64(counts[i])
		}
	}
	return sums
}

// checkConvergence tests if the centroids have changed significantly.
func checkConvergence(oldCentroids, newCentroids []float64, threshold float64) bool {
	for i := range oldCentroids {
		if math.Abs(oldCentroids[i]-newCentroids[i]) > threshold {
			return false
		}
	}
	return true
}

