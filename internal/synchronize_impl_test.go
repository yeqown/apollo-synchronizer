package internal

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
		want []diff1
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
			want: []diff1{
				{
					diff0: diff0{
						key:  "app1.json",
						mode: diffMode_MODIFY,
					},
					absFilepath: "/tmp/app1.json",
				},
				{
					diff0: diff0{
						key:  "app2.json",
						mode: diffMode_CREATE,
					},
					absFilepath: "/tmp/app2.json",
				},
				{
					diff0: diff0{
						key:  "app3.json",
						mode: diffMode_DELETE,
					},
					absFilepath: "/tmp/app3.json",
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
			want: []diff1{
				{
					diff0: diff0{
						key:  "app3.json",
						mode: diffMode_CREATE,
					},
					absFilepath: "/tmp/app3.json",
				},
				{
					diff0: diff0{
						key:  "app1.json",
						mode: diffMode_MODIFY,
					},
					absFilepath: "/tmp/app1.json",
				},
				{
					diff0: diff0{
						key:  "app2.json",
						mode: diffMode_DELETE,
					},
					absFilepath: "/tmp/app2.json",
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
			want: []diff1{
				{
					diff0: diff0{
						key:  "app1.json",
						mode: diffMode_MODIFY,
					},
					absFilepath: "/tmp/app1.json",
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
			want: []diff1{
				{
					diff0: diff0{
						key:  "app3.json",
						mode: diffMode_CREATE,
					},
					absFilepath: "/tmp/app3.json",
				},
				{
					diff0: diff0{
						key:  "app2.json",
						mode: diffMode_DELETE,
					},
					absFilepath: "/tmp/app2.json",
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
			want: []diff1{},
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
