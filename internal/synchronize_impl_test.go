package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_synchronizer_compare(t *testing.T) {
	type args struct {
		mode       SynchronizeMode
		localFiles []string
		namespaces []string
	}
	tests := []struct {
		name string
		args args
		want []diff
	}{
		{
			name: "case1",
			args: args{
				mode: SynchronizeMode_UPLOAD,
				localFiles: []string{
					"app1.json",
					"app2.json",
				},
				namespaces: []string{
					"app1.json",
				},
			},
			want: []diff{
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
				mode: SynchronizeMode_UPLOAD,
				localFiles: []string{
					"app1.json",
				},
				namespaces: []string{
					"app1.json",
					"app2.json",
				},
			},
			want: []diff{
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
		{
			name: "case3",
			args: args{
				mode: SynchronizeMode_DOWNLOAD,
				localFiles: []string{
					"app1.json",
				},
				namespaces: []string{
					"app1.json",
					"app2.json",
				},
			},
			want: []diff{
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
			name: "case4",
			args: args{
				mode: SynchronizeMode_DOWNLOAD,
				localFiles: []string{
					"app1.json",
					"app2.json",
				},
				namespaces: []string{
					"app1.json",
				},
			},
			want: []diff{
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
			got := (synchronizer{}).compare(tt.args.mode, tt.args.localFiles, tt.args.namespaces)
			assert.Equal(t, tt.want, got)
		})
	}
}
