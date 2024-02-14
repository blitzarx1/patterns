package cluster

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_generateOptimizationParamsVariations(t *testing.T) {
	type args struct {
		startingParams []int
		validator      func(params []int) error
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "test 1",
			args: args{
				startingParams: []int{1, 0},
				validator: func(params []int) error {
					if params[0] >= 3 || params[1] >= 3 {
						return fmt.Errorf("params must be less than 3")
					}

					return nil
				},
			},
			want: [][]int{
				{1, 0},
				{1, 1},
				{1, 2},
				{2, 0},
				{2, 1},
				{2, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateOptimizationParamsVariations(tt.args.startingParams, tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateOptimizationParamsVariations() = %v, want %v", got, tt.want)
			}
		})
	}
}
