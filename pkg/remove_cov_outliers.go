package covplots

import (
	"slices"
	"iter"
	"math"
	"os"
	"io"
	"strings"
	"flag"
	"fmt"
	"strconv"

	"github.com/jgbaldwinbrown/csvh"
	"github.com/jgbaldwinbrown/fastats/pkg"
	"github.com/jgbaldwinbrown/iterh"
	"github.com/montanaflynn/stats"
)

func ToBedGraphEntry[B fastats.BedEnter[[]string]](col int) func(B) (fastats.BedEntry[float64], error) {
	return func(b B) (fastats.BedEntry[float64], error) {
		var ent fastats.BedEntry[float64]
		ent.ChrSpan = fastats.ToChrSpan(b)
		fields := b.BedFields()
		if len(fields) <= col {
			return ent, fmt.Errorf("ToBedGraphEntry: col %v < len(fields); fields %v", col, fields)
		}
		var e error
		ent.Fields, e = strconv.ParseFloat(fields[col], 64)
		return ent, e
	}
}

func ToBedGraph[B fastats.BedEnter[[]string]](it iter.Seq[B], f func(B) (fastats.BedEntry[float64], error)) iter.Seq2[fastats.BedEntry[float64], error] {
	return func(y func(fastats.BedEntry[float64], error) bool) {
		for b := range it {
			ent, e := f(b)
			if !y(ent, e) {
				return
			}
		}
	}
}

func ToBedGraphEntryColSum[B fastats.BedEnter[[]string]](cols ...int) func(B) (fastats.BedEntry[float64], error) {
	return func(b B) (fastats.BedEntry[float64], error) {
		sum := 0.0
		for _, col := range cols {
			if len(b.BedFields()) <= col {
				return fastats.BedEntry[float64]{}, fmt.Errorf("GetColSums: b.Fields %v <= col %v", b.BedFields(), col)
			}
			f, err := strconv.ParseFloat(b.BedFields()[col], 64)
			if err != nil {
				return fastats.BedEntry[float64]{}, fmt.Errorf("GetColSums: %w", err)
			}
			sum += f
		}
		return fastats.BedEntry[float64]{ChrSpan: fastats.ToChrSpan(b), Fields: sum}, nil
	}
}

func ToBedGraphEntryMustColSum[B fastats.BedEnter[[]string]](cols ...int) func(B) (fastats.BedEntry[float64], error) {
	return func(b B) (fastats.BedEntry[float64], error) {
		sum := 0.0
		for _, col := range cols {
			if len(b.BedFields()) <= col {
				return fastats.BedEntry[float64]{ChrSpan: fastats.ToChrSpan(b), Fields: math.NaN()}, nil
			}
			f, err := strconv.ParseFloat(b.BedFields()[col], 64)
			if err != nil {
				return fastats.BedEntry[float64]{ChrSpan: fastats.ToChrSpan(b), Fields: math.NaN()}, nil
			}
			sum += f
		}
		return fastats.BedEntry[float64]{ChrSpan: fastats.ToChrSpan(b), Fields: sum}, nil
	}
}

func FilterByStdevs(col int, stdevs float64, entries iter.Seq[fastats.BedEntry[[]string]]) ([]fastats.BedEntry[[]string], error) {
	h := csvh.Handle0("GetBedVals: %w")

	getval := ToBedGraphEntry[fastats.BedEntry[[]string]](col)
	bedGraph := ToBedGraph(entries, getval)
	vals := []float64{}
	for ent, err := range bedGraph {
		if err != nil {
			return nil, h(err)
		}
		vals = append(vals, ent.Fields)
	}

	stdev, err := stats.StandardDeviation(vals)
	if err != nil { return nil, h(err) }

	mean, err := stats.Mean(vals)
	if err != nil { return nil, h(err) }

	lowthresh := mean - (stdev * stdevs)
	hithresh := mean + (stdev * stdevs)

	filt := func(b fastats.BedEntry[[]string]) bool {
		val, err := getval(b)
		if err != nil { return false }
		return (val.Fields >= lowthresh) && (val.Fields <= hithresh)
	}

	filtered := iterh.Filter(entries, filt)
	out := slices.Collect(filtered)
	return out, nil
}

func FilterByColsumStdevs(cols []int, stdevs float64, entries iter.Seq[fastats.BedEntry[[]string]]) ([]fastats.BedEntry[[]string], error) {
	h := csvh.Handle0("GetBedVals: %w")

	getval := ToBedGraphEntryColSum[fastats.BedEntry[[]string]](cols...)
	bedGraph := ToBedGraph(entries, getval)
	vals := []float64{}
	for ent, err := range bedGraph {
		if err != nil {
			return nil, h(err)
		}
		vals = append(vals, ent.Fields)
	}

	stdev, err := stats.StandardDeviation(vals)
	if err != nil { return nil, h(err) }

	mean, err := stats.Mean(vals)
	if err != nil { return nil, h(err) }

	lowthresh := mean - (stdev * stdevs)
	hithresh := mean + (stdev * stdevs)

	filt := func(b fastats.BedEntry[[]string]) bool {
		val, err := getval(b)
		if err != nil { return false }
		return (val.Fields >= lowthresh) && (val.Fields <= hithresh)
	}

	filtered := iterh.Filter(entries, filt)
	out := slices.Collect(filtered)
	return out, nil
}

func WriteBedEntries(w io.Writer, entries ...fastats.BedEntry[[]string]) error {
	for _, entry := range entries {
		fmt.Fprintf(w, "%v\t%v\t%v", entry.Chr, entry.Start, entry.End)
		for _, field := range entry.Fields {
			fmt.Fprintf(w, "\t%v", field)
		}
		fmt.Fprintf(w, "\n")
	}
	return nil
}

func ReadCovFiltFlags() ([]int, float64) {
	colsp := flag.String("c", "", "comma-separated columns to filter by")
	stdevp := flag.String("s", "", "Standard deviations to keep")
	flag.Parse()
	if *colsp == "" {
		panic(fmt.Errorf("missing column specifier -c"))
	}
	stdevs, err := strconv.ParseFloat(*stdevp, 64)
	if err != nil {
		panic(fmt.Errorf("Could not parse stdev %v", *stdevp))
	}

	colstrs := strings.Split(*colsp, ",")
	var cols []int
	for _, colstr := range colstrs {
		col, err := strconv.ParseInt(colstr, 0, 64)
		if err != nil {
			panic(fmt.Errorf("Could not parse col %v", colstr))
		}
		cols = append(cols, int(col))
	}
	return cols, stdevs
}

func FullFilterCov() error {
	cols, stdevs := ReadCovFiltFlags()

	entries_it := fastats.ParseBedFlat(os.Stdin)
	entries, err := iterh.CollectWithError(entries_it)
	if err != nil { return err }

	fentries, err := FilterByColsumStdevs(cols, stdevs, slices.Values(entries))
	if err != nil { return err }

	err = WriteBedEntries(os.Stdout, fentries...)
	if err != nil { return err }

	return nil
}
