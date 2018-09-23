package filter

import (
	"bufio"
	"io"
)

// Reader is a filtered line-based io.Reader.
type Reader struct {
	s *bufio.Scanner
	f Func

	nl *bool

	pos int
}

func (r *Reader) copyAndAdvance(p []byte) int {
	n := copy(p, r.s.Bytes()[r.pos:])
	r.pos += n

	if r.pos == len(r.s.Bytes()) && (len(p)-n >= 1 || !*r.nl) {
		if *r.nl {
			p[n] = '\n'
			n++
		}

		r.pos = -1
	}

	return n
}

// Read implements io.Reader and reads from the underlying reader.
func (r *Reader) Read(p []byte) (int, error) {
	if r.pos >= 0 {
		return r.copyAndAdvance(p), nil
	}

	for r.s.Scan() {
		if !r.f(r.s.Bytes()) {
			continue
		}

		r.pos = 0
		return r.copyAndAdvance(p), nil
	}

	err := r.s.Err()
	if err == nil {
		err = io.EOF
	}

	return 0, err
}

// Buffer sets the initial buffer to use when scanning and the maximum
// size of buffer that may be allocated during scanning. The maximum
// token size is the larger of max and cap(buf). If max <= cap(buf),
// Read will use this buffer only and do no allocation.
//
// By default, Read uses an internal buffer and sets the
// maximum token size to bufio.MaxScanTokenSize.
//
// Buffer panics if it is called after reading has started.
func (r *Reader) Buffer(buf []byte, max int) {
	r.s.Buffer(buf, max)
}

// NewReader wraps r and returns a new reader that will only pass
// through reads where the line is matched by f. It reads line by
// line and preserves newlines. It does not preserve carriage returns.
func NewReader(r io.Reader, f Func) *Reader {
	s := bufio.NewScanner(r)

	nl := new(bool)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		*nl = advance > 1 && data[advance-1] == '\n'
		return
	})

	return &Reader{
		s: s,
		f: f,

		nl: nl,

		pos: -1,
	}
}
