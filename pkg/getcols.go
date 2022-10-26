package covplots

import (
	"bufio"
	"strings"
	"io"
	"fmt"
	"github.com/jgbaldwinbrown/lscan/pkg"
)

func Columns(rs []io.Reader, args any) ([]io.Reader, error) {
	cols, ok := args.([]int)
	if !ok {
		return nil, fmt.Errorf("wrong argument %v to Columns", args)
	}
	return GetMultipleCols(rs, cols), nil
}

func HicSelfColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	return GetMultipleCols(rs, []int{0,1,2,6}), nil
}

func HicPairColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicPairColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,5}), nil
}

func WindowCovColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("WindowCovColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,3}), nil
}

func GetMultipleCols(rs []io.Reader, cols []int) []io.Reader {
	out := make([]io.Reader, len(rs))
	for i, r := range rs {
		out[i] = GetCols(r, cols)
	}
	return out
}

func GetCols(r io.Reader, cols []int) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Println("running GetCols internal func")
		var line []string
		var colvals []string
		split := lscan.ByByte('\t')

		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		for s.Scan() {
			// fmt.Printf("scanning line: %v\n", s.Text())
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
