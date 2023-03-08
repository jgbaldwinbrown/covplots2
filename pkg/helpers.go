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
	chrmins := make(map[string]int)
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

		min, ok := chrmins[entry.Chr]
		if !ok {
			chrmins[entry.Chr] = entry.Start
			min = entry.Start
		}
		if min > entry.Start {
			chrmins[entry.Chr] = entry.Start
		}

		data = append(data, entry)
	}

	if len(chrs) < 1 {
		return nil
	}


	// chrs = []string{"2L", "2R", "3L", "3R"}

	bpused := []int{chrlens[chrs[0]] - chrmins[chrs[0]]}
	chrnums := make(map[string]int)
	chrnums[chrs[0]] = 0
	chroffs := make(map[string]int)
	chroffs[chrs[0]] = -chrmins[chrs[0]]

	/*
	idx	start	end	fstart	fend	offset
	0	5	30	0	25	-5
	1	15	20	25	30	25 - 15
	2	2	4	30	32	30 - 2
	*/

	for i:=1; i<len(chrs); i++ {
		chr := chrs[i]
		chrnums[chr] = i

		bpused = append(bpused, bpused[i-1] + chrlens[chr] - chrmins[chr])
		// offsets = append(offsets, offsets[i-1] + chrlens[chrs[i-1]] - chrmins[chrs[i-1] - chrmins[chrs[i]])
		chroffs[chr] = bpused[i-1] - chrmins[chr]
	}

	// chrset := map[string]struct{}{
	// 	"2L": struct{}{},
	// 	"2R": struct{}{},
	// 	"3L": struct{}{},
	// 	"3R": struct{}{},
	// }
	for _, e := range data {
		// if _, ok := chrset[e.Chr]; !ok {
		// 	continue
		// }
		e.StartOff = chroffs[e.Chr] + e.Start
		e.EndOff = chroffs[e.Chr] + e.End
		e.ChrNum = chrnums[e.Chr]
		fmt.Fprintf(bw, "%s\t%d\t%d\t%d\n", e.Text, e.ChrNum, e.StartOff, e.EndOff)
	}
	return nil
}

func PlotSingle(outpre string, subtract bool) error {
	subtxt := ""
	if subtract {
		subtxt = "_sub"
	}
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot%s_single_cov %v %v
`,
		subtxt,
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
		line = lscan.SplitByFunc(line, s, splitter)
		if len(line) < 3 {
			return false
		}
		if !re.MatchString(line[0]) {
			return false
		}

		if end != -1 {
			lstart, err := strconv.ParseInt(line[1], 0, 64)
			if err != nil {
				return false
			}
			if int(lstart) >= end {
				return false
			}
		}

		if start != -1 {
			lend, err := strconv.ParseInt(line[2], 0, 64)
			if err != nil {
				return false
			}
			if int(lend) <= start {
				return false
			}
		}

		return true
	}
	return f, nil
}

func ReChr(rs []io.Reader, abiolines any) ([]io.Reader, error) {
	biolines, ok := abiolines.([]string)
	if !ok {
		return nil, fmt.Errorf("abiolines %v not of type []string", abiolines)
	}
	var outs []io.Reader
	for _, r := range rs {
		outs = append(outs, ReChrSingle(r, biolines))
	}
	return outs, nil
}

func ReChrSingle(r io.Reader, biolines []string) (io.Reader) {
	chrre := regexp.MustCompile(`^[^	]*`)
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	rout := PipeWrite(func(w io.Writer) {
		for s.Scan() {
			out := s.Text()
			for _, l := range biolines {
				out = chrre.ReplaceAllString(out, `&` + "_" + l)
			}
			fmt.Println(w, out)
		}
	})
	return rout
}

func ChrGrep(rs []io.Reader, apattern any) ([]io.Reader, error) {
	pattern, ok := apattern.(string)
	if !ok {
		return nil, fmt.Errorf("pattern %v not of type string", apattern)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("ChrGrep: could not compile pattern %v with error %w", pattern, err)
	}

	var outs []io.Reader
	for _, r := range rs {
		outs = append(outs, ChrGrepSingle(r, re))
	}
	return outs, nil
}

func ChrGrepSingle(r io.Reader, re *regexp.Regexp) (io.Reader) {
	chrre := regexp.MustCompile(`^[^	]*`)
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	rout := PipeWrite(func(w io.Writer) {
		i := 0
		j := 0
		for s.Scan() {
			chrstr := chrre.FindString(s.Text())
			if re.MatchString(chrstr) {
				fmt.Fprintln(w, s.Text())
				j++
			}
			i++
		}
		fmt.Fprintf(os.Stderr, "ChrGrep: printed %v of %v lines\n", j, i)
	})
	return rout
}

func RunSingle() {
	f := GetFlags()

	err := SinglePlot(os.Stdin, f.Outpre, f.Chr, f.Start, f.End)
	if err != nil {
		panic(err)
	}
}
