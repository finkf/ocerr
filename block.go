package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/finkf/lev"
)

const (
	minBlockLines = 3
	maxBlockLines = 4
	endOfBlock    = ""
)

type block struct {
	a         lev.Alignment
	fn, stats string
}

type readBlocksFunc func(block) error

func readBlocks(r io.Reader, f readBlocksFunc) error {
	s := bufio.NewScanner(r)
	buf := make([]string, 0, maxBlockLines)
	for s.Scan() {
		str := s.Text()
		if str != endOfBlock {
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
	if len(buf) > maxBlockLines || len(buf) < minBlockLines {
		return block{}, fmt.Errorf("invalid block: %v", buf)
	}
	var b block
	if len(buf) == maxBlockLines {
		b.fn = buf[0]
		buf = buf[1:]
	}
	a, err := lev.NewAlignment(
		deleteDottedCircles(buf[0]),
		deleteDottedCircles(buf[2]),
		buf[1],
	)
	b.a = a
	return b, err
}

func writeBlock(b block, w io.Writer) error {
	if len(b.fn) > 0 {
		if _, err := fmt.Fprintln(w, b.fn); err != nil {
			return err
		}
	}
	trace := string(b.a.Trace)
	if b.stats != "" {
		trace += " " + b.stats
	}
	_, err := fmt.Fprintf(w, "%s\n%s\n%s\n%s\n",
		string(addDottedCircles(b.a.S1)),
		trace,
		string(addDottedCircles(b.a.S2)),
		endOfBlock,
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
