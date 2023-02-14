package covplots

import (
	"github.com/jgbaldwinbrown/shellout/pkg"
	"os"
	"fmt"
	"bufio"
	"io"
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

func AddFacet(rs []io.Reader, args any) ([]io.Reader, error) {
	var out []io.Reader
	anysl, ok := args.([]any)
	if !ok {
		return nil, fmt.Errorf("AddFacet: input args %v not []any", args)
	}

	var facetnames []string
	for _, arg := range anysl {
		name, ok := arg.(string)
		if !ok {
			return nil, fmt.Errorf("AddFacet: input arg %v not string", arg)
		}
		facetnames = append(facetnames, name)
	}

	for i, r := range rs {
		out = append(out, AddFacetToOneReader(r, facetnames[i]))
	}
	return out, nil
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

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacetScalesAny(outpre string, ylim []float64, args any) error {
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

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacetnameScalesAny(outpre string, ylim []float64, args any) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiFacetnameScalesAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiFacetnameScalesAny: args %v == \"\"", args)
	}
	return PlotMultiFacetnameScales(outpre, scalespath)
}

