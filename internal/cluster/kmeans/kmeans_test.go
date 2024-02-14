package kmeans

// import (
// 	pkgContext "context"
// 	"github.com/boson-research/patterns/internal/context"
// 	"reflect"
// 	"slices"
// 	"testing"
// )

// var ctx = context.New(pkgContext.Background())

// func Test_calculateSilhouetteScore(t *testing.T) {
// 	type args struct {
// 		data   []float64
// 		labels []int
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantMin float64
// 	}{
// 		{
// 			name: "4 clusters 4 centroids",
// 			args: args{
// 				data:   []float64{1, 2, 4, 5, 7, 8, 10},
// 				labels: []int{2, 2, 3, 3, 1, 1, 0},
// 			},
// 			wantMin: 0.5,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := calculateSilhouetteScore(tt.args.data, tt.args.labels); got < tt.wantMin {
// 				t.Errorf("calculateSilhouetteScore() = %v, want %v", got, tt.wantMin)
// 			}
// 		})
// 	}
// }

// func TestKMeans(t *testing.T) {

// 	type args struct {
// 		ctx  context.Context
// 		data []float64
// 	}
// 	tests := []struct {
// 		name                  string
// 		args                  args
// 		wantLabelCountsSorted []int
// 	}{
// 		{
// 			name: "1 cluster 1 element",
// 			args: args{
// 				ctx:  ctx,
// 				data: []float64{1},
// 			},
// 			wantLabelCountsSorted: []int{1},
// 		},
// 		{
// 			name: "2 clusters 1 distance",
// 			args: args{
// 				ctx:  ctx,
// 				data: []float64{1, 2, 3, 5, 6, 7},
// 			},
// 			wantLabelCountsSorted: []int{3, 3},
// 		},
// 		{
// 			name: "2 clusters 2 distance different sizes",
// 			args: args{
// 				ctx:  ctx,
// 				data: []float64{1, 2, 3, 6, 7},
// 			},
// 			wantLabelCountsSorted: []int{2, 3},
// 		},
// 		{
// 			name: "2 clusters 3 distance different sizes",
// 			args: args{
// 				ctx:  ctx,
// 				data: []float64{1, 2, 3, 7, 8},
// 			},
// 			wantLabelCountsSorted: []int{2, 3},
// 		},
// 		{
// 			name: "3 clusters 1 distance",
// 			args: args{
// 				ctx:  ctx,
// 				data: []float64{1, 2, 3, 5, 6, 7, 9, 10, 11},
// 			},
// 			wantLabelCountsSorted: []int{3, 3, 3},
// 		},
// 		{
// 			name: "4 clusters different distance different sizes",
// 			args: args{
// 				ctx:  ctx,
// 				data: []float64{1, 2, 3, 5, 6, 8, 9, 11, 12},
// 			},
// 			wantLabelCountsSorted: []int{2, 2, 2, 3},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, gotLabels := KMeans(tt.args.ctx, tt.args.data)

// 			gotLabelCountsSorted := labelsToCountsSorted(gotLabels)

// 			if !reflect.DeepEqual(gotLabelCountsSorted, tt.wantLabelCountsSorted) {
// 				t.Errorf("KMeans() gotLabelCountsSorted = %v, wantLabelCountsSorted %v", gotLabelCountsSorted, tt.wantLabelCountsSorted)
// 			}
// 		})
// 	}
// }

// func labelsToCountsSorted(labels []int) []int {
// 	counts := make(map[int]int)
// 	for _, label := range labels {
// 		counts[label]++
// 	}
// 	countsSorted := make([]int, 0, len(counts))
// 	for _, count := range counts {
// 		countsSorted = append(countsSorted, count)
// 	}
// 	slices.Sort(countsSorted)
// 	return countsSorted
// }
