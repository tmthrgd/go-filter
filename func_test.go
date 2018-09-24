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

func TestPreviousBufferReuse(t *testing.T) {
	f := Previous(HasPrefixString("test"))
	buf := []byte("test buffer")

	assert.False(t, f(buf))
	copy(buf, []byte("nope"))
	assert.True(t, f(buf))
	assert.False(t, f(buf))
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

func TestAlternate(t *testing.T) {
	newF := func(ok bool) (Func, *bool) {
		wasCalled := new(bool)
		return func(line []byte) bool {
			*wasCalled = true
			return ok
		}, wasCalled
	}

	f1, c1 := newF(true)
	f2, c2 := newF(false)
	f3, c3 := newF(true)
	f4, c4 := newF(true)
	f5, c5 := newF(false)

	f := Alternate(f1, f2, f3, f4, f5)

	for i := 0; i < 3; i++ {
		assert.True(t, f(test))
		assert.True(t, *c1, "Func #1 should have been called (i == %d)", i)

		assert.False(t, f(test))
		assert.True(t, *c2, "Func #2 should have been called (i == %d)", i)

		assert.True(t, f(test))
		assert.True(t, *c3, "Func #3 should have been called (i == %d)", i)

		assert.True(t, f(test))
		assert.True(t, *c4, "Func #4 should have been called (i == %d)", i)

		assert.False(t, f(test))
		assert.True(t, *c5, "Func #5 should have been called (i == %d)", i)

		*c1, *c2, *c3, *c4, *c5 = false, false, false, false, false
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
