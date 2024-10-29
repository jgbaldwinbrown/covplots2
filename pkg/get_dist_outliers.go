package covplots

import (
	"slices"
	"iter"
	"os"
	"strings"
	"flag"
	"fmt"
	"strconv"

	"github.com/jgbaldwinbrown/csvh"
	"github.com/jgbaldwinbrown/iterh"
	"github.com/jgbaldwinbrown/fastats/pkg"
	"github.com/montanaflynn/stats"
)

func GetBedEntryLabels[B any](labeller func(B) string, entries iter.Seq[B]) iter.Seq[string] {
	return func(y func(string) bool) {
		for e := range entries {
			if !y(labeller(e)) {
				return
			}
		}
	}
}

func LabelByDistOutliers(cols, distcols []int, outlierPerc, lowOutlierPerc float64, entries, dist iter.Seq[fastats.BedEntry[[]string]]) (iter.Seq[string], error) {
	h := csvh.Handle0("LbelByDistOutliers: %w")

	getval := ToBedGraphEntryMustColSum[fastats.BedEntry[[]string]](cols...)
	getdistval := ToBedGraphEntryMustColSum[fastats.BedEntry[[]string]](distcols...)

	valsit := ToBedGraph(dist, getdistval)
	valsit2, errp := iterh.BreakWithError(valsit)
	vals := slices.Collect(iterh.Transform(valsit2, func(b fastats.BedEntry[float64]) float64 {
		return b.Fields
	}))

	if *errp != nil { return nil, h(*errp) }

	lowthresh, err := stats.Percentile(vals, lowOutlierPerc)
	if err != nil { return nil, h(err) }

	hithresh, err := stats.Percentile(vals, 100.0 - outlierPerc)
	if err != nil { return nil, h(err) }

	labeller := func(b fastats.BedEntry[[]string]) string {
		val, err := getval(b)
		if err != nil { return "NA" }
		if val.Fields > hithresh { return "high" }
		if val.Fields < lowthresh { return "low" }
		return "mid"
	}

	return GetBedEntryLabels(labeller, entries), nil
}

func AppendLabels(entries []fastats.BedEntry[[]string], labels []string) error {
	if len(labels) != len(entries) {
		return fmt.Errorf("AppendLabels: len(labels) %v != len(entries) %v", len(labels), len(entries))
	}

	for i, _ := range entries {
		ep := &entries[i]
		ep.Fields = append(ep.Fields, labels[i])
	}

	return nil
}

type LabellerArgs struct {
	TestCols []int
	TestPath string
	ThreshPerc float64
	LowThreshPerc float64
	DistCols []int
	DistPath string
}

func SplitCommas(in string) []int {
	colstrs := strings.Split(in, ",")
	var cols []int
	for _, colstr := range colstrs {
		col, err := strconv.ParseInt(colstr, 0, 64)
		if err != nil {
			panic(fmt.Errorf("Could not parse col %v", colstr))
		}
		cols = append(cols, int(col))
	}
	return cols
}

func ReadLabellerFlags() LabellerArgs {
	var a LabellerArgs
	var tcstring string
	var dcstring string
	var percstring string
	var lowpercstring string
	flag.StringVar(&tcstring, "tc", "", "comma-separated columns to filter by")
	flag.StringVar(&a.TestPath, "tp", "", "test input path")
	flag.StringVar(&dcstring, "dc", "", "comma-separated columns to establish distribution")
	flag.StringVar(&a.DistPath, "dp", "", "distribution input path")
	flag.StringVar(&percstring, "p", "", "Percent to keep")
	flag.StringVar(&lowpercstring, "lp", "", "Lower percent to keep (if different from high)")
	flag.Parse()

	if a.TestPath == "" { panic(fmt.Errorf("missing -tp")) }
	if a.DistPath == "" { panic(fmt.Errorf("missing -dp")) }

	a.TestCols = SplitCommas(tcstring)
	a.DistCols = SplitCommas(dcstring)

	var err error
	a.ThreshPerc, err = strconv.ParseFloat(percstring, 64)
	if err != nil {
		panic(fmt.Errorf("Could not parse stdev %v", percstring))
	}

	a.LowThreshPerc = a.ThreshPerc
	if lowpercstring != "" {
		a.LowThreshPerc, err = strconv.ParseFloat(lowpercstring, 64)
		if err != nil {
			panic(fmt.Errorf("Could not parse stdev %v", percstring))
		}
	}

	return a
}

func FullLabelOutliers() error {
	args := ReadLabellerFlags()

	testEntries, err := iterh.CollectWithError(iterh.PathIter(args.TestPath, fastats.ParseBedFlat))
	if err != nil { return err }

	distEntries, err := iterh.CollectWithError(iterh.PathIter(args.DistPath, fastats.ParseBedFlat))
	if err != nil { return err }

	labels, err := LabelByDistOutliers(args.TestCols, args.DistCols, args.ThreshPerc, args.LowThreshPerc, slices.Values(testEntries), slices.Values(distEntries))
	if err != nil { return err }

	AppendLabels(testEntries, slices.Collect(labels))

	err = WriteBedEntries(os.Stdout, testEntries...)
	if err != nil { return err }

	return nil
}
