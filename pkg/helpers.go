package covplots

import (
	"strconv"
	"regexp"
	"bufio"
	"bytes"
	"flag"
	"os"
	"github.com/jgbaldwinbrown/lscan/pkg"
	"github.com/jgbaldwinbrown/shellout/pkg"
	"io"
	"fmt"
)

func Plfmt(r io.Reader, outpre, chrbedpath string) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plfmt_flex -c 0 -b 1 -b2 2 -C %v -n > %v
`,
		chrbedpath,
		fmt.Sprintf("%v_plfmt.bed", outpre),
)

	return shellout.ShellOutPiped(script, r, os.Stdout, os.Stderr)
}

func PlotSingle(outpre string) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_single_cov %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotWin(outpre string) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_window_cov %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

type Flags struct {
	Outpre string
	Chrbedpath string
	Chr string
	Start int
	End int
}

func GetFlags() Flags {
	var f Flags
	flag.StringVar(&f.Outpre, "o", "single_cov_plot", "Output prefix")
	flag.StringVar(&f.Chrbedpath, "C", "", "chromosome lengths bed path")
	flag.StringVar(&f.Chr, "c", "", "chromosome to plot")
	flag.IntVar(&f.Start, "s", -1, "Starting coordinate to plot")
	flag.IntVar(&f.End, "e", -1, "End coordinate to plot")
	flag.Parse()

	if f.Chrbedpath == "" {
		panic(fmt.Errorf("missing chrbedpath"))
	}

	return f
}

type Filterer struct {
	s *bufio.Scanner
	buf *bytes.Buffer
	filt func(string) bool
}

func (f *Filterer) Read(out []byte) (n int, err error) {
	n, err = f.buf.Read(out)
	for n < len(out) {
		if !f.s.Scan() {
			return n, io.EOF
		}
		if f.filt(f.s.Text()) {
			f.buf.WriteString(f.s.Text() + "\n")
		}
		var n2 int
		n2, err = f.buf.Read(out[n:])
		n += n2
	}
	fmt.Printf("writing %v chars\n", n)
	return n, err
}

func Filter(r io.Reader, chr string, start, end int) (*Filterer, error) {
	re, err := regexp.Compile("^" + chr + "_")
	if err != nil {
		return nil, err
	}

	f := new(Filterer)
	f.s = bufio.NewScanner(r)
	f.buf = bytes.NewBuffer([]byte{})

	var line []string
	splitter := lscan.ByByte('\t')
	f.filt = func(s string) bool {
		// fmt.Println("matching")
		line = lscan.SplitByFunc(line, s, splitter)
		if len(line) < 3 {
			// fmt.Println("too short")
			return false
		}
		if !re.MatchString(line[0]) {
			// fmt.Println("wrong chr")
			return false
		}

		if end != -1 {
			// fmt.Println("checking end")
			lstart, err := strconv.ParseInt(line[1], 0, 64)
			if err != nil {
				// fmt.Println("couldn't parse start")
				return false
			}
			if int(lstart) >= end {
				// fmt.Println("lstart >= end")
				return false
			}
		}

		if start != -1 {
			// fmt.Println("checking start")
			lend, err := strconv.ParseInt(line[2], 0, 64)
			if err != nil {
				// fmt.Println("couldn't parse end")
				return false
			}
			if int(lend) <= start {
				// fmt.Println("lend <= start")
				return false
			}
		}

		fmt.Println("matched")
		return true
	}
	return f, nil
}

func RunSingle() {
	f := GetFlags()

	r, err := Filter(os.Stdin, f.Chr, f.Start, f.End)
	if err != nil {
		panic(err)
	}

	err = Plfmt(r, f.Outpre, f.Chrbedpath)
	if err != nil {
		panic(err)
	}

	err = PlotSingle(f.Outpre)
	if err != nil {
		panic(err)
	}
}
