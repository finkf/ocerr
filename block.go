package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/finkf/lev"
)

const nBlockLines = 3

type block struct {
	a             lev.Alignment
	p1, p2, stats string
}

func (b *block) init(buf [nBlockLines]string, sep string) error {
	// handle possible prefixes if a non empty separator is given
	if sep != "" {
		p1 := strings.Index(buf[0], sep)
		p2 := strings.Index(buf[2], sep)
		if p1 == -1 || p2 == -1 {
			return fmt.Errorf("missing separator: %q", sep)
		}
		n := max(p1, p2) // prefixes are justified; so use position
		b.p1 = buf[0][0 : p1+len(global.separator)]
		b.p2 = buf[2][0 : p2+len(global.separator)]
		buf[0] = buf[0][n+len(global.separator):]
		buf[1] = buf[1][n+len(global.separator):]
		buf[2] = buf[2][n+len(global.separator):]
	}
	a, err := lev.NewAlignment(
		deleteDottedCircles(buf[0]),
		deleteDottedCircles(buf[2]),
		buf[1],
	)
	b.a = a
	return err
}

func (b block) write(w io.Writer) error {
	n := maxlen(b.p1, b.p2)
	trace := string(b.a.Trace)
	if b.stats != "" {
		trace += " " + b.stats
	}
	_, err := fmt.Fprintf(w, "%s%s\n%s%s\n%s%s\n%s\n",
		prefix(b.p1, n),
		string(addDottedCircles(b.a.S1)),
		prefix("", n),
		trace,
		prefix(b.p2, n),
		string(addDottedCircles(b.a.S2)),
		global.endOfBlock,
	)
	return err
}

type blockScanner struct {
	scanner *bufio.Scanner
	_block  block
	_err    error
}

func newBlockScanner(in io.Reader) *blockScanner {
	return &blockScanner{scanner: bufio.NewScanner(in)}
}

func (s *blockScanner) scan() bool {
	var buf [nBlockLines]string
	var i int
	for s.scanner.Scan() {
		str := s.scanner.Text()
		if str != global.endOfBlock && i < nBlockLines {
			buf[i] = str
			i++
			continue
		}
		if i != nBlockLines {
			s._err = fmt.Errorf("invalid block")
			return false
		}
		s._err = s._block.init(buf, global.separator)
		if s._err != nil {
			return false
		}
		return true
	}
	if s._err == nil {
		s._err = s.scanner.Err()
	}
	return false
}

func (s *blockScanner) block() block {
	return s._block
}

func (s *blockScanner) err() error {
	return s._err
}

func readBlocks(r io.Reader, f func(block) error) error {
	s := newBlockScanner(r)
	for s.scan() {
		if err := f(s.block()); err != nil {
			return err
		}
	}
	return s.err()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxlen(a, b string) int {
	return max(len(a), len(b))
}

const (
	dottedCircle = '◌'
)

func addDottedCircles(rs []rune) []rune {
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if unicode.Is(unicode.Mn, r) {
			rs = append(rs[:i], append([]rune{dottedCircle}, rs[i:]...)...)
			i = i + 1
		}
	}
	return rs
}

func deleteDottedCircles(str string) string {
	return strings.Replace(str, string(dottedCircle), "", -1)
}

func prefix(str string, max int) string {
	if max > 0 {
		// left-justifiy and fill with spaces
		return fmt.Sprintf("%-*s", max, str)
	}
	return ""
}
