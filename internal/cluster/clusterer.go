package cluster

import "github.com/boson-research/patterns/internal/context"

type Clusterer interface {
	Cluster(ctx context.Context, data []float64) (clusters []float64, labels []int)
}

