package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = cobra.Command{
		Use:   "ocerr",
		Long:  `tools for ocr error examination`,
		Short: `tools for ocr error examination`,
	}
	global struct {
		endOfBlock string
		separator  string
	}
)

func init() {
	rootCmd.AddCommand(&alignCmd)
	rootCmd.AddCommand(&splitCmd)
	rootCmd.AddCommand(&statCmd)
	rootCmd.AddCommand(&matchCmd)
	rootCmd.AddCommand(&pairCmd)
	rootCmd.AddCommand(&countCmd)
	rootCmd.AddCommand(&catCmd)
	rootCmd.PersistentFlags().StringVarP(
		&global.endOfBlock, "eob", "b", "%%", "Set the end of block marker")
	rootCmd.PersistentFlags().StringVarP(
		&global.separator, "separator", "F", "",
		"Set the separator for prefixes (empty string means no prefix)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// no need to print error message
		// since cobra takes care of this
		os.Exit(1)
	}
}
