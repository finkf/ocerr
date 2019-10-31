package main

import (
	"io"
	"os"

	"github.com/finkf/lev"
	"github.com/spf13/cobra"
)

var (
	splitCmd = cobra.Command{
		Use:   "split",
		Long:  `Splits blocks of alignments at a set of characters`,
		Short: `Split blocks into tokens`,
		RunE:  runSplit,
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

func runSplit(cmd *cobra.Command, args []string) error {
	return split(os.Stdin, os.Stdout)
}

func split(stdin io.Reader, stdout io.Writer) error {
	return readBlocks(stdin, func(b block) error {
		return splitBlocks(b, stdout)
	})
}

func splitBlocks(b block, stdout io.Writer) error {
	i := 0
	for j := indexAny(b.a.S1[i:], splitCharSet); j > 0; {
		if err := splitBlock(b, i, j).write(stdout); err != nil {
			return err
		}
		i, j = nextSplitBlock(b, splitCharSet, j)
	}
	return splitBlock(b, i, len(b.a.S1)).write(stdout)
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
		p1: b.p1,
		p2: b.p2,
		a: lev.Alignment{
			S1:    b.a.S1[i:j],
			S2:    b.a.S2[i:j],
			Trace: b.a.Trace[i:j],
		},
	}
}

func nextSplitBlock(b block, set string, j int) (int, int) {
	i := j + 1
	j = indexAny(b.a.S1[i:], splitCharSet)
	if j == -1 {
		return i, j
	}
	return i, j + i
}
