package filter

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type byteReader struct{ io.Reader }

func (br *byteReader) Read(p []byte) (int, error) {
	return br.Reader.Read(p[:1])
}

func TestReader(t *testing.T) {
	r := strings.NewReader(`this is
a very
normal test
case for
this pacakge`)

	rr := NewReader(r, Not(Or(ContainsString("for"), Previous(ContainsString("is")))))

	b, err := ioutil.ReadAll(&byteReader{rr})
	require.NoError(t, err)

	assert.Equal(t, `this is
normal test
this pacakge`, string(b))
}
