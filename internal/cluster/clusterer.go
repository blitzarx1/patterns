package cluster

import "github.com/boson-research/patterns/internal/context"

type Clusterer interface {
	Init(ctx context.Context, data []float64)
	GetOptimizationParams(ctx context.Context) []int
	SetOptimizationParams(ctx context.Context, params []int) error
	Cluster(ctx context.Context) (clusters []float64, labels []int)
	ValidateOptimizationParams(params []int) error
}
