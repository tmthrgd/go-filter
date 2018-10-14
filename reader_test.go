package filter

import (
	"io/ioutil"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderRead(t *testing.T) {
	r := strings.NewReader(`this is
a very
normal test` + "\r" + `
case for
this package`)

	rr := NewReader(r, Not(Any(ContainsString("for"), Previous(HasSuffixString("is")))))

	b, err := ioutil.ReadAll(iotest.OneByteReader(rr))
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

func TestReaderLastLineNoNL(t *testing.T) {
	r := strings.NewReader("this is\na test")

	rr := NewReader(r, Not(HasSuffixString("test")))

	b, err := ioutil.ReadAll(iotest.OneByteReader(rr))
	require.NoError(t, err)

	assert.Equal(t, "this is\n", string(b))
}
