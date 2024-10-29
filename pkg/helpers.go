package covplots

import (
	"slices"
	"os/exec"
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"iter"
	"os"
	"regexp"
	"strings"

	"github.com/jgbaldwinbrown/csvh"
	"github.com/jgbaldwinbrown/fastats/pkg"
	"github.com/jgbaldwinbrown/iterh"
)

func ParseBedgraphEntry(line string) (fastats.BedEntry[float64], error) {
	b := fastats.BedEntry[float64]{}
	l := strings.Split(line, "\t")
	_, e := csvh.Scan(l, &b.Chr, &b.Start, &b.End, &b.Fields)
	return b, e
}

func ParseBedgraph(r io.Reader) iter.Seq2[fastats.BedEntry[float64], error] {
	return func(y func(fastats.BedEntry[float64], error) bool) {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		for s.Scan() {
			if s.Err() != nil {
				if !y(fastats.BedEntry[float64]{}, s.Err()) {
					return
				}
			}
			b, e := ParseBedgraphEntry(s.Text())
			if !y(b, e) {
				return
			}
		}
	}
}

type Plformatter struct {
	UseManualChrs bool
	Chrset        map[string]struct{}
	Chroffs       map[string]int
	Chrnums       map[string]int
	Nfields       int
}

func PlfmtSmall(it iter.Seq[fastats.BedEntry[[]string]], outpath string, manualChrs iter.Seq[string], useManualChrs bool) (*Plformatter, error) {
	out := PlfmtSmallRead(it, manualChrs, useManualChrs)

	if err := PlfmtSmallWrite(outpath, it, out); err != nil {
		return nil, err
	}
	return out, nil
}

func PlfmtSmallRead[B fastats.BedEnter[[]string]](it iter.Seq[B], manualChrs iter.Seq[string], useManualChrs bool) *Plformatter {
	out := &Plformatter{
		UseManualChrs: useManualChrs,
		Chrset:        map[string]struct{}{},
		Chroffs:       map[string]int{},
		Chrnums:       map[string]int{},
	}

	chrlens := make(map[string]int)
	chrmins := make(map[string]int)

	chrs := []string{}

	for b := range it {
		length, ok := chrlens[b.SpanChr()]
		if !ok {
			chrlens[b.SpanChr()] = 0
			length = 0
			chrs = append(chrs, b.SpanChr())
		}
		if int(b.SpanEnd()) > length {
			chrlens[b.SpanChr()] = int(b.SpanEnd())
		}

		cmin, ok := chrmins[b.SpanChr()]
		if !ok {
			chrmins[b.SpanChr()] = int(b.SpanStart())
			cmin = int(b.SpanStart())
		}
		if cmin > int(b.SpanStart()) {
			chrmins[b.SpanChr()] = int(b.SpanStart())
		}
		if out.Nfields < len(b.BedFields()) {
			out.Nfields = len(b.BedFields())
		}
	}

	if len(chrs) < 1 {
		return out
	}

	if useManualChrs {
		chrs = slices.Collect(manualChrs)
	}

	bpused := []int{chrlens[chrs[0]] - chrmins[chrs[0]]}
	out.Chrnums[chrs[0]] = 0
	out.Chroffs[chrs[0]] = -chrmins[chrs[0]]

	for i := 1; i < len(chrs); i++ {
		chr := chrs[i]
		out.Chrnums[chr] = i

		bpused = append(bpused, bpused[i-1]+chrlens[chr]-chrmins[chr])
		out.Chroffs[chr] = bpused[i-1] - chrmins[chr]
	}

	if useManualChrs {
		for _, chr := range chrs {
			out.Chrset[chr] = struct{}{}
		}
	}

	return out
}

func PlfmtSmallWrite(outpath string, it iter.Seq[fastats.BedEntry[[]string]], f *Plformatter) (err error) {
	h := func(e error) error {
		return fmt.Errorf("PlfmtSmallWrite: %w", e)
	}

	w, err := os.Create(outpath)
	if err != nil {
		return h(err)
	}
	defer func() { csvh.DeferE(&err, w.Close()) }()
	bw := bufio.NewWriter(w)
	defer func() { csvh.DeferE(&err, bw.Flush()) }()

	for b := range it {
		if f.UseManualChrs {
			if _, ok := f.Chrset[b.Chr]; !ok {
				continue
			}
		}

		if _, e := fmt.Fprintf(bw, "%v\t%v\t%v", b.Chr, b.Start, b.End); e != nil {
			return h(e)
		}
		for _, field := range b.Fields {
			if _, e := fmt.Fprintf(bw, "\t%v", field); e != nil {
				return h(e)
			}
		}
		if len(b.Fields) < f.Nfields {
			for i := 0; i < f.Nfields - len(b.Fields); i++ {
				if _, e := fmt.Fprintf(bw, "\t"); e != nil {
					return h(e)
				}
			}
		}
		if _, err := fmt.Fprintf(bw, "\t%d\t%d\t%d\n", f.Chrnums[b.Chr], f.Chroffs[b.Chr] + int(b.Start), f.Chroffs[b.Chr] + int(b.End)); err != nil {
			return h(err)
		}
	}
	return nil
}

func PlotSingle(outpre string) error {
	cmd := exec.Command("plot_single_cov", fmt.Sprintf("%v_plfmt.bed", outpre), fmt.Sprintf("%v_plotted.png", outpre))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func PlotWin(outpre string) error {
	cmd := exec.Command("plot_window_cov", fmt.Sprintf("%v_plfmt.bed", outpre), fmt.Sprintf("%v_plotted.png", outpre))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type Flags struct {
	Outpre     string
	Chrbedpath string
	Chr        string
	Start      int
	End        int
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

func Filter[B fastats.ChrSpanner](it iter.Seq[B], chr string, start, end int) (iter.Seq[B], error) {
	re, err := regexp.Compile("^" + chr + "_|$")
	if err != nil {
		return nil, err
	}
	return func(y func(B) bool) {
		for b := range it {
			if !re.MatchString(b.SpanChr()) || int(b.SpanStart()) >= end || int(b.SpanEnd()) < start {
				continue
			}
			if !y(b) {
				return
			}
		}
	}, nil
}

func ReChr[F any](it iter.Seq[fastats.BedEntry[F]], biolines []string) iter.Seq[fastats.BedEntry[F]] {
	return func(y func(fastats.BedEntry[F]) bool) {
		chrre := regexp.MustCompile(`^[^	]*`)
		for b := range it {
			for _, l := range biolines {
				b.Chr = chrre.ReplaceAllString(b.Chr, `&`+"_"+l)
			}
			if !y(b) {
				return
			}
		}
	}
}

func ChrGrep[B fastats.ChrSpanner](it iter.Seq[B], re *regexp.Regexp) iter.Seq[B] {
	return func(y func(B) bool) {
		for b := range it {
			if re.MatchString(b.SpanChr()) {
				if !y(b) {
					return
				}
			}
		}
	}
}

func FieldGrep[B fastats.BedEnter[[]string]](it iter.Seq[B], col int, re *regexp.Regexp) iter.Seq[B] {
	return func(y func(B) bool) {
		for b := range it {
			if len(b.BedFields()) > col && re.MatchString(b.BedFields()[col]) {
				if !y(b) {
					return
				}
			}
		}
	}
}

func ColGrepSingle(r io.Reader, col int, re *regexp.Regexp) io.Reader {
	h := csvh.Handle0("ColGrepSingle: %w")

	rout := PipeWrite(func(w io.Writer) {
		cr := csv.NewReader(r)
		cr.LazyQuotes = true
		cr.ReuseRecord = true
		cr.FieldsPerRecord = -1
		cr.Comma = rune('\t')

		cw := csv.NewWriter(w)
		cw.Comma = rune('\t')
		defer cw.Flush()

		i := 0
		j := 0

		for l, e := cr.Read(); e != io.EOF; l, e = cr.Read() {
			if e != nil {
				panic(h(e))
			}
			if len(l) <= col {
				panic(h(fmt.Errorf("len(l) %v <= col %v", len(l), col)))
			}
			if re.MatchString(l[col]) {
				cw.Write(l)
				j++
			}
			i++
		}
		fmt.Fprintf(os.Stderr, "ColGrep: printed %v of %v lines\n", j, i)
	})
	return rout
}

func RunSingle() {
	f := GetFlags()

	it, errp := iterh.BreakWithError(fastats.ParseBedFlat(os.Stdin))
	err := SinglePlot(it, f.Outpre, f.Chr, f.Start, f.End)
	if *errp != nil {
		panic(*errp)
	}
	if err != nil {
		panic(err)
	}
}
