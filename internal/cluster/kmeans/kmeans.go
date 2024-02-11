package kmeans

import (
	"context"
	"math"
	"math/rand"

	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

// initializeCentroids selects k unique random points from the data as the initial centroids.
func initializeCentroids(data []float64, k int) []float64 {
	centroids := make([]float64, k)
	perm := rand.Perm(len(data))
	for i := 0; i < k; i++ {
		centroids[i] = data[perm[i]]
	}
	return centroids
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

// KMeans performs 1D k-means clustering with refactored helper functions.
func KMeans(ctx context.Context, data []float64, k int, maxIterations int) ([]float64, []int) {
	ctx, span := otel.Tracer("kmeans").Start(ctx, "KMeans")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debugf("clustering %d points into %d clusters", len(data), k)

	centroids := initializeCentroids(data, k)
	for i := 0; i < maxIterations; i++ {
		labels := assignPointsToCentroids(data, centroids)
		newCentroids := updateCentroids(data, labels, k)
		if checkConvergence(centroids, newCentroids, 1e-5) {
			return newCentroids, labels
		}
		centroids = newCentroids
	}
	return centroids, assignPointsToCentroids(data, centroids)
}
