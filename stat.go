package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/finkf/lev"
	"github.com/spf13/cobra"
)

type spair struct {
	first, second string
}

type counts map[byte]uint

var (
	statCmd = cobra.Command{
		Use:   "stat",
		Short: `Display error statistics`,
		RunE:  runStat,
		Args:  cobra.ExactArgs(0),
		Long: `Calculate and display error statistics of
alignment blocks`,
	}
	statGlobal        = make(counts)
	statLocal         = make(counts)
	statErrorPatterns = make(map[spair]uint)
	statMax           = 0
	statCut           = 0
)

func init() {
	statCmd.Flags().IntVarP(&statMax, "max", "m",
		0, "set maximal number of printed error patterns (0=all, -1=none)")
	statCmd.Flags().IntVarP(&statCut, "cut", "c",
		0, "do not print error patterns with a count smaller than cut")
}

func runStat(cmd *cobra.Command, args []string) error {
	return stat(os.Stdin, os.Stdout)
}

func stat(stdin io.Reader, stdout io.Writer) error {
	err := readBlocks(stdin, func(b block) error {
		return statBlock(b, stdout)
	})
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Global: %s\n", statGlobal); err != nil {
		return err
	}
	return printErrorPatterns(stdout)
}

func statBlock(b block, stdout io.Writer) error {
	// clear local stats
	for _, b := range []byte{lev.Del, lev.Sub, lev.Ins, lev.Nop} {
		statLocal[b] = 0
	}
	for i := 0; i < len(b.a.Trace); i++ {
		op := b.a.Trace[i]
		statLocal[op]++
		statGlobal[op]++
	}
	addErrorPatterns(b)
	b.stats = statLocal.String()
	return b.write(stdout)
}

func addErrorPatterns(b block) {
	if statMax < 0 {
		return
	}
	for i := 0; i < len(b.a.Trace); {
		op := b.a.Trace[i]
		if op == lev.Nop {
			i++
			continue
		}
		j := i + 1
		for ; j < len(b.a.Trace); j++ {
			if b.a.Trace[j] != op {
				break
			}
		}
		p := spair{string(b.a.S1[i:j]), string(b.a.S2[i:j])}
		statErrorPatterns[p]++
		i = j
	}
}

func printErrorPatterns(stdout io.Writer) error {
	if statMax < 0 {
		return nil
	}
	var ep []struct {
		p spair
		c uint
	}
	var total uint
	for p, c := range statErrorPatterns {
		total += c
		ep = append(ep, struct {
			p spair
			c uint
		}{p, c})
	}
	// sort descending by counts
	// sort ascending by first and second string
	sort.Slice(ep, func(i, j int) bool {
		if ep[i].c == ep[j].c {
			if ep[i].p.first == ep[j].p.first {
				return ep[i].p.second < ep[j].p.second
			}
			return ep[i].p.first < ep[j].p.first
		}
		return ep[i].c > ep[j].c
	})
	if _, err := fmt.Fprintf(stdout, "Errors: %d\n", total); err != nil {
		return err
	}
	var n int
	for i := range ep {
		if statMax > 0 && n >= statMax {
			break
		}
		if ep[i].c <= uint(statCut) {
			continue
		}
		n++
		if _, err := fmt.Fprintf(stdout, "%d:\t{%s}\t{%s}\n",
			ep[i].c, ep[i].p.first, ep[i].p.second); err != nil {
			return err
		}
	}
	return nil
}

func (c counts) String() string {
	x := c[lev.Nop]
	s := c[lev.Sub]
	i := c[lev.Ins]
	d := c[lev.Del]
	res := float64(x) / float64(x+s+i+d)
	return fmt.Sprintf("Acc=c/(c+s+i+d)=%d/(%d+%d+%d+%d)=%f",
		x, x, s, i, d, res)
}
