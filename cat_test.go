package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

var (
	testUpdateGoldFile = false
)

func init() {
	flag.BoolVar(&testUpdateGoldFile, "update", false, "update gold files")
}

type subCmdFunc func(*cobra.Command, []string) error

func withInput(t *testing.T, fn string, f subCmdFunc) subCmdFunc {
	t.Helper()
	in, err := os.Open(filepath.Join("testdata", fn))
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	return func(cmd *cobra.Command, args []string) error {
		os.Stdin = in
		defer in.Close()
		return f(cmd, args)
	}
}

func runSubCmd(t *testing.T, f subCmdFunc, args []string) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	defer func() { _ = r.Close() }()
	os.Stdout = w
	if err = f(nil, args); err != nil {
		t.Fatalf("got error: %v", err)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("got error: %v", err)
	}
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	os.Stdout = oldStdout
	return string(bs)
}

func checkGoldFile(t *testing.T, gold, got string) {
	// update the gold file with the given output
	if testUpdateGoldFile {
		outfile := filepath.Join("testdata", gold)
		if err := ioutil.WriteFile(outfile, []byte(got), os.ModePerm); err != nil {
			t.Fatalf("got error: %v", err)
		}
	}
	t.Helper()
	in, err := os.Open(filepath.Join("testdata", gold))
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	defer func() { _ = in.Close() }()
	want, err := ioutil.ReadAll(in)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if string(want) != got {
		t.Fatalf("expected %q; got %q in %s", want, got, gold)
	}
}

func TestSubCmds(t *testing.T) {
	tests := []struct {
		gold string
		f    subCmdFunc
		args []string
	}{
		{"cat_gold.txt", cat, []string{"testdata/0001.gt.txt", "testdata/0002.gt.txt"}},
		{"align_gold.txt", withInput(t, "cat_gold.txt", align), nil},
	}
	for _, tc := range tests {
		t.Run(tc.gold, func(t *testing.T) {
			got := runSubCmd(t, tc.f, tc.args)
			checkGoldFile(t, tc.gold, got)
		})
	}
}
