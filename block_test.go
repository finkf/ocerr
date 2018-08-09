package main

import (
	"bytes"
	"testing"

	"github.com/finkf/lev"
)

func TestPrintBlock(t *testing.T) {
	tests := []struct {
		fn, s1, s2, want string
	}{
		{"", "abc", "def", "abc\n###\ndef\n" + endOfBlock + "\n"},
		{"a.txt", "x", "x", "a.txt\nx\n|\nx\n" + endOfBlock + "\n"},
	}
	for _, tc := range tests {
		t.Run(tc.s1+" "+tc.s2, func(t *testing.T) {
			var l lev.Lev
			a, err := l.Alignment(l.Trace(tc.s1, tc.s2))
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			b := block{a: a, fn: tc.fn}
			var buf bytes.Buffer
			if err := printBlock(b, &buf); err != nil {
				t.Fatalf("got error: %v", err)
			}
			if got := buf.String(); got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}
