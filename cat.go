package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	catCmd = cobra.Command{
		Use: "cat [input-files...]",
		Long: `Concatenate pairs of files and output them ` +
			`in a way for the align command to consume`,
		Short: `Concatenate pairs of files`,
		RunE:  cat,
		//		Args:  cobra.ExactArgs(0),
	}
	ext = ".txt"
)

func init() {
	catCmd.Flags().StringVarP(&ext, "ext", "e",
		ext, "set extension for other input file")
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

func catFiles(f1, f2 string, out io.Writer) error {
	log.Printf("f1 = %s, f2 = %s", f1, f2)
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
	return catReaders(in1, in2, out)
}

func catReaders(in1, in2 io.Reader, out io.Writer) error {
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
		if _, err := fmt.Fprintf(out, "%s\n%s\n", t1, t2); err != nil {
			return err
		}
	}
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

func otherFileName(f1 string) string {
	ext1 := filepath.Ext(f1)
	return f1[0:len(f1)-len(ext1)] + ext
}
