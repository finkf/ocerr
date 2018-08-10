package main

import (
	"bufio"
	"fmt"
	"io"

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
		b, err := makeBlock(buf)
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

func makeBlock(buf []string) (block, error) {
	if len(buf) > maxBlockLines || len(buf) < minBlockLines {
		return block{}, fmt.Errorf("invalid block: %v", buf)
	}
	var b block
	if len(buf) == maxBlockLines {
		b.fn = buf[0]
		buf = buf[1:]
	}
	a, err := lev.NewAlignment(buf[0], buf[2], buf[1])
	b.a = a
	return b, err
}

func printBlock(b block, w io.Writer) error {
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
		string(b.a.S1), trace, string(b.a.S2), endOfBlock)
	return err
}
