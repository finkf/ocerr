package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	pairCmd = cobra.Command{
		Use:   "pair",
		Short: `Pairs multiple files from stdin to stdout`,
		RunE:  runPair,
	}
)

func init() {
}

func runPair(cmd *cobra.Command, args []string) error {
	return pair(os.Stdin, os.Stdout)
}

func pair(stdin io.Reader, stdout io.Writer) error {
	s := bufio.NewScanner(stdin)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		return s.Err()
	}
	n := len(lines) / 2
	for i := 0; i < n; i++ {
		if _, err := fmt.Fprintln(stdout, lines[i]); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, lines[i+n]); err != nil {
			return err
		}
	}
	return nil
}
