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

func TestAny(t *testing.T) {
	assert.True(t, Any(True, True, True)(test), "Any(False, True, False)")
	assert.True(t, Any(False, True, False)(test), "Any(False, True, False)")
	assert.True(t, Any(True)(test), "Any(True)")
	assert.False(t, Any(False)(test), "Any(False)")
	assert.False(t, Any()(test), "Any()")
}

func TestAll(t *testing.T) {
	assert.True(t, All(True, True, True)(test), "All(False, True, False)")
	assert.False(t, All(False, True, False)(test), "All(False, True, False)")
	assert.True(t, All(True)(test), "All(True)")
	assert.False(t, All(False)(test), "All(False)")
	assert.False(t, All()(test), "All()")
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
	testFn := func(t *testing.T, typ MatchType) {
		f := Before(HasSuffixString("2"), typ)

		for i := 0; i < 3; i++ {
			assert.True(t, f(test), "before, not matching: i == %d", i)
		}

		assert.Equal(t, typ == IncludeCurrent, f(test2), "matching condition")

		for i := 0; i < 3; i++ {
			assert.False(t, f(test), "after, not matching: i == %d", i)
		}

		for i := 0; i < 3; i++ {
			assert.False(t, f(test2), "after, matching: i == %d", i)
		}
	}
	t.Run("exclude", func(t *testing.T) { testFn(t, ExcludeCurrent) })
	t.Run("include", func(t *testing.T) { testFn(t, IncludeCurrent) })
}

func TestAfter(t *testing.T) {
	testFn := func(t *testing.T, typ MatchType) {
		f := After(HasSuffixString("2"), typ)

		for i := 0; i < 3; i++ {
			assert.False(t, f(test), "before, not matching: i == %d", i)
		}

		assert.Equal(t, typ == IncludeCurrent, f(test2), "matching condition")

		for i := 0; i < 3; i++ {
			assert.True(t, f(test), "after, not matching: i == %d", i)
		}

		for i := 0; i < 3; i++ {
			assert.True(t, f(test2), "after, matching: i == %d", i)
		}
	}
	t.Run("exclude", func(t *testing.T) { testFn(t, ExcludeCurrent) })
	t.Run("include", func(t *testing.T) { testFn(t, IncludeCurrent) })
}
