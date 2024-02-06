package processor

import (
	// "reflect"
	"testing"

	"github.com/boson-research/patterns/internal/models"
	// "github.com/boson-research/patterns/internal/models/neighbourhood"
)

func Test_checkPattern(t *testing.T) {
	type args struct {
		pattern *models.Pattern
		text    []byte
		it      int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "out of bounds",
			args: args{
				pattern: models.NewPattern([]byte{'a', 'a', 'a'}),
				text:    []byte("a"),
				it:      0,
			},
			want: false,
		},
		{
			name: "out of bounds 2",
			args: args{
				pattern: models.NewPattern([]byte{'a', 'a', 'a'}),
				text:    []byte("aaa"),
				it:      1,
			},
			want: false,
		},
		{
			name: "success",
			args: args{
				pattern: models.NewPattern([]byte{'a', 'b', 'c'}),
				text:    []byte("dabc"),
				it:      1,
			},
			want: true,
		},
		{
			name: "fail",
			args: args{
				pattern: models.NewPattern([]byte{'a', 'b', 'c'}),
				text:    []byte("dabc"),
				it:      0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPattern(tt.args.pattern, tt.args.text, tt.args.it); got != tt.want {
				t.Errorf("Processor.checkPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_mergeStatsNeighbourhoods(t *testing.T) {
// 	type args struct {
// 		a *neighbourhood.Stat
// 		b *neighbourhood.Stat
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *neighbourhood.Neighbourhood
// 	}{
// 		{
// 			name: "a is empty",
// 			args: args{
// 				a: neighbourhood.NewStatNeighbourhood(),
// 				b: &neighbourhood.Stat{
// 					locations: []int{1, 2, 3},
// 					patterns: []*models.Pattern{
// 						models.NewPattern([]byte("abc")),
// 						models.NewPattern([]byte("def")),
// 						models.NewPattern([]byte("ghi")),
// 					},
// 				},
// 			},
// 			want: &neighbourhood.Stat{
// 				locations: []int{1, 2, 3},
// 				patterns: []*models.Pattern{
// 					models.NewPattern([]byte("abc")),
// 					models.NewPattern([]byte("def")),
// 					models.NewPattern([]byte("ghi")),
// 				},
// 			},
// 		},
// 		{
// 			name: "b is empty",
// 			args: args{
// 				a: &neighbourhood.Stat{
// 					locations: []int{1, 2, 3},
// 					patterns: []*models.Pattern{
// 						models.NewPattern([]byte("abc")),
// 						models.NewPattern([]byte("def")),
// 						models.NewPattern([]byte("ghi")),
// 					},
// 				},
// 				b: neighbourhood.NewStatNeighbourhood(),
// 			},
// 			want: &neighbourhood.Stat{
// 				locations: []int{1, 2, 3},
// 				patterns: []*models.Pattern{
// 					models.NewPattern([]byte("abc")),
// 					models.NewPattern([]byte("def")),
// 					models.NewPattern([]byte("ghi")),
// 				},
// 			},
// 		},
// 		{
// 			name: "a and b are nil",
// 			args: args{
// 				a: nil,
// 				b: nil,
// 			},
// 			want: nil,
// 		},
// 		{
// 			name: "a is nil and b is empty",
// 			args: args{
// 				a: nil,
// 				b: neighbourhood.NewStatNeighbourhood(),
// 			},
// 			want: neighbourhood.NewStatNeighbourhood(),
// 		},
// 		{
// 			name: "merge",
// 			args: args{
// 				a: &neighbourhood.Stat{
// 					locations: []int{1, 5, 10, 12},
// 					patterns: []*models.Pattern{
// 						models.NewPattern([]byte("1")),
// 						models.NewPattern([]byte("5")),
// 						models.NewPattern([]byte("10")),
// 						models.NewPattern([]byte("12")),
// 					},
// 				},
// 				b: &neighbourhood.Stat{
// 					locations: []int{2, 3, 4, 8, 11},
// 					patterns: []*models.Pattern{
// 						models.NewPattern([]byte("2")),
// 						models.NewPattern([]byte("3")),
// 						models.NewPattern([]byte("4")),
// 						models.NewPattern([]byte("8")),
// 						models.NewPattern([]byte("11")),
// 					},
// 				},
// 			},
// 			want: &neighbourhood.Stat{
// 				locations: []int{1, 2, 3, 4, 5, 8, 10, 11, 12},
// 				patterns: []*models.Pattern{
// 					models.NewPattern([]byte("1")),
// 					models.NewPattern([]byte("2")),
// 					models.NewPattern([]byte("3")),
// 					models.NewPattern([]byte("4")),
// 					models.NewPattern([]byte("5")),
// 					models.NewPattern([]byte("8")),
// 					models.NewPattern([]byte("10")),
// 					models.NewPattern([]byte("11")),
// 					models.NewPattern([]byte("12")),
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := mergeStatsNeighbourhoods(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MergeStatsNeighbourhoods() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
