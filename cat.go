package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	catCmd = cobra.Command{
		Use:   "cat FILE [FILES...]",
		Short: `Cat lines from files`,
		RunE:  runCat,
		Args:  cobra.MinimumNArgs(1),
		Long: `Concatenates one-line files with their sibling files.
The sibling of a file is a file with the same name but
with another extension.`,
	}
	cat struct {
		ext string
	}
)

func init() {
	catCmd.Flags().StringVarP(&cat.ext, "other", "o",
		".gt.txt", "other file extension")
}

func runCat(_ *cobra.Command, args []string) error {
	for _, file := range args {
		if err := catFile(file); err != nil {
			return err
		}
	}
	return nil
}

func catFile(file string) error {
	// this file
	this, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot cat: %v", err)
	}
	if len(this) > 0 && this[len(this)-1] == '\n' {
		this = this[0 : len(this)-1]
	}
	// other file
	otherFile := stem(file) + cat.ext
	other, err := ioutil.ReadFile(otherFile)
	if err != nil {
		return fmt.Errorf("cannot cat: %v", err)
	}
	if len(other) > 0 && other[len(other)-1] == '\n' {
		other = other[0 : len(other)-1]
	}
	catWrite(this, other, file, otherFile)
	return nil
}

func catWrite(this, other []byte, file, otherFile string) {
	if global.separator == "" {
		fmt.Printf("%s\n%s\n", this, other)
		return
	}
	// n := max(len(file), len(otherFile))
	fmt.Printf("%s%s%s\n%s%s%s\n",
		file, global.separator, this,
		otherFile, global.separator, other)
}

func stem(path string) string {
	base := filepath.Base(path)
	pos := strings.Index(base, ".")
	if pos <= 0 {
		return path
	}
	return filepath.Join(filepath.Dir(path), base[0:pos])
}
