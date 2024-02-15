package kmeans

import (
	"context"
	"fmt"
	"math"

	"github.com/boson-research/patterns/internal/telemetry/logger"
	"go.opentelemetry.io/otel"
)

type KMeans struct {
	clustersNum         int
	maxIterations       int
	clustersNumBounds   [2]int
	centroidsIniterType CentroidsIniterType
	data                []float64
}

func (k *KMeans) Init(ctx context.Context, data []float64) {
	k.data = data
	k.centroidsIniterType = RandomCentroidsIniter
	k.clustersNumBounds = [2]int{1, len(data) / 2}
	k.maxIterations = 100
}

func (k *KMeans) GetOptimizationParams(ctx context.Context) []int {
	return []int{k.clustersNum, int(k.centroidsIniterType)}
}

func (k *KMeans) SetOptimizationParams(ctx context.Context, params []int) error {
	if err := k.ValidateOptimizationParams(params); err != nil {
		return err
	}

	k.clustersNum = params[0]
	k.maxIterations = params[1]

	return nil
}

func (k *KMeans) ValidateOptimizationParams(params []int) error {
	if params[0] < k.clustersNumBounds[0] || params[0] > k.clustersNumBounds[1] {
		return fmt.Errorf("number of clusters must be between %d and %d", k.clustersNumBounds[0], k.clustersNumBounds[1])
	}

	if CentroidsIniterType(params[1]).String() == "unknown" {
		return fmt.Errorf("unknown centroids initer type")
	}

	return nil
}

func (k *KMeans) Cluster(ctx context.Context) ([]float64, []int) {
	ctx, span := otel.Tracer("").Start(ctx, "KMeans")
	defer span.End()

	logger.MustFromContext(ctx).Tracef("clustering %d points into %d clusters", len(k.data), k.clustersNum)

	centroids := getCentroidsIniter(k.centroidsIniterType)(k.data, k.clustersNum)
	for i := 0; i < k.maxIterations; i++ {
		labels := assignPointsToCentroids(k.data, centroids)
		newCentroids := updateCentroids(k.data, labels, k.clustersNum)
		if checkConvergence(centroids, newCentroids, 1e-5) {
			logger.MustFromContext(ctx).Tracef("converged after %d iterations", i+1)

			return newCentroids, labels
		}

		centroids = newCentroids
	}

	return centroids, assignPointsToCentroids(k.data, centroids)
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
