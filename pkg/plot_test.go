package covplots

import (
	"testing"

	"github.com/jgbaldwinbrown/fastats/pkg"
)

func exampleBedEntries(y func(fastats.BedEntry[[]string]) bool) {
	b := fastats.BedEntry[[]string]{
		ChrSpan: fastats.ChrSpan{Chr: "chr1", Span: fastats.Span{Start: 5, End: 6}},
		Fields: []string{"0.35", "i1"},
	}
	if !y(b) {
		return
	}
	b = fastats.BedEntry[[]string]{
		ChrSpan: fastats.ChrSpan{Chr: "chr1", Span: fastats.Span{Start: 7, End: 8}},
		Fields: []string{"0.45", "i1"},
	}
	if !y(b) {
		return
	}
	b = fastats.BedEntry[[]string]{
		ChrSpan: fastats.ChrSpan{Chr: "chr2", Span: fastats.Span{Start: 3, End: 9}},
		Fields: []string{"0.25", "i1"},
	}
	if !y(b) {
		return
	}
}

func TestPlot(t *testing.T) {
	c := UltimateConfig{}
	c.Input = exampleBedEntries
	c.Chrlens = "example_chrlens.bed"
	c.Outpre = "example_out/example_out"
	c.PlotFunc = PlotMulti([]float64{0, 1})
	c.Winsize = 100
	c.Winstep = 5

	if err := AllMultiplotParallel([]UltimateConfig{c}, 8); err != nil {
		t.Errorf("TestPlot: %v", err)
	}
}
