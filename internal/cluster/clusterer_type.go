package cluster

import "github.com/boson-research/patterns/internal/cluster/kmeans"

type ClustererType int

const (
	KMeans ClustererType = iota
)

func (t ClustererType) String() string {
	switch t {
	case KMeans:
		return "kmeans"
	default:
		return "unknown"
	}
}

func getClusterer(t ClustererType) Clusterer {
	switch t {
	case KMeans:
		return new(kmeans.KMeans)
	}
	return nil
}
