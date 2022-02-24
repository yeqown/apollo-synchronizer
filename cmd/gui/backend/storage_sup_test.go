package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveReadInBinary(t *testing.T) {
	err := save("./testdata/bin", []byte("test"), _ext_binary)
	assert.NoError(t, err)

	out := make([]byte, 4)
	err = read("./testdata/bin", &out, _ext_binary)
	assert.NoError(t, err)

	assert.Equal(t, []byte("test"), out)
}
