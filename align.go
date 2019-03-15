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
		Long:  `Align pairs of input lines and convert them into alignment blocks`,
		Short: `Align pairs of input lines`,
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
	return readAlignInput(stdin, func(p1, p2, s1, s2 string) error {
		a, err := l.Alignment(l.Trace(s1, s2))
		if err != nil {
			return err
		}
		return writeBlock(block{p1: p1, p2: p2, a: a}, stdout)
	})
}

func readAlignInput(r io.Reader, f func(p1, p2, s1, s2 string) error) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		s1 := s.Text()
		if !s.Scan() {
			break
		}
		s2 := s.Text()
		var p1, p2 string
		if separator != "" {
			n1 := strings.Index(s1, separator)
			n2 := strings.Index(s2, separator)
			if n1 == -1 || n2 == -1 {
				return fmt.Errorf("missing separtor: %q", separator)
			}
			p1 = s1[0 : n1+len(separator)]
			p2 = s2[0 : n2+len(separator)]
			s1 = s1[n1+len(separator):]
			s2 = s2[n2+len(separator):]
		}
		if err := f(p1, p2, s1, s2); err != nil {
			return err
		}
	}
	return s.Err()
}
