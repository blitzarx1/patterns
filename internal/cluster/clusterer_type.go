package cluster

type ClustererType int

const (
	KMeans ClustererType = iota
	KMeansPlusPLus
)

func (t ClustererType) String() string {
	switch t {
	case KMeans:
		return "kmeans"
	case KMeansPlusPLus:
		return "kmeans++"
	}
	return "unknown"
}
