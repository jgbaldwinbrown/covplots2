package covplots

import (
	"strings"
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

type PlfmtEntry struct {
	Chr string
	Start int
	End int
	Text string
	Line []string
	ChrNum int
	StartOff int
	EndOff int
}

func PlfmtSmall(r io.Reader, outpre string) error {
	w, err := os.Create(outpre + "_plfmt.bed")
	if err != nil {
		return err
	}
	defer w.Close()
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	data := []PlfmtEntry{}
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	chrlens := make(map[string]int)
	chrs := []string{}

	for s.Scan() {
		line := strings.Split(s.Text(), "\t")
		if len(line) < 3 {
			continue
		}
		start, err := strconv.ParseInt(line[1], 0, 64)
		if err != nil {
			continue
		}
		end, err := strconv.ParseInt(line[2], 0, 64)
		if err != nil {
			continue
		}
		entry := PlfmtEntry{Chr: line[0], Start: int(start), End: int(end), Text: s.Text(), Line: line}
		length, ok := chrlens[entry.Chr]
		if !ok {
			chrlens[entry.Chr] = 0
			length = 0
			chrs = append(chrs, entry.Chr)
		}
		if entry.End > length {
			chrlens[entry.Chr] = entry.End
		}
		data = append(data, entry)
	}

	if len(chrs) < 1 {
		return nil
	}

	offsets := []int{0}
	chrnums := make(map[string]int)
	chrnums[chrs[0]] = 0
	chroffs := make(map[string]int)
	chroffs[chrs[0]] = 0
	for i:=1; i<len(chrs); i++ {
		chr := chrs[i]
		offsets = append(offsets, offsets[i-1] + chrlens[chrs[i-1]])
		chrnums[chr] = i
		chroffs[chr] = offsets[i]
	}
	for _, e := range data {
		e.StartOff = chroffs[e.Chr] + e.Start
		e.EndOff = chroffs[e.Chr] + e.End
		e.ChrNum = chrnums[e.Chr]
		fmt.Fprintf(bw, "%s\t%d\t%d\t%d\n", e.Text, e.ChrNum, e.StartOff, e.EndOff)
	}
	return nil
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
	// fmt.Printf("writing %v chars\n", n)
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

		// fmt.Println("matched")
		return true
	}
	return f, nil
}

func RunSingle() {
	f := GetFlags()

	err := SinglePlot(os.Stdin, f.Outpre, f.Chr, f.Start, f.End)
	if err != nil {
		panic(err)
	}
}
