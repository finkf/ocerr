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

type readBlocksFunc func(block) error

func readBlocks(r io.Reader, f readBlocksFunc) error {
	s := bufio.NewScanner(r)
	buf := make([]string, 0, nBlockLines)
	for s.Scan() {
		str := s.Text()
		if str != global.endOfBlock {
			buf = append(buf, str)
			continue
		}
		b, err := newBlock(buf)
		if err != nil {
			return err
		}
		if err := f(b); err != nil {
			return err
		}
		buf = buf[:0]
	}
	return s.Err()
}

func newBlock(buf []string) (block, error) {
	if len(buf) != nBlockLines {
		return block{}, fmt.Errorf("invalid block: %v", buf)
	}
	var b block
	// handle possible prefixes if a non empty separator is given
	if global.separator != "" {
		p1 := strings.Index(buf[0], global.separator)
		p2 := strings.Index(buf[2], global.separator)
		if p1 == -1 || p2 == -1 {
			return block{}, fmt.Errorf("missing separator: %q", global.separator)
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
	return b, err
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

func writeBlock(b block, w io.Writer) error {
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
