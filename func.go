package filter

import "bytes"

// Func represents a line-matching function.
type Func func(line []byte) bool

// Any matches the line if any of f match the current line.
func Any(f ...Func) Func {
	return func(line []byte) bool {
		for _, f := range f {
			if f(line) {
				return true
			}
		}

		return false
	}
}

// All matches the line if all of f match the current line.
func All(f ...Func) Func {
	return func(line []byte) bool {
		for _, f := range f {
			if !f(line) {
				return false
			}
		}

		return len(f) > 0
	}
}

// Not inverts f and matches the current line if f doesn't match the line.
func Not(f Func) Func {
	return func(line []byte) bool {
		return !f(line)
	}
}

// Previous matches the current line iff f matches the previous line.
// It never matches the first line.
//
// It is not safe to call concurrently or reuse.
func Previous(f Func) Func {
	var next bool
	return func(line []byte) bool {
		ok := next
		next = f(line)
		return ok
	}
}

// Contains matches the current line if it contains tok.
func Contains(tok []byte) Func {
	return func(line []byte) bool {
		return bytes.Contains(line, tok)
	}
}

// ContainsString matches the current line if it contains tok.
func ContainsString(tok string) Func {
	return Contains([]byte(tok))
}

// HasPrefix matches the current line if it starts with prefix.
func HasPrefix(prefix []byte) Func {
	return func(line []byte) bool {
		return bytes.HasPrefix(line, prefix)
	}
}

// HasPrefixString matches the current line if it starts with prefix.
func HasPrefixString(prefix string) Func {
	return HasPrefix([]byte(prefix))
}

// HasSuffix matches the current line if it ends with suffix.
func HasSuffix(suffix []byte) Func {
	return func(line []byte) bool {
		return bytes.HasSuffix(line, suffix)
	}
}

// HasSuffixString matches the current line if it ends with suffix.
func HasSuffixString(suffix string) Func {
	return HasSuffix([]byte(suffix))
}

// Odd matches every second line including the first.
//
// It is not safe to call concurrently or reuse.
func Odd() Func {
	var n int
	return func(line []byte) bool {
		n++
		return n%2 == 1
	}
}

// Even matches every second line excluding the first.
//
// It is not safe to call concurrently or reuse.
func Even() Func {
	var n int
	return func(line []byte) bool {
		n++
		return n%2 == 0
	}
}

// MatchType controls whether Before and After match the current line.
type MatchType bool

const (
	// ExcludeCurrent instructs the matching function to exclude
	// the current line.
	ExcludeCurrent MatchType = true
	// IncludeCurrent instructs the matching function to include
	// the current line.
	IncludeCurrent MatchType = false
)

// Before matches every line before f matches the current line.
//
// Whether to include or exclude the current line is controlled by
// typ. Either ExcludeCurrent or IncludeCurrent may be passed.
//
// It is not safe to call concurrently or reuse.
func Before(f Func, typ MatchType) Func {
	ok := true
	return func(line []byte) bool {
		old := ok
		ok = ok && !f(line)

		if typ == IncludeCurrent {
			return old
		}

		return ok
	}
}

// After matches every line after f matches the current line.
//
// Whether to include or exclude the current line is controlled by
// typ. Either ExcludeCurrent or IncludeCurrent may be passed.
//
// It is not safe to call concurrently or reuse.
func After(f Func, typ MatchType) Func {
	var ok bool
	return func(line []byte) bool {
		old := ok
		ok = ok || f(line)

		if typ == ExcludeCurrent {
			return old
		}

		return ok
	}
}
