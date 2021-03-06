module github.com/yeqown/apollo-synchronizer

go 1.16

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/yeqown/log v1.1.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

retract (
	v1.3.3
	v1.3.2
)
