package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = cobra.Command{
		Use:   "gocr",
		Long:  `tools for ocr error examination`,
		Short: `tools for ocr error examination`,
	}
)

func init() {
	rootCmd.AddCommand(&alignCmd)
	rootCmd.AddCommand(&splitCmd)
	rootCmd.AddCommand(&statCmd)
	rootCmd.AddCommand(&matchCmd)
	rootCmd.AddCommand(&pairCmd)
	rootCmd.AddCommand(&countCmd)
	rootCmd.PersistentFlags().StringVarP(
		&endOfBlock, "eob", "b", endOfBlock, "Set the end of block marker")
	rootCmd.PersistentFlags().StringVarP(
		&separator, "separator", "F", separator, "Set the separator for prefixes")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// no need to print error message
		// since cobra takes care of this
		os.Exit(1)
	}
}
