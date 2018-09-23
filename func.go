package filter

import "bytes"

// Func represents a line-matching function.
type Func func(line []byte) bool

// Or matches the line if any of f match the current line.
func Or(f ...Func) Func {
	return func(line []byte) bool {
		for _, f := range f {
			if f(line) {
				return true
			}
		}

		return false
	}
}

// And matches the line if all of f match the current line.
func And(f ...Func) Func {
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
	var prev []byte
	return func(line []byte) bool {
		ok := prev != nil && f(prev)
		prev = line
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
