package covplots

import (
	"math"
	"github.com/montanaflynn/stats"
	"strconv"
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
	fmt.Println("one")
	f := GetAllMultiplotFlags()
	fmt.Println(f)
	cfg, err := GetUltimateConfig(f.Config)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	err = AllMultiplotParallel(cfg, f.WinSize, f.WinStep, f.Threads)
	if err != nil {
		panic(err)
	}
}

func FilterMulti(chr string, start, end int, rs ...io.Reader) ([]io.Reader, error) {
	var out []io.Reader
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
	fmt.Printf("len(rs): %v; names: %v\n", len(rs), names)
	var out strings.Builder
	for i, r := range rs {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		nlines := 0
		for s.Scan() {
			fmt.Fprintf(&out, "%s\t%s\n", s.Text(), names[i])
			nlines++
		}
		fmt.Printf("rs[%v] nlines: %v\n", i, nlines)
	}
	return strings.NewReader(out.String()), nil
}

func PlotMulti(outpre string, ylim []float64) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func Nop([]io.Reader, any) (io.Reader, error) {return nil, nil}

func GetFunc(fstr string) func(rs []io.Reader, args any) (io.Reader, error) {
	switch fstr {
	case "subtract_two": return SubtractTwo
	case "unchanged": return Unchanged
	case "normalize": return Normalize
	default: return Nop
	}
	return Nop
}

func OpenPaths(paths ...string) ([]io.Reader, error) {
	var out []io.Reader
	for _, path := range paths {
		r, err := os.Open(path)
		if err != nil {
			CloseReaders(out...)
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CloseReaders(files ...io.Reader) {
	for _, r := range files {
		if c, ok := r.(io.Closer); ok {
			c.Close()
		}
	}
}

func CloseFiles(files ...io.Closer) {
	for _, c := range files {
		c.Close()
	}
}

func MultiplotInputSet(cfg InputSet, chr string, start, end int) (io.Reader, []io.Closer, error) {
	f := GetFunc(cfg.Function)
	rs, err := OpenPaths(cfg.Paths...)
	if err != nil {
		return nil, nil, err
	}
	var closers []io.Closer
	for _, r := range rs {
		closers = append(closers, r.(io.Closer))
	}

	frs, err := FilterMulti(chr, start, end, rs...)
	if err != nil {
		return nil, nil, err
	}

	readers, err := f(frs, cfg.FunctionArgs)
	return readers, closers, err
}

func Multiplot(cfg UltimateConfig, chr string, start, end int) error {
	outpre := fmt.Sprintf("%s_%v_%v_%v", cfg.Outpre, chr, start, end)
	var rs []io.Reader
	for _, set := range cfg.InputSets {
		r, closers, err := MultiplotInputSet(set, chr, start, end)
		if err != nil {
			return err
		}
		defer CloseFiles(closers...)
		rs = append(rs, r)
	}

	var names []string
	for _, set := range cfg.InputSets {
		names = append(names, set.Name)
	}

	combined, err := CombineSinglebpPlots(names, rs...)
	if err != nil {
		return err
	}

	err = PlfmtSmall(combined, outpre)
	if err != nil {
		return err
	}

	ylim := []float64{-300,300}
	if cfg.Ylim != nil {
		ylim = cfg.Ylim
	}
	err = PlotMulti(outpre, ylim)
	if err != nil {
		return err
	}
	return nil
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

func SubtractTwo(rs []io.Reader, args any) (io.Reader, error) {
	newreader, err := Subtract(rs[0], rs[1])
	if err != nil {
		return nil, fmt.Errorf("SubtractSome: %w", err)
	}
	return newreader, nil
}

func Unchanged(rs []io.Reader, args any) (io.Reader, error) {
	if len(rs) != 1 {
		return nil, fmt.Errorf("Unchanged: wrong number of paths (%v)", len(rs))
	}
	return rs[0], nil
}

func NormalizeFloats(in []float64) []float64 {
	var nanfree []float64
	for _, f := range in {
		if !math.IsNaN(f) {
			nanfree = append(nanfree, f)
		}
	}
	m, err := stats.Mean(nanfree)
	if err != nil {
		m = 0
	}
	s, err := stats.StdDevP(nanfree)
	if err != nil {
		s = 1
	}
	out := make([]float64, len(in))
	for i, f := range in {
		out[i] = (f-m) / s
	}
	return out
}

func Normalize(rs []io.Reader, args any) (io.Reader, error) {
	if len(rs) != 1 {
		return nil, fmt.Errorf("Normalize: wrong number of paths (%v)", len(rs))
	}
	s := bufio.NewScanner(rs[0])
	s.Buffer([]byte{}, 1e12)
	var lines [][]string
	var vals []float64
	for s.Scan() {
		line := strings.Split(s.Text(), "\t")
		if len(line) < 4 {
			return nil, fmt.Errorf("Normalize: line %v has length %v < 4", line, len(line))
		}
		lines = append(lines, line)
		f, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			f = math.NaN()
		}
		vals = append(vals, f)
	}
	vals = NormalizeFloats(vals)
	if len(vals) != len(lines) {
		return nil, fmt.Errorf("Normalize: len(vals) %v != len(lines) %v", len(vals), len(lines))
	}

	var out strings.Builder
	for i, line := range lines {
		line[3] = fmt.Sprintf("%f", vals[i])
		fmt.Fprintln(&out, strings.Join(line, "\t"))
	}
	return strings.NewReader(out.String()), nil
}

func MultiplotSlide(cfg UltimateConfig, winsize, winstep int) error {
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

func AllMultiplotParallel(cfgs []UltimateConfig, winsize, winstep, threads int) error {
	jobs := make(chan UltimateConfig, len(cfgs))
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


