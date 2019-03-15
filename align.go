package main

import (
	"bufio"
	"io"
	"os"

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
	alignCmd.Flags().BoolVarP(&gocrFileName, "header", "H",
		false, "read the filename as additional first line from input")
}

func runAlign(cmd *cobra.Command, args []string) error {
	return align(os.Stdin, os.Stdout)
}

func align(stdin io.Reader, stdout io.Writer) error {
	var l lev.Lev
	return readAlignInput(stdin, func(fn, s1, s2 string) error {
		a, err := l.Alignment(l.Trace(s1, s2))
		if err != nil {
			return err
		}
		return writeBlock(block{fn: fn, a: a}, stdout)
	})
}

type readAlignInputFunc func(string, string, string) error

func readAlignInput(r io.Reader, f readAlignInputFunc) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		var fn string
		if gocrFileName {
			fn = s.Text()
			if !s.Scan() {
				break
			}
		}
		s1 := s.Text()
		if !s.Scan() {
			break
		}
		s2 := s.Text()
		if err := f(fn, s1, s2); err != nil {
			return err
		}
	}
	return s.Err()
}
