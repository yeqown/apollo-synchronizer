package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type fsTestSuite struct {
	suite.Suite
}

func quickMkdir(p string) {
	if err := os.MkdirAll(p, 0777); err != nil {
		panic(err)
	}
}

func quickTouch(f string) {
	fd, err := os.Create(f)
	if err != nil {
		panic(err)
	}
	fd.Close()
}

func quickRemove(p string) {
	os.RemoveAll(p)
}

func (f fsTestSuite) SetupSuite() {
	// mkdir -p ./testdata/deeper
	quickMkdir("./testdata/deeper")
	// touch ./testdata/deeper/file1.txt
	quickTouch("./testdata/deeper/file1.txt")
	// touch ./testdata/deeper/file2.txt
	quickTouch("./testdata/deeper/file2.txt")
	// touch ./testdata/file1.txt
	quickTouch("./testdata/file1.txt")
	// touch ./testdata/file2.txt
	quickTouch("./testdata/file2.txt")
	// touch ./testdata/file3.txt
	quickTouch("./testdata/file3.txt")
}

func (f fsTestSuite) TearDownSuite() {
	quickRemove("./testdata")
}

func (f fsTestSuite) Test_TravelDir_NoRecursive() {
	out, err := TravelDirectory("./testdata", false)
	f.NoError(err)

	f.Equal([]string{
		"testdata/file1.txt",
		"testdata/file2.txt",
		"testdata/file3.txt",
	}, out)

}

func (f fsTestSuite) Test_TravelDir_Recursive() {
	out, err := TravelDirectory("./testdata", true)
	f.NoError(err)

	f.Equal([]string{
		"testdata/deeper/file1.txt",
		"testdata/deeper/file2.txt",
		"testdata/file1.txt",
		"testdata/file2.txt",
		"testdata/file3.txt",
	}, out)
}

func Test_fs(t *testing.T) {
	suite.Run(t, new(fsTestSuite))
}
