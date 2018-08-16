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
	rootCmd.AddCommand(&catCmd)
	rootCmd.AddCommand(&matchCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// no need to print error message
		// since cobra takes care of this
		os.Exit(1)
	}
}
