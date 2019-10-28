package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

var (
	updateGoldFiles = false
)

func init() {
	flag.BoolVar(&updateGoldFiles, "update", false, "update gold files")
}

func withCatFiles(a, b string, f func(io.Reader)) {
	abs, err := ioutil.ReadFile(a)
	ensureReadTestFile(a, err)
	bbs, err := ioutil.ReadFile(b)
	ensureReadTestFile(b, err)
	var buf bytes.Buffer
	if _, err := buf.Write(abs); err != nil {
		panic(err)
	}
	if _, err := buf.Write(bbs); err != nil {
		panic(err)
	}
	f(&buf)
}

func withFile(a string, f func(io.Reader)) {
	abs, err := ioutil.ReadFile(a)
	ensureReadTestFile(a, err)
	buf := bytes.NewBuffer(abs)
	f(buf)
}

func ensureReadTestFile(path string, err error) {
	if err != nil {
		panic(fmt.Sprintf("cannot read testfile: %s: %v", path, err))
	}
}

func gold(t *testing.T, in io.Reader, gold string) {
	t.Helper()
	got, err := ioutil.ReadAll(in)
	ensureReadTestFile("input", err)
	if updateGoldFiles {
		if err = ioutil.WriteFile(gold, got, 0666); err != nil {
			panic(fmt.Sprintf("cannot write %s: %v", gold, err))
		}
	}
	want, err := ioutil.ReadFile(gold)
	ensureReadTestFile(gold, err)
	if !bytes.Equal(want, got) {
		t.Fatalf("input not equal to %s", gold)
	}
}

func TestPair(t *testing.T) {
	withCatFiles("testdata/a.txt", "testdata/b.txt", func(in io.Reader) {
		var b bytes.Buffer
		if err := pair(in, &b); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/pair.gold.txt")
	})
}

func TestAlign(t *testing.T) {
	withFile("testdata/pair.gold.txt", func(in io.Reader) {
		var b bytes.Buffer
		if err := align(in, &b); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/align.gold.txt")
	})
}

func TestAlignPrefix(t *testing.T) {
	withFile("testdata/pair-prefix.txt", func(in io.Reader) {
		global.separator = " "
		defer func() { // reset
			global.separator = ""
		}()
		var b bytes.Buffer
		if err := align(in, &b); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/align-prefix.gold.txt")
	})
}

func TestStat(t *testing.T) {
	withFile("testdata/align.gold.txt", func(in io.Reader) {
		var b bytes.Buffer
		if err := stat(in, &b); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/stat.gold.txt")
	})
}

func TestSplit(t *testing.T) {
	withFile("testdata/align.gold.txt", func(in io.Reader) {
		var b bytes.Buffer
		if err := split(in, &b); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/split.gold.txt")
	})
}

func TestMatch(t *testing.T) {
	withFile("testdata/align.gold.txt", func(in io.Reader) {
		var b bytes.Buffer
		if err := match(in, &b, "ine|||ine", ".#("); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/match.gold.txt")
	})
}

func TestCount(t *testing.T) {
	withFile("testdata/split.gold.txt", func(in io.Reader) {
		var b bytes.Buffer
		if err := count(in, &b); err != nil {
			t.Fatalf("got error: %v", err)
		}
		gold(t, &b, "testdata/count.gold.txt")
	})
}
