package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	catCmd = cobra.Command{
		Use: "cat [input-files...]",
		Long: `Concatenate pairs of files and output them ` +
			`in a way for the align command to consume`,
		Short: `Concatenate pairs of files`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  cat,
	}
	catPattern1       = `\.gt\.txt$`
	catReplacePattern = ".txt"
)

func init() {
	catCmd.Flags().StringVarP(&catPattern1, "p1", "1",
		catPattern1, "set regex pattern for first input file")
	catCmd.Flags().StringVarP(&catReplacePattern, "p2", "2",
		catReplacePattern, "set replacement pattern for second input file")
	catCmd.Flags().BoolVarP(&gocrFileName, "file-names", "f",
		false, "output first filename")
}

func cat(cmd *cobra.Command, args []string) error {
	for _, f1 := range args {
		f2 := otherFileName(f1)
		if err := catFiles(f1, f2, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

/* #nosec */
func catFiles(f1, f2 string, out io.Writer) error {
	in1, err := os.Open(f1)
	if err != nil {
		return err
	}
	defer func() { _ = in1.Close() }()

	in2, err := os.Open(f2)
	if err != nil {
		return err
	}
	defer func() { _ = in2.Close() }()
	return catReaders(f1, in1, in2, out)
}

func catReaders(fn string, in1, in2 io.Reader, out io.Writer) error {
	s1 := bufio.NewScanner(in1)
	s2 := bufio.NewScanner(in2)
	for {
		t1, ok, err := nextNonEmptyLine(s1)
		if !ok {
			return err
		}
		t2, ok, err := nextNonEmptyLine(s2)
		if !ok {
			return err
		}
		if err := print2lines(fn, t1, t2, out); err != nil {
			return err
		}
	}
}

func print2lines(fn, t1, t2 string, out io.Writer) error {

	if gocrFileName {
		if _, err := fmt.Println(fn); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(out, "%s\n%s\n", t1, t2)
	return err
}

func nextNonEmptyLine(s *bufio.Scanner) (string, bool, error) {
	for s.Scan() {
		t := s.Text()
		if len(t) > 0 {
			return t, true, nil
		}
	}
	return "", false, s.Err()
}

var catRegexPattern *regexp.Regexp

func otherFileName(f1 string) string {
	if catRegexPattern == nil {
		catRegexPattern = regexp.MustCompile(catPattern1)
	}
	return catRegexPattern.ReplaceAllString(f1, catReplacePattern)
}
