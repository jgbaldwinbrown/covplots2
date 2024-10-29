package covplots

import (
	"bufio"
	"strings"
	"io"
	"fmt"
	"github.com/jgbaldwinbrown/lscan/pkg"
)

func FourColumns(rs io.Reader) io.Reader {
	return GetCols(rs, []int{0,1,2,3})
}

func HicSelfColumns(rs io.Reader) io.Reader {
	return GetCols(rs, []int{0,1,2,6})
}

func HicPairColumns(rs io.Reader) io.Reader {
	return GetCols(rs, []int{0,1,2,5})
}

func HicPairPropFpkmColumns(rs io.Reader) io.Reader {
	return GetCols(rs, []int{0,1,2,16})
}

func HicPairPropColumns(rs io.Reader) io.Reader {
	return GetCols(rs, []int{0,1,2,7})
}

func WindowCovColumns(rs io.Reader) io.Reader {
	return GetCols(rs, []int{0,1,2,3})
}

func GetCols(r io.Reader, cols []int) io.Reader {
	return PipeWrite(func(w io.Writer) {
		var line []string
		var colvals []string
		split := lscan.ByByte('\t')

		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		for s.Scan() {
			line = lscan.SplitByFunc(line, s.Text(), split)
			colvals = colvals[:0]
			for _, col := range cols {
				if len(line) > col {
					colvals = append(colvals, line[col])
				} else {
					colvals = append(colvals, "")
				}
			}
			fmt.Fprintln(w, strings.Join(colvals, "\t"))
		}
	})
}
