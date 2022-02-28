package asy

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_synchronizer_compare(t *testing.T) {
	type args struct {
		mode      SynchronizeMode
		force     bool
		overwrite bool
		parent    string
		local     []string
		remote    []string
	}
	tests := []struct {
		name string
		args args
		want []Diff1
	}{
		{
			name: "case1",
			args: args{
				mode:      SynchronizeMode_UPLOAD,
				force:     true,
				overwrite: true,
				parent:    "/tmp",
				local: []string{
					"app1.json",
					"app2.json",
				},
				remote: []string{
					"app1.json",
					"app3.json",
				},
			},
			want: []Diff1{
				{
					Diff0: Diff0{
						Key:  "app1.json",
						Mode: DiffMode_MODIFY,
					},
					AbsFilepath: "/tmp/app1.json",
				},
				{
					Diff0: Diff0{
						Key:  "app2.json",
						Mode: DiffMode_CREATE,
					},
					AbsFilepath: "/tmp/app2.json",
				},
				{
					Diff0: Diff0{
						Key:  "app3.json",
						Mode: DiffMode_DELETE,
					},
					AbsFilepath: "/tmp/app3.json",
				},
			},
		},
		{
			name: "case2",
			args: args{
				mode:      SynchronizeMode_DOWNLOAD,
				force:     true,
				overwrite: true,
				parent:    "/tmp",
				local: []string{
					"app1.json",
					"app2.json",
				},
				remote: []string{
					"app1.json",
					"app3.json",
				},
			},
			want: []Diff1{
				{
					Diff0: Diff0{
						Key:  "app3.json",
						Mode: DiffMode_CREATE,
					},
					AbsFilepath: "/tmp/app3.json",
				},
				{
					Diff0: Diff0{
						Key:  "app1.json",
						Mode: DiffMode_MODIFY,
					},
					AbsFilepath: "/tmp/app1.json",
				},
				{
					Diff0: Diff0{
						Key:  "app2.json",
						Mode: DiffMode_DELETE,
					},
					AbsFilepath: "/tmp/app2.json",
				},
			},
		},
		{
			name: "case3",
			args: args{
				mode:      SynchronizeMode_DOWNLOAD,
				force:     false,
				overwrite: true,
				parent:    "/tmp",
				local: []string{
					"app1.json",
					"app2.json",
				},
				remote: []string{
					"app1.json",
					"app3.json",
				},
			},
			want: []Diff1{
				{
					Diff0: Diff0{
						Key:  "app1.json",
						Mode: DiffMode_MODIFY,
					},
					AbsFilepath: "/tmp/app1.json",
				},
			},
		},
		{
			name: "case2",
			args: args{
				mode:      SynchronizeMode_DOWNLOAD,
				force:     true,
				overwrite: false,
				parent:    "/tmp",
				local: []string{
					"app1.json",
					"app2.json",
				},
				remote: []string{
					"app1.json",
					"app3.json",
				},
			},
			want: []Diff1{
				{
					Diff0: Diff0{
						Key:  "app3.json",
						Mode: DiffMode_CREATE,
					},
					AbsFilepath: "/tmp/app3.json",
				},
				{
					Diff0: Diff0{
						Key:  "app2.json",
						Mode: DiffMode_DELETE,
					},
					AbsFilepath: "/tmp/app2.json",
				},
			},
		},
		{
			name: "case2",
			args: args{
				mode:      SynchronizeMode_DOWNLOAD,
				force:     false,
				overwrite: false,
				parent:    "/tmp",
				local: []string{
					"app1.json",
					"app2.json",
				},
				remote: []string{
					"app1.json",
					"app3.json",
				},
			},
			want: []Diff1{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (synchronizer{}).compare(
				tt.args.mode,
				tt.args.parent,
				tt.args.force,
				tt.args.overwrite,
				tt.args.local,
				tt.args.remote,
			)

			sort.Sort(sorterDiff1(got))
			sort.Sort(sorterDiff1(tt.want))
			assert.Equal(t, tt.want, got)
		})
	}
}
