package cluster

import "math"

type QualityEstimationMethod int

const (
	Silhouette QualityEstimationMethod = iota
	Elbow
)

func (t QualityEstimationMethod) String() string {
	switch t {
	case Silhouette:
		return "silhouette"
	case Elbow:
		return "elbow"
	}
	return "unknown"
}

type qualityEstimator func(data []float64, labels []int) float64

func getQualityEstimator(t QualityEstimationMethod) qualityEstimator {
	switch t {
	case Silhouette:
		return calculateSilhouette
	case Elbow:
		panic("elbow method not implemented yet")
	}
	return nil
}

func calcAvgDistOwn(p float64, cluster []float64) float64 {
	if len(cluster) == 1 {
		return 0.0
	}

	sumDistance := 0.0
	for _, otherPoint := range cluster {
		sumDistance += math.Abs(p - otherPoint)
	}

	return sumDistance / float64(len(cluster)-1)
}

func calcAvgDistOther(p float64, cluster []float64) float64 {
	sumDistance := 0.0
	for _, otherPoint := range cluster {
		sumDistance += math.Abs(p - otherPoint)
	}

	return sumDistance / float64(len(cluster))
}

// calculateSilhouetteScore calculates the silhouette score for each point and returns the average score.
func calculateSilhouette(data []float64, labels []int) float64 {
	// create clusters from labels
	clusters := make(map[int][]float64)
	for i, label := range labels {
		clusters[label] = append(clusters[label], data[i])
	}

	totalScore := 0.0
	for i, point := range data {
		// calculate a(i)
		a := calcAvgDistOwn(point, clusters[labels[i]])

		// calculate b(i)
		b := math.Inf(1)
		for label, cluster := range clusters {
			if label == labels[i] {
				continue // skip own cluster
			}

			dist := calcAvgDistOther(point, cluster)
			if dist < b {
				b = dist
			}
		}

		// calculate silhouette score for point i
		si := (b - a) / math.Max(a, b)
		totalScore += si
	}

	// return average silhouette score
	return totalScore / float64(len(data))
}
