package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	countCmd = cobra.Command{
		Use:   "count",
		Short: `Counts valid and invalid lines`,
		RunE:  runCount,
	}
)

func runCount(cmd *cobra.Command, args []string) error {
	return count(os.Stdin, os.Stdout)
}

func count(stdin io.Reader, stdout io.Writer) error {
	var total, errors int
	s := newBlockScanner(stdin)
	for s.scan() {
		total++
		if bytes.ContainsAny(s.block().a.Trace, "-#+") {
			errors++
		}
	}
	if s.err() != nil {
		return s.err()
	}
	acc := float64(total-errors) / float64(total)
	ers := float64(errors) / float64(total)
	if _, err := fmt.Fprintf(stdout, "total: %d, errors: %d, accuracy: %f, errors: %f\n",
		total, errors, acc, ers); err != nil {
		return err
	}
	return nil
}
