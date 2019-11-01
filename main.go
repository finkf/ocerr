package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = cobra.Command{
		Use:   "ocerr",
		Short: `Tool for ocr error examination`,
		Long: `Tool for ocr error examination

This tool can be used to examine pairs of ground-truth test lines.
The base are always pairs of input lines (first line is the
ground-truth, second line the test).  These lines can be aligned (see
the align command) for further processing.

Most of the command operate on blocks.  A block is a sequence of 3
lines separated by a end of block marker (EOB) on a single line.  Each
line of a block can have a prefix seperated by a separator.  If the
separator is the empty string, no prefix is assumed.

You can either use the --eob|-b, --sep|-F command line options to
change them or set the environment variales OCERREOB or OCERRSEP
accordingly.
`,
	}
	global struct {
		endOfBlock string
		separator  string
	}
)

func init() {
	defEOB := "%%"
	if val, set := os.LookupEnv("OCERREOB"); set {
		defEOB = val
	}
	defSep := ""
	if val, set := os.LookupEnv("OCERRSEP"); set {
		defSep = val
	}
	rootCmd.AddCommand(&alignCmd)
	rootCmd.AddCommand(&splitCmd)
	rootCmd.AddCommand(&statCmd)
	rootCmd.AddCommand(&matchCmd)
	rootCmd.AddCommand(&pairCmd)
	rootCmd.AddCommand(&countCmd)
	rootCmd.PersistentFlags().StringVarP(
		&global.endOfBlock, "eob", "b", defEOB,
		"Set the end of block marker")
	rootCmd.PersistentFlags().StringVarP(
		&global.separator, "separator", "F", defSep,
		"Set the sep for prefixes")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// no need to print error message
		// since cobra takes care of this
		os.Exit(1)
	}
}
