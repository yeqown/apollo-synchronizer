package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compare(t *testing.T) {
	type args struct {
		from []string
		to   []string
	}
	tests := []struct {
		name string
		args args
		want []diff0
	}{
		{
			name: "case1",
			args: args{
				from: []string{
					"app1.json",
					"app2.json",
				},
				to: []string{
					"app1.json",
				},
			},
			want: []diff0{
				{
					key:  "app1.json",
					mode: diffMode_MODIFY,
				},
				{
					key:  "app2.json",
					mode: diffMode_CREATE,
				},
			},
		},
		{
			name: "case2",
			args: args{
				from: []string{
					"app1.json",
				},
				to: []string{
					"app1.json",
					"app2.json",
				},
			},
			want: []diff0{
				{
					key:  "app1.json",
					mode: diffMode_MODIFY,
				},
				{
					key:  "app2.json",
					mode: diffMode_DELETE,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compare(tt.args.from, tt.args.to)
			assert.Equal(t, tt.want, got)
		})
	}
}
