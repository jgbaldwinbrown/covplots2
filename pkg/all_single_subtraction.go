package covplots

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"iter"
	"strings"

	"github.com/jgbaldwinbrown/fastats/pkg"
	"github.com/jgbaldwinbrown/iterh"
)

func GetAllSubtractSingleFlags() AllSingleFlags {
	var f AllSingleFlags
	flag.StringVar(&f.Config, "i", "", "Input config file. Tab-separated columns containing input bed path 1, input bed path 2, chromosome length bed path, and output prefix. Default stdin.")
	flag.IntVar(&f.WinSize, "w", 1000000, "Sliding window plot size (default = 1000000).")
	flag.IntVar(&f.WinStep, "s", 1000000, "Sliding window step distance (default = 1000000).")
	flag.Parse()

	return f
}

func SubtractSinglePlotWinsParallel(cfgs []Config, winsize, winstep, threads int) error {
	jobs := make(chan Config, len(cfgs))
	for _, cfg := range cfgs {
		jobs <- cfg
	}
	close(jobs)

	errs := make(chan error, len(cfgs))

	for i := 0; i < threads; i++ {
		go func() {
			for cfg := range jobs {
				errs <- SubtractSinglePlotWins(cfg.Inpath, cfg.Inpath2, cfg.Chrlenpath, cfg.Outpre, winsize, winstep)
			}
		}()
	}

	var out Errors
	for i := 0; i < len(cfgs); i++ {
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

func RunAllSubtractSinglePlots() {
	f := GetAllSubtractSingleFlags()
	cfgs, err := GetConfig(f.Config, true)
	if err != nil {
		panic(err)
	}

	err = SubtractSinglePlotWinsParallel(cfgs, f.WinSize, f.WinStep, 8)
	if err != nil {
		panic(err)
	}
}

func SubtractSinglePlotWins(inpath1, inpath2, chrlenpath, outpre string, winsize, winstep int) error {
	chrlens, err := GetChrLens(chrlenpath)
	if err != nil {
		return fmt.Errorf("SubtactSinglePlotWins: %w", err)
	}

	for _, chrlenset := range chrlens {
		chr, chrlen := chrlenset.Chr, chrlenset.Len
		for start := 0; start < chrlen; start += winstep {
			end := start + winsize
			outpre2 := fmt.Sprintf("%s_%v_%v_%v", outpre, chr, start, end)
			err = SubtractSinglePlotPath(inpath1, inpath2, outpre2, chr, start, end)
			if err != nil {
				return fmt.Errorf("SubtractSinglePlotWins: %w", err)
			}
		}
	}

	return nil
}

func SubtractSinglePlotPath(path1, path2 string, outpre, chr string, start, end int) error {
	it1, ep1 := iterh.BreakWithError(iterh.PathIter(path1, ParseBedgraph))
	it2, ep2 := iterh.BreakWithError(iterh.PathIter(path2, ParseBedgraph))

	err := SubtractSinglePlot(it1, it2, outpre, chr, start, end)
	if err != nil {
		return fmt.Errorf("SubtractSinglePlotPath: %w", err)
	}
	if *ep1 != nil {
			return *ep1
	}
	if *ep2 != nil {
			return *ep2
	}
	return nil
}

type Pos struct {
	Chr string
	Bp  int
}

type SubVal struct {
	Val        float64
	Subtracted bool
}

func AllPosVals[B fastats.BedEnter[float64]](b B) iter.Seq[PosEntry] {
	start := b.SpanStart()
	end := b.SpanEnd()
	chr := b.SpanChr()
	return func(y func(PosEntry) bool) {
		for i := start; i < end; i++ {
			if !y(PosEntry{Pos{chr, int(i)}, b.BedFields()}) {
				return
			}
		}
	}
}

func CollectVals(r io.Reader) (map[Pos]float64, error) {
	out := make(map[Pos]float64)
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		var chr string
		var start int
		var end int
		var v float64
		_, err := fmt.Sscanf(s.Text(), "%s	%d	%d	%f", &chr, &start, &end, &v)
		if err != nil {
			return nil, fmt.Errorf("CollectVals: %w", err)
		}
		for i := start; i < end; i++ {
			out[Pos{chr, i}] = v
		}
	}
	return out, nil
}

func SubtractInternal[B fastats.BedEnter[float64]](it1, it2 iter.Seq[B]) map[Pos]SubVal {
	out := map[Pos]SubVal{}
	for b := range it1 {
		posvals := AllPosVals(b)
		for pv := range posvals {
			out[pv.Pos] = SubVal{pv.Val, false}
		}
	}

	for b := range it2 {
		posvals := AllPosVals(b)
		for pv2 := range posvals {
			if sv1, ok := out[pv2.Pos]; ok {
				out[pv2.Pos] = SubVal{sv1.Val - pv2.Val, true}
			}
		}
	}
	return out
}

func Subtract[B fastats.BedEnter[float64]](it1, it2 iter.Seq[B]) iter.Seq[fastats.BedEntry[float64]] {
	return func(y func(fastats.BedEntry[float64]) bool) {
		sub := SubtractInternal(it1, it2)

		for pos, sval := range sub {
			if sval.Subtracted {
				be := fastats.BedEntry[float64]{
					ChrSpan: fastats.ChrSpan{
						Span: fastats.Span{Start: int64(pos.Bp), End: int64(pos.Bp + 1)},
						Chr:  pos.Chr,
					},
					Fields: sval.Val,
				}
				if !y(be) {
					return
				}
			}
		}
	}
}

func SubtractOld(r1, r2 io.Reader) (*strings.Reader, error) {
	vals1, err := CollectVals(r1)
	if err != nil {
		return nil, fmt.Errorf("Subtract: %w", err)
	}
	vals2, err := CollectVals(r2)
	if err != nil {
		return nil, fmt.Errorf("Subtract: %w", err)
	}
	sub := make(map[Pos]float64)
	for pos, val := range vals1 {
		if val2, ok := vals2[pos]; ok {
			sub[pos] = val - val2
		}
	}

	var out strings.Builder
	for pos, val := range sub {
		fmt.Fprintf(&out, "%s\t%d\t%d\t%f\n", pos.Chr, pos.Bp, pos.Bp+1, val)
	}
	return strings.NewReader(out.String()), nil
}

func SubtractSinglePlot(r1, r2 iter.Seq[fastats.BedEntry[float64]], outpre, chr string, start, end int) error {
	fr1, err := Filter(r1, chr, start, end)
	if err != nil {
		return err
	}

	fr2, err := Filter(r2, chr, start, end)
	if err != nil {
		return err
	}

	fr3 := Subtract(fr1, fr2)
	fr4 := func(y func(fastats.BedEntry[[]string]) bool) {
		for b := range fr3 {
			b2 := fastats.BedEntry[[]string]{}
			b2.ChrSpan = b.ChrSpan
			b2.Fields = []string{fmt.Sprintf("%v", b.Fields)}
			if !y(b2) {
				return
			}
		}
	}

	_, err = PlfmtSmall(fr4, outpre, nil, false)
	if err != nil {
		return err
	}

	err = PlotSingle(outpre)
	if err != nil {
		return err
	}
	return nil
}
