package covplots

import (
	"iter"
	"regexp"
	"errors"
	"os/exec"
	"math"
	"os"
	"fmt"

	"golang.org/x/sync/errgroup"
	"github.com/montanaflynn/stats"
	"github.com/jgbaldwinbrown/fastats/pkg"
	"github.com/jgbaldwinbrown/iterh"
)

func FilterMulti[B fastats.ChrSpanner](chr string, start, end int, rs ...iter.Seq[B]) ([]iter.Seq[B], error) {
	var out []iter.Seq[B]
	for _, r := range rs {
		fr, err := Filter(r, chr, start, end)
		if err != nil {
			return nil, fmt.Errorf("FilterMulti: %w", err)
		}
		out = append(out, fr)
	}
	return out, nil
}

func CombineSinglebpPlots[B fastats.BedEnter[[]string]](names []string, rs ...iter.Seq[B]) iter.Seq[fastats.BedEntry[[]string]] {
	return func(y func(fastats.BedEntry[[]string]) bool) {
		for i, r := range rs {
			for b := range r {
				ent := fastats.ToBedEntry(b)
				ent.Fields = append(ent.Fields, names[i])
				if !y(ent) {
					return
				}
			}
		}
	}
}

func PlotGeneric(command, bedpath, outpath string, args ...string) error {
	full := make([]string, 0, 3 + len(args))
	full = append(full, command, bedpath, outpath)
	full = append(full, args...)
	cmd := exec.Command(full[0], full[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func PlotMulti(ylim []float64) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric("plot_singlebp_multiline_cov",
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
		)
	}
}

func PlotMultiFixedOrder(ylim []float64) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric("plot_singlebp_multiline_cov_fixed_order",
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
		)
	}
}

type PrettyCfg struct {
	Xlab string
	Ylab string
	Width float64
	Height float64
	Res float64
	TextSize float64
}

func PlotMultiPrettyGeneric(command string, ylim []float64, cfg PrettyCfg) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric(command,
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
			fmt.Sprint(cfg.Xlab),
			fmt.Sprint(cfg.Ylab),
			fmt.Sprint(cfg.Width),
			fmt.Sprint(cfg.Height),
			fmt.Sprint(cfg.Res),
			fmt.Sprint(cfg.TextSize),
		)
	}
}

func PlotMultiPretty(outpre string, ylim []float64, cfg PrettyCfg) func(outpre string) error {
	return PlotMultiPrettyGeneric("plot_multi_pretty", ylim, cfg)
}

func PlotMultiPrettyBlue(outpre string, ylim []float64, cfg PrettyCfg) func(outpre string) error {
	return PlotMultiPrettyGeneric("plot_multi_pretty_blue", ylim, cfg)
}

func PlotMultiPrettyColorseries(outpre string, ylim []float64, cfg PrettyCfg) func(outpre string) error {
	return PlotMultiPrettyGeneric("plot_multi_pretty_colorseries", ylim, cfg)
}

func PlotMultiFacet(outpre string, ylim []float64) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric("plot_singlebp_multiline_cov_facet",
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
		)
	}
}

func CheckPathExists(path string) bool {
	_, err := os.Stat("/path/to/whatever")
	return !errors.Is(err, os.ErrNotExist)
}

func GzPath(path string, threads int) error {
	cmd := exec.Command("pigz", "-f", "-p", fmt.Sprintf("%d", threads), path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StripParent[F any](it iter.Seq[fastats.BedEntry[F]]) iter.Seq[fastats.BedEntry[F]] {
	return func(y func(fastats.BedEntry[F]) bool) {
		re := regexp.MustCompile(`^([^_	]*)_([^	])*`)
		for b := range it {
			b.Chr = re.ReplaceAllString(b.Chr, "$1")
			if !y(b) {
				return
			}
		}
	}
}

func GetManualChrs(path string) (chrs []string, err error) {
	it := iterh.PathIter(path, fastats.ParseBedFlat)
	for b, e := range it {
		chrs = append(chrs, b.Chr)
		if e != nil {
			return nil, e
		}
	}

	return chrs, nil
}

type MultiplotPlotFuncArgs struct {
	Plformatter *Plformatter
	Cfg UltimateConfig
	Chr string
	Start int
	End int
	Fullchr bool
}

func Multiplot[C fastats.ChrSpanner](cfg UltimateConfig, chr string, start, end int) (func() error, *Plformatter, error) {
	outpre := fmt.Sprintf("%s_%v_%v_%v", cfg.Outpre, chr, start, end)
	bed := cfg.Input
	if !cfg.FullChr {
		var err error
		bed, err = Filter(bed, chr, start, end)
		if err != nil {
			return nil, nil, err
		}
	}

	if cfg.NoParent {
		bed = StripParent(bed)
	}

	pf := PlfmtSmallRead(bed, cfg.ManualChrs, cfg.UseManualChrs)

	f := func() error {
		if err := PlfmtSmallWrite(outpre + "_plfmt.bed", bed, pf); err != nil {
			return fmt.Errorf("Multiplot: during PlfmtSmallWrite: %w", err)
		}
		if err := cfg.PlotFunc(outpre); err != nil {
			return fmt.Errorf("Multiplot: during plotfunc: %w", err)
		}
		if err := GzPath(outpre + "_plfmt.bed", 8); err != nil {
			return fmt.Errorf("Multiplot: during GzPath: %w", err)
		}
		return nil
	}

	return f, pf, nil
}

func NormalizeFloats(in []float64) {
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
	for i, f := range in {
		in[i] = (f-m) / s
	}
}

func Normalize[B fastats.BedEnter[float64]](r iter.Seq[B]) iter.Seq[fastats.BedEntry[float64]] {
	return func(y func(fastats.BedEntry[float64]) bool) {
		var bs []fastats.BedEnter[float64]
		var vals []float64
		for b := range r {
			bs = append(bs, b)
			vals = append(vals, b.BedFields())
		}
		NormalizeFloats(vals)

		for i, b := range bs {
			ent := fastats.ToBedEntry(b)
			ent.Fields = vals[i]
			if !y(ent) {
				return
			}
		}
	}
}

func MultiplotFullchr(cfg UltimateConfig) error {
	f, _, err := Multiplot[fastats.ChrSpan](cfg, "full_genome", 0, 0)
	if err != nil {
		return err
	}
	if err := f(); err != nil {
		return fmt.Errorf("MultiplotFullchr: %w", err)
	}

	return nil
}

func MultiplotSelectWins(cfg UltimateConfig) iter.Seq2[func() error, error] {
	return func(y func(func() error, error) bool) {
		for entry := range cfg.SelectedWins {
			f, _, err := Multiplot[fastats.ChrSpan](cfg, entry.SpanChr(), int(entry.SpanStart()), int(entry.SpanEnd()))
			if !y(f, err) {
				return
			}
		}
	}
}

func MultiplotSlide(cfg UltimateConfig) (iter.Seq2[func() error, error], error) {
	chrlens, err := GetChrLens(cfg.Chrlens)
	if err != nil {
		return nil, fmt.Errorf("MultiplotSlide: %w", err)
	}

	return func(y func(func() error, error) bool) {
		for _, chrlenset := range chrlens {
			chr, chrlen := chrlenset.Chr, chrlenset.Len
			for start := 0; start < chrlen; start += cfg.Winstep {
				end := start + cfg.Winsize
				f, _, err := Multiplot[fastats.ChrSpan](cfg, chr, start, end)
				if !y(f, err) {
					return
				}
			}
		}
	}, nil
}

func MultiplotFlex(cfg UltimateConfig) (iter.Seq2[func() error, error], error) {
	if cfg.FullChr {
		return func(y func(func() error, error) bool) {
			f := func() error {
				return MultiplotFullchr(cfg)
			}
			if !y(f, nil) {
				return
			}
		}, nil
	}
	if cfg.UseSelectedWins {
		return MultiplotSelectWins(cfg), nil
	}
	return MultiplotSlide(cfg)
}

func AllMultiplotParallel(cfgs []UltimateConfig, threads int) error {
	var g errgroup.Group
	if threads > 0 {
		g.SetLimit(threads)
	}

	for _, cfg := range cfgs {
		fs, e := MultiplotFlex(cfg)
		if e != nil {
			return e
		}
		for f, e := range fs {
			if e != nil {
				g.Wait()
				return e
			}
			g.Go(f)
		}
	}
	return g.Wait()
}
