package main

import (
	"os"

	"github.com/finkf/lev"
	"github.com/spf13/cobra"
)

var (
	splitCmd = cobra.Command{
		Use:   "split",
		Long:  `Splits blocks of alignments at a set of characters`,
		Short: `Split blocks into tokens`,
		RunE:  split,
		Args:  cobra.ExactArgs(0),
	}
	splitCharSet string
)

const (
	defaultSplitCharSet = "\t "
)

func init() {
	splitCmd.Flags().StringVarP(&splitCharSet, "chars", "c",
		defaultSplitCharSet, "set the character set used to split blocks")
}

func split(cmd *cobra.Command, args []string) error {
	return readBlocks(os.Stdin, func(b block) error {
		return splitBlocks(b)
	})
}

func splitBlocks(b block) error {
	i := 0
	for j := indexAny(b.a.S1[i:], splitCharSet); j > 0; {
		printBlock(splitBlock(b, i, j), os.Stdout)
		i, j = nextSplitBlock(b, splitCharSet, i, j)
	}
	printBlock(splitBlock(b, i, len(b.a.S1)), os.Stdout)
	return nil
}

func indexAny(rs []rune, set string) int {
	for i, r := range rs {
		for _, c := range set {
			if r == c {
				return i
			}
		}
	}
	return -1
}

func splitBlock(b block, i, j int) block {
	return block{
		fn: b.fn,
		a: lev.Alignment{
			S1:    b.a.S1[i:j],
			S2:    b.a.S2[i:j],
			Trace: b.a.Trace[i:j],
		},
	}
}

func nextSplitBlock(b block, set string, i, j int) (int, int) {
	i = j + 1
	j = indexAny(b.a.S1[i:], splitCharSet)
	if j == -1 {
		return i, j
	}
	return i, j + i
}
