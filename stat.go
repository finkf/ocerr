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
	global = make(map[byte]uint)
	local  = make(map[byte]uint)
)

// Args: --max (max number of stats 0=all)
//       --error-stats print error stats
func init() {
	// catCmd.Flags().StringVarP(&catPattern1, "p1", "1",
	// 	catPattern1, "set regex pattern for first input file")
	// catCmd.Flags().StringVarP(&catReplacePattern, "p2", "2",
	// 	catReplacePattern, "set replacement pattern for second input file")
	// catCmd.Flags().BoolVarP(&catFileName, "file-names", "f",
	// 	false, "output first filename")
}

func stat(cmd *cobra.Command, args []string) error {
	err := readBlocks(os.Stdin, func(b block) error {
		return statBlock(b)
	})
	if err != nil {
		return err
	}
	_, err = fmt.Printf("global: %s\n", countsToString(global))
	return err
}

func statBlock(b block) error {
	// clear local stats
	for _, b := range []byte{lev.Del, lev.Sub, lev.Ins, lev.Nop} {
		local[b] = 0
	}
	for i := 0; i < len(b.a.Trace); i++ {
		op := b.a.Trace[i]
		local[op]++
		global[op]++
	}
	b.stats = countsToString((local))
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
