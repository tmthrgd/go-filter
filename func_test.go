package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func True(line []byte) bool  { return true }
func False(line []byte) bool { return false }

var (
	test  = []byte("test")
	test2 = []byte("test2")
)

func TestOr(t *testing.T) {
	assert.True(t, Or(True, True, True)(test), "Or(False, True, False)")
	assert.True(t, Or(False, True, False)(test), "Or(False, True, False)")
	assert.True(t, Or(True)(test), "Or(True)")
	assert.False(t, Or(False)(test), "Or(False)")
	assert.False(t, Or()(test), "Or()")
}

func TestAnd(t *testing.T) {
	assert.True(t, And(True, True, True)(test), "And(False, True, False)")
	assert.False(t, And(False, True, False)(test), "And(False, True, False)")
	assert.True(t, And(True)(test), "And(True)")
	assert.False(t, And(False)(test), "And(False)")
	assert.False(t, And()(test), "And()")
}

func TestNot(t *testing.T) {
	assert.False(t, Not(True)(test), "Not(True)")
	assert.True(t, Not(False)(test), "Not(False)")
}

func TestPrevious(t *testing.T) {
	f := Previous(HasSuffixString("2"))
	assert.False(t, f(test2))
	assert.True(t, f(test))
	assert.False(t, f(test2))
	assert.True(t, f(test2))
	assert.True(t, f(test))
}

func TestContains(t *testing.T) {
	assert.True(t, Contains(test)(test2), "Contains(test)(test2)")
	assert.True(t, Contains(test2)(test2), "Contains(test)(test2)")
	assert.False(t, Contains(test2)(test), "Contains(test)(test2)")

	assert.True(t, ContainsString("test")(test2), `ContainsString("test")(test2)`)
	assert.True(t, ContainsString("test2")(test2), `ContainsString("test")(test2)`)
	assert.False(t, ContainsString("test2")(test), `ContainsString("test")(test2)`)
}

func TestHasPrefix(t *testing.T) {
	assert.True(t, HasPrefix(test)(test2), "HasPrefix(test)(test2)")
	assert.True(t, HasPrefix(test2)(test2), "HasPrefix(test)(test2)")
	assert.False(t, HasPrefix(test2)(test), "HasPrefix(test)(test2)")

	assert.True(t, HasPrefixString("test")(test2), `HasPrefixString("test")(test2)`)
	assert.True(t, HasPrefixString("test2")(test2), `HasPrefixString("test")(test2)`)
	assert.False(t, HasPrefixString("test2")(test), `HasPrefixString("test")(test2)`)
}

func TestHasSuffix(t *testing.T) {
	two := []byte("2")

	assert.True(t, HasSuffix(test)(test), "HasSuffix(test)(test)")
	assert.True(t, HasSuffix(test2)(test2), "HasSuffix(test)(test2)")
	assert.False(t, HasSuffix(test2)(test), "HasSuffix(test)(test2)")
	assert.True(t, HasSuffix(two)(test2), "HasSuffix(two)(test2)")
	assert.False(t, HasSuffix(two)(test), "HasSuffix(two)(test)")

	assert.True(t, HasSuffixString("test")(test), `HasSuffixString("test")(test)`)
	assert.True(t, HasSuffixString("test2")(test2), `HasSuffixString("test")(test2)`)
	assert.False(t, HasSuffixString("test2")(test), `HasSuffixString("test")(test2)`)
	assert.True(t, HasSuffixString("2")(test2), `HasSuffixString("2")(test2)`)
	assert.False(t, HasSuffixString("2")(test), `HasSuffixString("2")(test)`)
}

func TestOdd(t *testing.T) {
	f := Odd()

	for i := 0; i < 10; i++ {
		assert.True(t, f(test), "odd")
		assert.False(t, f(test), "even")
	}
}

func TestEven(t *testing.T) {
	f := Even()

	for i := 0; i < 10; i++ {
		assert.False(t, f(test), "odd")
		assert.True(t, f(test), "even")
	}
}

func TestBefore(t *testing.T) {
	f := Before(HasSuffixString("2"))

	for i := 0; i < 3; i++ {
		assert.True(t, f(test), "before, not matching: i == %d", i)
	}

	assert.False(t, f(test2), "matching condition")

	for i := 0; i < 3; i++ {
		assert.False(t, f(test), "after, not matching: i == %d", i)
	}

	for i := 0; i < 3; i++ {
		assert.False(t, f(test2), "after, matching: i == %d", i)
	}
}

func TestAfter(t *testing.T) {
	f := After(HasSuffixString("2"))

	for i := 0; i < 3; i++ {
		assert.False(t, f(test), "before, not matching: i == %d", i)
	}

	assert.True(t, f(test2), "matching condition")

	for i := 0; i < 3; i++ {
		assert.True(t, f(test), "after, not matching: i == %d", i)
	}

	for i := 0; i < 3; i++ {
		assert.True(t, f(test2), "after, matching: i == %d", i)
	}
}
