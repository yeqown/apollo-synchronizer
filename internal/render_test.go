package internal

import "testing"

var (
	_diffs = []diff1{
		{
			diff0: diff0{
				key:  "app.json",
				mode: diffMode_MODIFY,
			},
			absFilepath: "/var/apollo-synchronizer/app",
		},
		{
			diff0: diff0{
				key:  "app2.json",
				mode: diffMode_CREATE,
			},
			absFilepath: "/var/apollo-synchronizer/app",
		},
		{
			diff0: diff0{
				key:  "app3.json",
				mode: diffMode_DELETE,
			},
			absFilepath: "/var/apollo-synchronizer/app",
		},
	}

	_results = []*synchronizeResult{
		{
			key:       "app.json",
			mode:      diffMode_MODIFY,
			error:     "you failed",
			succeeded: false,
			published: false,
		},
		{
			key:       "app2.json",
			mode:      diffMode_CREATE,
			error:     "",
			succeeded: true,
			published: false,
		},
		{
			key:       "app3.json",
			mode:      diffMode_DELETE,
			error:     "",
			succeeded: true,
			published: true,
		},
	}
)

func Test_renderer_terminal(t *testing.T) {
	r := terminalRenderer{}

	r.renderingDiffs(_diffs)
	r.renderingResult(_results)
}

func Test_renderer_termui(t *testing.T) {
	r := newTermUI(&SynchronizeScope{
		Mode:              0,
		Path:              "",
		LocalFiles:        nil,
		ApolloSecret:      "",
		ApolloAppID:       "",
		ApolloEnv:         "",
		ApolloClusterName: "",
		ApolloPortalAddr:  "",
		ApolloAccount:     "",
		ApolloAutoPublish: false,
		Overwrite:         false,
		Force:             false,
	})

	r.renderingDiffs(_diffs)
	r.renderingResult(_results)
}
