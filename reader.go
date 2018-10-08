// Package filter provides a filtered line-based io.Reader.
package filter

import (
	"bufio"
	"bytes"
	"io"
)

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

// dropCRLF drops a terminal \r\n from the data.
func dropCRLF(data []byte) []byte {
	data = data[:len(data)-1] // data will always be terminated with a \n

	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[:len(data)-1]
	}

	return data
}

// Reader is a filtered line-based io.Reader.
type Reader struct {
	s *bufio.Scanner
	f Func

	pos int
}

// Read implements io.Reader and reads from the underlying reader.
func (r *Reader) Read(p []byte) (n int, err error) {
	if b := r.s.Bytes(); r.pos >= 0 && r.pos < len(b) {
		n := copy(p, b[r.pos:])
		r.pos += n
		return n, nil
	}

	for r.s.Scan() {
		b := r.s.Bytes()
		if !r.f(dropCRLF(b)) {
			continue
		}

		n := copy(p, b)
		r.pos = n
		return n, nil
	}

	err = r.s.Err()
	if err == nil {
		err = io.EOF
	}

	return 0, err
}

// WriteTo implements io.WriterTo and reads from the underlying reader.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	if r.pos != -1 {
		panic("go-filter: inconsistent usage of io.Reader and io.WriterTo")
	}

	for r.s.Scan() {
		b := r.s.Bytes()
		if !r.f(dropCRLF(b)) {
			continue
		}

		nn, err := w.Write(b)
		if nn > len(b) {
			panic("go-filter: invalid Write count")
		}
		n += int64(nn)
		if err != nil {
			return n, err
		}
		if nn != len(b) {
			return n, io.ErrShortWrite
		}
	}

	return n, r.s.Err()
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
// line and preserves newlines and carriage returns.
func NewReader(r io.Reader, f Func) *Reader {
	s := bufio.NewScanner(r)
	s.Split(scanLines)

	return &Reader{
		s: s,
		f: f,

		pos: -1,
	}
}
