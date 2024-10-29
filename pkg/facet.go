package covplots

import (
	"os"
	"fmt"
	"bufio"
	"io"
	"iter"

	"github.com/jgbaldwinbrown/fastats/pkg"
	"github.com/jgbaldwinbrown/iterh"
	"github.com/jgbaldwinbrown/shellout/pkg"
)

func AddFacetToOneReader(r io.Reader, facetname string) (io.Reader) {
	out := PipeWrite(func(w io.Writer) {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)

		for s.Scan() {
			fmt.Fprintf(w, "%v\t%v\n", s.Text(), facetname)
		}
	})
	return out
}

func AddFacet(r iter.Seq[fastats.BedEntry[[]string]], facetname string) iter.Seq[fastats.BedEntry[[]string]] {
	return func(y func(fastats.BedEntry[[]string]) bool) {
		for b := range r {
			b.Fields = append(b.Fields, facetname)
			if !y(b) {
				return
			}
		}
	}
}

func PlotMultiFacetScales(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiFacetScales\n")
	fmt.Fprintf(os.Stderr, "PlotMultiFacetScales scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facetscales %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacetScalesAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiFacetScalesAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiFacetScalesAny: args %v == \"\"", args)
	}
	return PlotMultiFacetScales(outpre, scalespath)
}

func PlotMultiFacetnameScales(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiFacetnameScales\n")
	fmt.Fprintf(os.Stderr, "PlotMultiFacetnameScales scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facetname_scales %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacetnameScalesAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiFacetnameScalesAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiFacetnameScalesAny: args %v == \"\"", args)
	}
	return PlotMultiFacetnameScales(outpre, scalespath)
}

type PlotMultiFacetScalesBoxedArgs struct {
	Scales string
	Boxes string
}

func PlfmtPath(inpath, outpre string, margs MultiplotPlotFuncArgs) error {
	it := iterh.PathIter(inpath, fastats.ParseBedFlat)
	it2, errp := iterh.BreakWithError(it)

	if !margs.Fullchr {
		var err error
		it2, err = Filter(it2, margs.Chr, margs.Start, margs.End)
		if err != nil {
			return err
		}
		if *errp != nil {
			return *errp
		}
	}

	if err := PlfmtSmallWrite(outpre, it2, margs.Plformatter); err != nil {
		return err
	}
	if *errp != nil {
		return *errp
	}
	return nil
}

func PlotMultiFacetScalesBoxed(outpre string, args PlotMultiFacetScalesBoxedArgs, margs MultiplotPlotFuncArgs) error {
	h := func(e error) error {
		return fmt.Errorf("PlotMultiFacetScalesBoxed: %w", e)
	}

	boxpre := fmt.Sprintf("%v_boxes", outpre)
	boxpath := fmt.Sprintf("%v_boxes_plfmt.bed", outpre)
	err := PlfmtPath(args.Boxes, boxpre, margs)
	if err != nil {
		return h(err)
	}

	fmt.Fprintf(os.Stderr, "running PlotMultiFacetScalesBoxed\n")
	fmt.Fprintf(os.Stderr, "PlotMultiFacetScalesBoxed scalespath: %v\n", args.Scales);
	fmt.Fprintf(os.Stderr, "PlotMultiFacetScalesBoxed boxes path: %v\n", boxpath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facetscales_boxed %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		args.Scales,
		boxpath,
	)

	err = shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return h(err)
	}
	return nil
}
