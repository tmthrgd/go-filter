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

func TestReaderRead(t *testing.T) {
	r := strings.NewReader(`this is
a very
normal test` + "\r" + `
case for
this package`)

	rr := NewReader(r, Not(Any(ContainsString("for"), Previous(HasSuffixString("is")))))

	b, err := ioutil.ReadAll(&byteReader{rr})
	require.NoError(t, err)

	assert.Equal(t, `this is
normal test`+"\r"+`
this package`, string(b))
}

func TestReaderWriteTo(t *testing.T) {
	r := strings.NewReader(`this is
a very
normal test` + "\r" + `
case for
this package`)

	rr := NewReader(r, Not(Any(ContainsString("for"), Previous(HasSuffixString("is")))))

	var s strings.Builder
	n, err := rr.WriteTo(&s)
	require.NoError(t, err)
	assert.Equal(t, int64(s.Len()), n, "WriteTo returned wrong bytes written count")

	assert.Equal(t, `this is
normal test`+"\r"+`
this package`, s.String())
}
