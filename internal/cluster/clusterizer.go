package cluster

import (
	"fmt"
	"math"

	"github.com/boson-research/patterns/internal/cluster/kmeans"
	"github.com/boson-research/patterns/internal/context"
)

type Clusterizer struct {
	clusterer        Clusterer
	qualityEstimator qualityEstimator
}

func New(qualityEstimator QualityEstimationMethod) *Clusterizer {
	return &Clusterizer{
		// change for cluster builder with explicit parameter/s for optimization
		clusterer:        kmeans.New(1, 100, kmeans.RandomCentroidsIniter),
		qualityEstimator: getQualityEstimator(qualityEstimator),
	}
}

func (c *Clusterizer) Clusterize(ctx context.Context, data []float64) ([]float64, []int) {
	return c.optimize(ctx, data)
}

func (c *Clusterizer) optimize(ctx context.Context, data []float64) ([]float64, []int) {
	if len(data) == 1 {
		ctx.Logger().Debug("skipping optimization for number of clusters")

		return c.clusterer.Cluster(ctx, data)
	}

	var bestCentroids []float64
	var bestLabels []int
	var bestClustersNum int
	bestSilhouetteScore := math.Inf(-1)
	maxClusters := len(data) / 2
	for clustersNum := 1; clustersNum <= maxClusters; clustersNum++ {
		clusterer := kmeans.New(clustersNum, 100, kmeans.RandomCentroidsIniter)

		centroids, labels := clusterer.Cluster(ctx, data)
		silhouetteScore := c.qualityEstimator(data, labels)

		fmt.Println("silhouette score for", clustersNum, "clusters:", silhouetteScore)

		if silhouetteScore > bestSilhouetteScore {
			bestSilhouetteScore = silhouetteScore
			bestCentroids = centroids
			bestLabels = labels
			bestClustersNum = clustersNum
		}

		ctx.Logger().Tracef("silhouette score for %d clusters: %f", clustersNum, silhouetteScore)
	}

	ctx.Logger().Debugf("found optimal number of clusters: %d", bestClustersNum)

	return bestCentroids, bestLabels
}
