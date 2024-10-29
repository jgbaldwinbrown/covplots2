package covplots

import (
	"strconv"
	"math"
	"iter"

	"github.com/jgbaldwinbrown/fastats/pkg"
)

func Log10[B fastats.BedEnter[float64]](r iter.Seq[B]) iter.Seq[fastats.BedEntry[float64]] {
	return OneArgArith(r, math.Log10)
}

func Abs[B fastats.BedEnter[float64]](r iter.Seq[B]) iter.Seq[fastats.BedEntry[float64]] {
	return OneArgArith(r, math.Abs)
}

func Add[B fastats.BedEnter[Tuple2[float64, float64]]](r iter.Seq[B]) iter.Seq[fastats.BedEntry[float64]] {
	return TwoArgArith(r, func(x, y float64) float64 { return x + y })
}

func OneArgArith[B fastats.BedEnter[float64]](r iter.Seq[B], f func(float64) float64) iter.Seq[fastats.BedEntry[float64]] {
	return func(y func(fastats.BedEntry[float64]) bool) {
		for b := range r {
			var ent fastats.BedEntry[float64]
			ent.ChrSpan = fastats.ToChrSpan(b)
			ent.Fields = f(b.BedFields())
			if !y(ent) {
				return
			}
		}
	}
}

func AlwaysParseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return math.NaN()
	}
	return val
}

type Tuple2[T, U any] struct {
	V0 T
	V1 U
}

func TwoArgArith[B fastats.BedEnter[Tuple2[float64, float64]]](r iter.Seq[B], f func(float64, float64) float64) iter.Seq[fastats.BedEntry[float64]] {
	return func(y func(fastats.BedEntry[float64]) bool) {
		for b := range r {
			var ent fastats.BedEntry[float64]
			ent.ChrSpan = fastats.ToChrSpan(b)
			ent.Fields = f(b.BedFields().V0, b.BedFields().V1)
			if !y(ent) {
				return
			}
		}
	}
}
