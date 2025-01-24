package covplots

import (
	"iter"
	"github.com/jgbaldwinbrown/fastats/pkg"
)

func AppendNames(it iter.Seq[fastats.BedEntry[[]string]], names ...string) iter.Seq[fastats.BedEntry[[]string]] {
	return func(y func(fastats.BedEntry[[]string]) bool) {
		for ent := range it {
			ent.Fields = append(ent.Fields, names...)
			if !y(ent) {
				return
			}
		}
	}
}

