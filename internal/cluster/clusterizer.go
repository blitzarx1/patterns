package cluster

import (
	"context"
	"fmt"
	"math"

	"github.com/boson-research/patterns/internal/telemetry/logger"
)

type Clusterizer struct {
	clusterer        Clusterer
	qualityEstimator qualityEstimator
}

func New(clusterer ClustererType, qualityEstimator QualityEstimationMethod) *Clusterizer {
	return &Clusterizer{
		clusterer:        getClusterer(clusterer),
		qualityEstimator: getQualityEstimator(qualityEstimator),
	}
}

func (c *Clusterizer) Clusterize(ctx context.Context, data []float64) ([]float64, []int) {
	c.clusterer.Init(ctx, data)
	return c.optimize(ctx, data)
}

func (c *Clusterizer) optimize(ctx context.Context, data []float64) ([]float64, []int) {
	if len(data) == 1 {
		logger.MustFromContext(ctx).Debug("skipping optimization for number of clusters")

		return c.clusterer.Cluster(ctx)
	}

	optimizationParams := c.clusterer.GetOptimizationParams(ctx)

	var bestCentroids []float64
	var bestLabels []int
	var bestParams []int
	bestScore := math.Inf(-1)

	// TODO: parallelize
	for _, params := range generateOptimizationParamsVariations(optimizationParams, c.clusterer.ValidateOptimizationParams) {
		if err := c.clusterer.SetOptimizationParams(ctx, params); err != nil {
			logger.MustFromContext(ctx).Errorf("failed to set optimization params: %v", err)
			continue
		}

		centroids, labels := c.clusterer.Cluster(ctx)
		score := c.qualityEstimator(data, labels)

		fmt.Printf("quality score for %v params: %.2f\n", params, score)

		if score > bestScore {
			bestScore = score
			bestCentroids = centroids
			bestLabels = labels
			bestParams = params
		}

		logger.MustFromContext(ctx).Tracef("quality score for %v params: %.2f", params, score)
	}

	logger.MustFromContext(ctx).Debugf("found optimal score for %v params: %.2f", bestParams, bestScore)

	fmt.Printf("found optimal score for %v params: %.2f\n", bestParams, bestScore)

	return bestCentroids, bestLabels
}

func generateOptimizationParamsVariations(startingParams []int, validator func(params []int) error) [][]int {
	var variations [][]int

	// create a copy of startingParams to avoid modifying the original slice
	params := make([]int, len(startingParams))
	copy(params, startingParams)

	var generate func(int)
	generate = func(idx int) {
		if idx == len(params) {
			if err := validator(params); err == nil {
				// make a copy of params to avoid modifying the slice later
				validParams := make([]int, len(params))
				copy(validParams, params)
				variations = append(variations, validParams)
			}
			return
		}

		for {
			generate(idx + 1)
			params[idx]++
			if err := validator(params); err != nil {
				params[idx] = startingParams[idx]
				break
			}
		}
	}

	generate(0)
	return variations
}
