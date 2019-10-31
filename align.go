package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/finkf/lev"
	"github.com/spf13/cobra"
)

var (
	alignCmd = cobra.Command{
		Use:   "align",
		Short: `Align pairs of input lines`,
		Long:  `Align pairs of input lines and convert them into alignment blocks.`,
		RunE:  runAlign,
		Args:  cobra.ExactArgs(0),
	}
	gocrFileName bool
)

func init() {
}

func runAlign(cmd *cobra.Command, args []string) error {
	return align(os.Stdin, os.Stdout)
}

func align(stdin io.Reader, stdout io.Writer) error {
	var l lev.Lev
	return readPairs(stdin, func(p1, p2, s1, s2 string) error {
		a, err := l.Alignment(l.Trace(s1, s2))
		if err != nil {
			return err
		}
		return block{p1: p1, p2: p2, a: a}.write(stdout)
	})
}

func readPairs(r io.Reader, f func(p1, p2, s1, s2 string) error) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		s1 := s.Text()
		if !s.Scan() {
			break
		}
		s2 := s.Text()
		var p1, p2 string
		if global.separator != "" {
			n1 := strings.Index(s1, global.separator)
			n2 := strings.Index(s2, global.separator)
			if n1 == -1 || n2 == -1 {
				return fmt.Errorf("missing separtor: %q", global.separator)
			}
			p1 = s1[0 : n1+len(global.separator)]
			p2 = s2[0 : n2+len(global.separator)]
			s1 = s1[n1+len(global.separator):]
			s2 = s2[n2+len(global.separator):]
		}
		if err := f(p1, p2, s1, s2); err != nil {
			return err
		}
	}
	return s.Err()
}
