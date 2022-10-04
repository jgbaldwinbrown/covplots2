package covplots

import (
	"github.com/jgbaldwinbrown/shellout/pkg"
	"strings"
	"io"
	"bufio"
	"os"
	"fmt"
	"flag"
)

func GetAllMultiplotFlags() AllSingleFlags {
	var f AllSingleFlags
	flag.StringVar(&f.Config, "i", "", "Input config file. JSON, following the documented format.")
	flag.IntVar(&f.WinSize, "w", 1000000, "Sliding window plot size (default = 1000000).")
	flag.IntVar(&f.WinStep, "s", 1000000, "Sliding window step distance (default = 1000000).")
	flag.IntVar(&f.Threads, "t", 8, "Threads to run simultaneously")
	flag.Parse()

	return f
}

func RunAllMultiplot() {
	f := GetAllMultiplotFlags()
	cfg, err := GetUltimateConfig(f.Config)

	err = AllMultiplotParallel(cfg, f.WinSize, f.WinStep, f.Threads)
	if err != nil {
		panic(err)
	}
}

func FilterMulti(chr string, start, end int, rs ...io.Reader) ([]*Filterer, error) {
	var out []*Filterer
	for _, r := range rs {
		fr, err := Filter(r, chr, start, end)
		if err != nil {
			return nil, fmt.Errorf("FilterMulti: %w", err)
		}
		out = append(out, fr)
	}
	return out, nil
}

func CombineSinglebpPlots(names []string, rs ...io.Reader) (*strings.Reader, error) {
	var out strings.Builder
	for i, r := range rs {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		for s.Scan() {
			fmt.Fprintf(&out, "%s\t%s\n", s.Text(), names[i])
		}
	}
	return strings.NewReader(out.String()), nil
}

func PlotMulti(outpre string, toplot io.Reader) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov %v %v
`,
		fmt.Sprintf("%v_multi_plfmt.bed", outpre),
		fmt.Sprintf("%v_multi_plotted.png", outpre),
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func Nop(string, []io.Reader, []string, any) error {return nil}

func GetFunc(fstr string) func(string, []io.Reader, []string, any) error {
	switch fstr {
	case "subtract_some": return SubtractSome
	default: return Nop
	}
	return Nop
}

func OpenPaths(paths ...string) ([]io.Reader, error) {
	var out []io.Reader
	for _, path := range paths {
		r, err := os.Open(path)
		if err != nil {
			CloseFiles(out...)
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CloseFiles(files ...io.Reader) {
	for _, f := range files {
		if c, ok := f.(io.Closer); ok {
			c.Close()
		}
	}
}

func Multiplot(

func MultiplotInputSet(outpre string, cfg InputSet, chr string, start, end int) error {
	f := GetFunc(cfg.Function)
	rs, err := OpenPaths(cfg.Paths...)
	if err != nil {
		return err
	}
	defer CloseFiles(rs...)

	frs, err := FilterMulti(chr, start, end, rs...)
	if err != nil {
		return err
	}

	return f(cfg.Outpre, frs, cfg.Names, cfg.FunctionArgs)
}

func ParseSubArgs(args any) ([][]int, error) {
	var out [][]int
	anysl, ok := args.([]any)
	if !ok {
		return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
	}
	for _, anypair := range anysl {
		pair, ok := anypair.([]any)
		if !ok {
			return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
		}
		if len(pair) != 2 {
			return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
		}
		entry := make([]int, 2)
		for i, val := range pair {
			ival, ok := val.(int)
			if !ok {
				return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
			}
			entry[i] = ival
		}
		out = append(out, entry)
	}
	return out, nil
}

func SubtractTwo(outpre string, rs []io.Reader, names []string, args any) error {
	sargs, err := ParseSubArgs(args)
	if err != nil {
		return fmt.Errorf("SubtractSome: %w", err)
	}
	var to_combine []io.Reader
	var newnames []string
	used := make(map[int]struct{})
	for _, entry := range sargs {
		newreader, err := Subtract(rs[entry[0]], rs[entry[1]])
		if err != nil {
			return fmt.Errorf("SubtractSome: %w", err)
		}

		newnames = append(newnames, fmt.Sprintf("%v-%v", names[entry[0]], names[entry[1]]))

		to_combine = append(to_combine, newreader)
		used[entry[0]] = struct{}{}
		used[entry[1]] = struct{}{}
	}
	for i, r := range rs {
		if _, ok := used[i]; !ok {
			to_combine = append(to_combine, r)
			newnames = append(newnames, names[i])
		}
	}

	combined, err := CombineSinglePlots(newnames, to_combine...)
	if err != nil {
		return fmt.Errorf("SubtractSome: %w", err)
	}

	err := PlotMulti(outpre, combined)
	if err != nil {
		return fmt.Errorf("SubtractSome: %w", err)
	}

	return nil
}

func MultiplotSlide(cfg UltimateConfig, chr string, winsize, winstep int) error {
	chrlens, err := GetChrLens(cfg.Chrlens)
	if err != nil {
		return fmt.Errorf("MultiplotSlide: %w", err)
	}

	for _, chrlenset := range chrlens {
		chr, chrlen := chrlenset.Chr, chrlenset.Len
		for start := 0; start < chrlen; start += winstep {
			end := start + winsize
			err := Multiplot(cfg, chr, start, end)
			if err != nil {
				return fmt.Errorf("SubtractSinglePlotWins: %w", err)
			}
		}
	}

	return nil
}

func AllMultiplotParallel(cfgs UltimateConfig, winsize, winstep, threads int) error {
	jobs := make(chan Config, len(cfgs))
	for _, cfg := range cfgs {
		jobs <- cfg
	}
	close(jobs)

	errs := make(chan error, len(cfgs))

	for i:=0; i<threads; i++ {
		go func() {
			for cfg := range jobs {
				errs <- MultiplotSlide(cfg, winsize, winstep)
			}
		}()
	}

	var out Errors
	for i:=0; i<len(cfgs); i++ {
		err := <-errs
		if err != nil {
			out = append(out, err)
		}
	}
	if len(out) < 0 {
		return out
	}
	return nil
}


