package main

import (
	"fmt"
	"os"

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
	statGlobal = make(map[byte]uint)
	statLocal  = make(map[byte]uint)
	statMax    = 0
)

// Args: --max (max number of stats 0=all)
//       --error-stats print error stats
func init() {
	statCmd.Flags().IntVarP(&statMax, "max", "m",
		0, "set maximal number of printed error patterns")
}

func stat(cmd *cobra.Command, args []string) error {
	err := readBlocks(os.Stdin, func(b block) error {
		return statBlock(b)
	})
	if err != nil {
		return err
	}
	_, err = fmt.Printf("global: %s\n", countsToString(statGlobal))
	return err
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
	b.stats = countsToString((statLocal))
	return printBlock(b, os.Stdout)
}

func countsToString(counts map[byte]uint) string {
	c := counts[lev.Nop]
	s := counts[lev.Sub]
	i := counts[lev.Ins]
	d := counts[lev.Del]
	res := float64(c) / float64(c+s+i+d)
	return fmt.Sprintf("Acc=c/(c+s+i+d)=%d/(%d+%d+%d+%d)=%f", c, c, s, i, d, res)
}
