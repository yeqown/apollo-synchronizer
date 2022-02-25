package asy

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
		want []Diff0
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
			want: []Diff0{
				{
					Key:  "app1.json",
					Mode: DiffMode_MODIFY,
				},
				{
					Key:  "app2.json",
					Mode: DiffMode_CREATE,
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
			want: []Diff0{
				{
					Key:  "app1.json",
					Mode: DiffMode_MODIFY,
				},
				{
					Key:  "app2.json",
					Mode: DiffMode_DELETE,
				},
			},
		},
		{
			name: "case 3",
			args: args{
				from: []string{
					"app1.json",
					"app2.json",
				},
				to: []string{},
			},
			want: []Diff0{
				{
					Key:  "app1.json",
					Mode: DiffMode_CREATE,
				},
				{
					Key:  "app2.json",
					Mode: DiffMode_CREATE,
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
