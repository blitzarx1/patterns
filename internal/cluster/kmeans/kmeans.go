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

// calculateAverageDistance calculates the average distance of a given point to other points in the cluster.
func calculateAverageDistance(point float64, cluster []float64) float64 {
	if len(cluster) <= 1 {
		// If the cluster contains only the point itself or is empty, return 0 because there are no other points to compare.
		return 0.0
	}

	sumDistance := 0.0
	for _, otherPoint := range cluster {
		sumDistance += math.Abs(point - otherPoint)
	}

	// To avoid biasing the average distance by including the point itself in its own cluster,
	// subtract the distance of the point to itself (which is 0) and decrement the count by 1.
	averageDistance := sumDistance / float64(len(cluster)-1)

	return averageDistance
}

// calculateSilhouetteScore calculates the silhouette score for each point and returns the average score.
func calculateSilhouetteScore(data []float64, labels []int) float64 {
	// Create clusters from labels
	clusters := make(map[int][]float64)
	for i, label := range labels {
		clusters[label] = append(clusters[label], data[i])
	}

	totalScore := 0.0
	for i, point := range data {
		// Calculate a(i)
		a := calculateAverageDistance(point, clusters[labels[i]])

		// Calculate b(i)
		b := math.Inf(1)
		for label, cluster := range clusters {
			if label == labels[i] {
				continue // Skip own cluster
			}
			dist := calculateAverageDistance(point, cluster)
			if dist < b {
				b = dist
			}
		}

		// Calculate silhouette score for point i
		si := (b - a) / math.Max(a, b)
		totalScore += si
	}

	// Return average silhouette score
	return totalScore / float64(len(data))
}

func kmeans(ctx context.Context, data []float64, k int, maxIterations int) ([]float64, []int) {
	ctx, span := otel.Tracer("kmeans").Start(ctx, "kmeans")
	defer span.End()
	l := logger.Logger(ctx)

	l.Tracef("clustering %d points into %d clusters", len(data), k)

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

// KMeans performs 1D k-means clustering with refactored helper functions.
func KMeans(ctx context.Context, data []float64) ([]float64, []int) {
	ctx, span := otel.Tracer("kmeans").Start(ctx, "KMeans")
	defer span.End()
	l := logger.Logger(ctx)

	l.Debug("clustering")

	var bestCentroids []float64
	var bestLabels []int
	var bestK int
	bestSilhouetteScore := math.Inf(-1)
	maxClusters := len(data) - 1

	if maxClusters == 0 {
		l.Debug("skipping silhouete number of clusters optimization")

		return kmeans(ctx, data, 1, 100)
	}

	for k := 1; k <= maxClusters; k++ {
		centroids, labels := kmeans(ctx, data, k, 100)
		silhouetteScore := calculateSilhouetteScore(data, labels)

		if silhouetteScore > bestSilhouetteScore {
			bestSilhouetteScore = silhouetteScore
			bestCentroids = centroids
			bestLabels = labels
			bestK = k
		}
	}

	l.Debugf("found optimal number of clusters: %d", bestK)

	return bestCentroids, bestLabels
}
