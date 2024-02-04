package processor

import (
	"context"
	"testing"

	"github.com/boson-research/patterns/internal/models"
)

func TestProcessor_checkPattern(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		pattern models.Pattern
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
				ctx:     ctx,
				pattern: models.Pattern{'a', 'a', 'a'},
				text:    []byte("a"),
				it:      0,
			},
			want: false,
		},
		{
			name: "out of bounds 2",
			args: args{
				ctx:     ctx,
				pattern: models.Pattern{'a', 'a', 'a'},
				text:    []byte("aaa"),
				it:      1,
			},
			want: false,
		},
		{
			name: "success",
			args: args{
				ctx:     ctx,
				pattern: models.Pattern{'a', 'b', 'c'},
				text:    []byte("dabc"),
				it:      1,
			},
			want: true,
		},
		{
			name: "fail",
			args: args{
				ctx:     ctx,
				pattern: models.Pattern{'a', 'b', 'c'},
				text:    []byte("dabc"),
				it:      0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPattern(tt.args.ctx, tt.args.pattern, tt.args.text, tt.args.it); got != tt.want {
				t.Errorf("Processor.checkPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
