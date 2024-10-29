package covplots

import (
	"iter"

	"github.com/jgbaldwinbrown/fastats/pkg"
)

type UltimateConfig struct {
	Input iter.Seq[fastats.BedEntry[[]string]]
	Chrlens string
	Outpre string
	NoParent bool
	PlotFunc func(outpre string) error
	FullChr bool

	ManualChrs iter.Seq[string]
	UseManualChrs bool

	SelectedWins iter.Seq[fastats.ChrSpan]
	UseSelectedWins bool

	Winsize int
	Winstep int
}
