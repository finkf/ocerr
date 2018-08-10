package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/finkf/lev"
	"github.com/spf13/cobra"
)

var (
	statCmd = cobra.Command{
		Use:   "stat",
		Long:  `Calculate error statistics of alignment blocks`,
		Short: `Calculate error statistics`,
		RunE:  stat,
		Args:  cobra.ExactArgs(0),
	}
	statGlobal        = make(map[byte]uint)
	statLocal         = make(map[byte]uint)
	statErrorPatterns = make(map[spair]uint)
	statMax           = 0
)

type spair struct {
	first, second string
}

func init() {
	statCmd.Flags().IntVarP(&statMax, "max", "m",
		0, "set maximal number of printed error patterns (0=all, -1=none)")
}

func stat(cmd *cobra.Command, args []string) error {
	err := readBlocks(os.Stdin, func(b block) error {
		return statBlock(b)
	})
	if err != nil {
		return err
	}
	if _, err := fmt.Printf("Global: %s\n", countsToString(statGlobal)); err != nil {
		return err
	}
	return printErrorPatterns()
}

func statBlock(b block) error {
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
	b.stats = countsToString((statLocal))
	return printBlock(b, os.Stdout)
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

func printErrorPatterns() error {
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
	if _, err := fmt.Printf("Errors: %d\n", total); err != nil {
		return err
	}
	for i := range ep {
		if statMax > 0 && i >= statMax {
			break
		}
		if _, err := fmt.Printf("%d:\t{%s}\t{%s}\n",
			ep[i].c, ep[i].p.first, ep[i].p.second); err != nil {
			return err
		}
	}
	return nil
}

func countsToString(counts map[byte]uint) string {
	c := counts[lev.Nop]
	s := counts[lev.Sub]
	i := counts[lev.Ins]
	d := counts[lev.Del]
	res := float64(c) / float64(c+s+i+d)
	return fmt.Sprintf("Acc=c/(c+s+i+d)=%d/(%d+%d+%d+%d)=%f",
		c, c, s, i, d, res)
}
