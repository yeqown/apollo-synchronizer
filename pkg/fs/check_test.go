package fs

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type checkTestSuite struct {
	suite.Suite
}

func (c checkTestSuite) SetupSuite() {
	// mkdir -p ./testdata/check
	quickMkdir("./testdata/check")
	// touch ./testdata/check/file1
	quickTouch("./testdata/check/file1")
	// mkdir ./testdata/check/dir1
	quickMkdir("./testdata/check/dir1")
}

func (c checkTestSuite) TearDownSuite() {
	quickRemove("./testdata")
}

func (c checkTestSuite) TestCheck() {
	// check ./testdata/check/file1
	c.NotNil(MakeSure("./testdata/check/file1"))
	// check ./testdata/check/dir1
	c.NoError(MakeSure("./testdata/check/dir1"))
	c.NoError(MakeSure("./testdata/check/dir2"))
}

func Test_MakeSure(t *testing.T) {
	suite.Run(t, new(checkTestSuite))
}
