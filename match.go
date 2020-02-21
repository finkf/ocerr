package main

import (
	"fmt"
	"io"
	"os"

	"github.com/finkf/lev"
	"github.com/spf13/cobra"
)

var (
	matchCmd = cobra.Command{
		Use:   "match PATTERN [PATTERN...]",
		Short: `Print blocks that match a pattern`,
		RunE:  runMatch,
		Args:  cobra.MinimumNArgs(1),
		Long: `Printblocks that match PATTERN or PATTERNS...
A pattern matches strings on the upper line, the
alignment line and the lower line.  To match any
character use '.' (it is impossible to match a
literal dot). So for examplethe pattern 'ab##..'
matches lines where 'ab' has a substitution with
any two characters.`,
	}
	grepInverted bool
)

func init() {
	matchCmd.Flags().BoolVarP(&grepInverted, "inverted", "v",
		false, "invert matching")
}

func runMatch(cmd *cobra.Command, args []string) error {
	return match(os.Stdin, os.Stdout, args...)
}

func match(stdin io.Reader, stdout io.Writer, args ...string) error {
	ms, err := newMatchers(args)
	if err != nil {
		return err
	}
	return readBlocks(stdin, func(b block) error {
		if ms.match(b.a) != grepInverted {
			return b.write(stdout)
		}
		return nil
	})
}

const (
	matchAll = 0
)

type matchers []matcher

func newMatchers(ps []string) (matchers, error) {
	var ms matchers
	for _, p := range ps {
		m, err := newMatcher(p)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

func (ms matchers) match(a lev.Alignment) bool {
	for _, m := range ms {
		for i := 0; i < len(a.S1); i++ {
			if m.matchAt(a, i) {
				return true
			}
		}
	}
	return false
}

type matcher struct {
	S1, S2 []rune
	Trace  []byte
}

func newMatcher(p string) (matcher, error) {
	rs := unescapePattern(p)
	n := len(rs)
	if n%3 != 0 {
		return matcher{}, fmt.Errorf("invalid pattern: %q", p)
	}
	m := matcher{
		S1: rs[0 : n/3],
		S2: rs[2*n/3:],
	}
	for i := n / 3; i < 2*n/3; i++ {
		switch byte(rs[i]) {
		case lev.Ins, lev.Del, lev.Sub, lev.Nop, matchAll:
			m.Trace = append(m.Trace, byte(rs[i]))
		default:
			return matcher{}, fmt.Errorf("invalid pattern: %q", p)
		}
	}
	return m, nil
}

func (m matcher) matchAt(a lev.Alignment, i int) bool {
	if i+len(m.S1) >= len(a.S1) {
		return false
	}
	return matchRunes(m.S1, a.S1[i:]) &&
		matchRunes(m.S2, a.S2[i:]) &&
		matchBytes(m.Trace, a.Trace[i:])
}

func matchRunes(a, b []rune) bool {
	for i := range a {
		if a[i] == matchAll {
			continue
		}
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func matchBytes(a, b []byte) bool {
	for i := range a {
		if a[i] == matchAll {
			continue
		}
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func unescapePattern(p string) []rune {
	var cnv []rune
	var escaped bool
	for _, r := range p {
		if !escaped && r == '\\' {
			escaped = true
			continue
		}
		if escaped {
			// use litteral character
			cnv = append(cnv, r)
			escaped = false
			continue
		}
		escaped = false
		if r == '.' {
			cnv = append(cnv, matchAll)
			continue
		}
		cnv = append(cnv, r)
	}
	return cnv
}
