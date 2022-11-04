package covplots

import (
	"math"
	"sort"
	"github.com/jgbaldwinbrown/shellout/pkg"
	"io"
	"os"
	"fmt"
)

func GetPlotFunc(fstr string) func(outpre string, ylim []float64, args any) error {
	switch fstr {
	case "plot_multi": return PlotMultiAny
	case "": return PlotMultiAny
	case "plot_cov_vs_pair": return PlotCovVsPair
	default: return PlotPanic
	}
	return PlotPanic
}
func PlotPanic(outpre string, ylim []float64, args any) error {
	panic(fmt.Errorf("trying to use an unimplemented plot function"))
	return nil
}

func PlotMultiAny(outpre string, ylim []float64, args any) error {
	return PlotMulti(outpre, ylim)
}

func PlotCovVsPair(outpre string, ylim []float64, args any) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_cov_vs_pair %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PosLess(p1, p2 Pos) bool {
	if p1.Chr < p2.Chr {
		return true
	}
	if p1.Chr > p2.Chr {
		return false
	}
	return p1.Bp < p2.Bp
}

func SortAndUniqPoses[T any](maps ...map[Pos]T) []Pos {
	uniques := map[Pos]struct{}{}
	for _, posmap := range maps {
		for pos, _ := range posmap {
			uniques[pos] = struct{}{}
		}
	}
	poses := make([]Pos, 0, len(uniques))
	for pos, _ := range uniques {
		poses = append(poses, pos)
	}
	sort.Slice(poses, func(i, j int) bool{return PosLess(poses[i], poses[j])})
	return poses
}

func CombineToOneLine(rs []io.Reader, args any) ([]io.Reader, error) {
	maps := []map[Pos]float64{}
	for _, r := range rs {
		posmap, err := CollectVals(r)
		if err != nil {
			return nil, fmt.Errorf("CombineToOneLine: error during CollectVals %w", err)
		}
		maps = append(maps, posmap)
	}

	fmt.Fprintf(os.Stderr, "%v\n", maps[1])

	poses := SortAndUniqPoses(maps...)
	// fmt.Fprintf(os.Stderr, "%v\n", poses)

	out := PipeWrite(func(w io.Writer) {
		vals := []float64{}
		for _, pos := range poses {
			vals = vals[:0]
			for _, posmap := range maps {
				if val, ok := posmap[pos]; ok {
					vals = append(vals, val)
				} else {
					vals = append(vals, math.NaN())
				}
			}
			fmt.Fprintf(w, "%v\t%v\t%v", pos.Chr, pos.Bp, pos.Bp+1)
			for _, val := range vals {
				fmt.Fprintf(w, "\t%v", val)
			}
			fmt.Fprintf(w, "\n");
		}
	})
	return []io.Reader{out}, nil
}
