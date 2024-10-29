package covplots

import (
	"bufio"
	"math"
	"sort"
	"io"
	"fmt"

	"github.com/jgbaldwinbrown/fastats/pkg"
)

func PlotCovVsPair(ylim []float64, args any, margs MultiplotPlotFuncArgs) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric("plot_cov_vs_pair",
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
		)
	}
}

func PlotSelfVsPair(ylim []float64) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric("plot_self_vs_pair",
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
		)
	}
}

func PlotSelfVsPairLim(ylim, xlim []float64) func(outpre string) error {
	return func(outpre string) error {
		return PlotGeneric("plot_self_vs_pair_lim",
			outpre + "_plfmt.bed",
			outpre + "_plotted.png",
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
			fmt.Sprint(xlim[0]),
			fmt.Sprint(xlim[1]),
		)
	}
}

type PlotSelfVsPairArgs struct {
	Xmin float64
	Xmax float64
	Ylab string
	Xlab string
	Width float64
	Height float64
	ResScale float64
	TextSize float64
}

func PlotSelfVsPairPretty(ylim []float64, args PlotSelfVsPairArgs) func(outpre string) error {
	a := args
	return func(outpre string) error {
		return PlotGeneric("plot_self_vs_pair_pretty",
			fmt.Sprintf("%v_plfmt.bed", outpre),
			fmt.Sprintf("%v_plotted.png", outpre),
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
			fmt.Sprint(a.Xmin),
			fmt.Sprint(a.Xmax),
			fmt.Sprint(a.Ylab),
			fmt.Sprint(a.Xlab),
			fmt.Sprint(a.Width),
			fmt.Sprint(a.Height),
			fmt.Sprint(a.ResScale),
			fmt.Sprint(a.TextSize),
		)
	}
}

func PlotSelfVsPairPrettyFixed(ylim []float64, args PlotSelfVsPairArgs) func(outpre string) error {
	a := args
	return func(outpre string) error {
		return PlotGeneric("plot_self_vs_pair_pretty_fixed",
			fmt.Sprintf("%v_plfmt.bed", outpre),
			fmt.Sprintf("%v_plotted.png", outpre),
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
			fmt.Sprint(a.Xmin),
			fmt.Sprint(a.Xmax),
			fmt.Sprint(a.Ylab),
			fmt.Sprint(a.Xlab),
			fmt.Sprint(a.Width),
			fmt.Sprint(a.Height),
			fmt.Sprint(a.ResScale),
			fmt.Sprint(a.TextSize),
		)
	}
}

func PosLess(p1, p2 Pos) bool {
	if p1.Chr < p2.Chr {
		return true
	}
	if p1.Chr > p2.Chr {
		return false
	}
	return p1.Bp < p2.Bp
}

func SortAndUniqPoses[T any](maps ...map[Pos]T) []Pos {
	uniques := map[Pos]struct{}{}
	for _, posmap := range maps {
		for pos, _ := range posmap {
			uniques[pos] = struct{}{}
		}
	}
	poses := make([]Pos, 0, len(uniques))
	for pos, _ := range uniques {
		poses = append(poses, pos)
	}
	sort.Slice(poses, func(i, j int) bool{return PosLess(poses[i], poses[j])})
	return poses
}

type PosEntry struct {
	Pos
	Val float64
}

func CollectEntry(text string, sl *[]PosEntry) error {
	*sl = (*sl)[:0]
	var e fastats.BedEntry[float64]
	_, err := fmt.Sscanf(text, "%s	%d	%d	%f", &e.Chr, &e.Start, &e.End, &e.Fields)
	if err != nil {
		return fmt.Errorf("CollectEntry: %w", err)
	}
	for i:=e.Start; i<e.End; i++ {
		*sl = append(*sl, PosEntry{Pos{e.Chr, int(i)}, e.Fields})
	}
	return nil
}

func InitNaNSl(length int) []float64 {
	out := make([]float64, length, length)
	for i:=0; i<length; i++ {
		out[i] = math.NaN()
	}
	return out
}

func CollectEntries(posmap map[Pos][]float64, idx, nidx int, ebuffer *[]PosEntry, r io.Reader) error {
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		err := CollectEntry(s.Text(), ebuffer)
		if err != nil {
			return err
		}
		for _, e := range *ebuffer {
			vals, ok := posmap[e.Pos]
			if !ok {
				vals = InitNaNSl(nidx)
				posmap[e.Pos] = vals
			}
			vals[idx] = e.Val
		}
	}
	return nil
}

func CombineToOneLine(rs []io.Reader, args any) ([]io.Reader, error) {
	posmap := map[Pos][]float64{}
	var es []PosEntry
	for i, r := range rs {
		err := CollectEntries(posmap, i, len(rs), &es, r)
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(posmap)
	out := PipeWrite(func(w io.Writer) {
		for pos, vals := range posmap {
			fmt.Fprintf(w, "%v\t%v\t%v", pos.Chr, pos.Bp, pos.Bp+1)
			for _, val := range vals {
				fmt.Fprintf(w, "\t%v", val)
			}
			fmt.Fprintf(w, "\n");
		}
	})
	return []io.Reader{out}, nil
}

func CollectEntriesDumb(posmap map[fastats.ChrSpan][]string, idx, nidx int, r io.Reader) error {
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		entry, err := CollectEntryDumb(s.Text())
		if err != nil {
			return err
		}
		vals, ok := posmap[entry.ChrSpan]
		if !ok {
			vals = make([]string, nidx)
			posmap[entry.ChrSpan] = vals
		}
		vals[idx] = entry.Fields
	}
	return nil
}

func CollectEntryDumb(text string) (fastats.BedEntry[string], error) {
	var e fastats.BedEntry[string]
	_, err := fmt.Sscanf(text, "%s	%d	%d	%s", &e.Chr, &e.Start, &e.End, &e.Fields)
	if err != nil { return e, fmt.Errorf("CollectEntryDumb: line %v: %w", text, err) }

	return e, nil
}

func CombineToOneLineDumb(rs []io.Reader, args any) ([]io.Reader, error) {
	if len(rs) < 1 {
		return []io.Reader{}, nil
	}
	fmt.Println("CombineToOneLineDumb len(rs):", len(rs))

	posmap := map[fastats.ChrSpan][]string{}
	for i, r := range rs {
		err := CollectEntriesDumb(posmap, i, len(rs), r)
		if err != nil {
			return nil, err
		}
	}

	out := PipeWrite(func(w io.Writer) {
		for span, vals := range posmap {
			fmt.Println("printing span", span, "with vals", vals)
			fmt.Fprintf(w, "%v\t%v\t%v", span.Chr, span.Start, span.End)
			for _, val := range vals {
				fmt.Fprintf(w, "\t%v", val)
			}
			fmt.Fprintf(w, "\n");
		}
	})
	return []io.Reader{out}, nil
}

type PlotCovHistArgs struct {
	Xmin float64
	Xmax float64
	Ylab string
	Xlab string
	Binwidth float64
	Width float64
	Height float64
	ResScale float64
}

func PlotCovHist(ylim []float64, args PlotCovHistArgs) func(outpre string) error {
	hargs := args
	return func(outpre string) error {
		return PlotGeneric("plot_cov_hist",
			fmt.Sprintf("%v_plfmt.bed", outpre),
			fmt.Sprintf("%v_plotted.png", outpre),
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
			fmt.Sprint(hargs.Xmin),
			fmt.Sprint(hargs.Xmax),
			fmt.Sprint(hargs.Ylab),
			fmt.Sprint(hargs.Xlab),
			fmt.Sprint(hargs.Binwidth),
			fmt.Sprint(hargs.Width),
			fmt.Sprint(hargs.Height),
			fmt.Sprint(hargs.ResScale),
		)
	}
}

type PlotBoxwhiskerArgs struct {
	Xmin float64
	Xmax float64
	Ylab string
	Xlab string
	Width float64
	Height float64
	ResScale float64
	TextSize float64
	FillName string
}

func PlotBoxwhisker(ylim []float64, args PlotBoxwhiskerArgs) func(outpre string) error {
	a := args
	return func(outpre string) error {
		return PlotGeneric("plot_boxwhisker",
			fmt.Sprintf("%v_plfmt.bed", outpre),
			fmt.Sprintf("%v_plotted.pdf", outpre),
			fmt.Sprintf("%v_plotted.png", outpre),
			fmt.Sprint(ylim[0]),
			fmt.Sprint(ylim[1]),
			fmt.Sprint(a.Xmin),
			fmt.Sprint(a.Xmax),
			fmt.Sprint(a.Ylab),
			fmt.Sprint(a.Xlab),
			fmt.Sprint(a.Width),
			fmt.Sprint(a.Height),
			fmt.Sprint(a.ResScale),
			fmt.Sprint(a.TextSize),
			fmt.Sprint(a.FillName),
		)
	}
}
