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

func ToIntSlice(a any) []int {
	var out []int
	as := a.([]any)
	for _, ai := range as {
		out = append(out, int(ai.(float64)))
	}
	return out
}

func ToIntSliceSlice(a any) [][]int {
	var out [][]int
	as := a.([]any)
	for _, ai := range as {
		out = append(out, ToIntSlice(ai))
	}
	return out
}

func ColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	colsandreaders := ToIntSliceSlice(args)
	cols := colsandreaders[0]
	readers := colsandreaders[1]
	// if !ok {
	// 	return nil, fmt.Errorf("wrong argument %v to Columns", args)
	// }
	return GetMultipleColsSome(rs, cols, readers), nil
}

func HicSelfColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	return GetMultipleCols(rs, []int{0,1,2,6}), nil
}

func HicSelfColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	return GetMultipleColsSome(rs, []int{0,1,2,6}, somereaders), nil
}

func HicPairColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicPairColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,5}), nil
}

func HicPairColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicPairColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,5}, somereaders), nil
}

func HicPairPropFpkmColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicPairPropFpkmColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,16}), nil
}

func HicPairPropFpkmColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicPairPropFpkmColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,16}, somereaders), nil
}

func HicPairPropColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicPairPropColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,7}), nil
}

func HicPairPropColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicPairPropColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,7}, somereaders), nil
}

func WindowCovColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("WindowCovColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,3}), nil
}

func WindowCovColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	rs_to_subset := ToIntSlice(args)
	fmt.Printf("WindowCovColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,3}, rs_to_subset), nil
}

func GetMultipleCols(rs []io.Reader, cols []int) []io.Reader {
	out := make([]io.Reader, len(rs))
	for i, r := range rs {
		out[i] = GetCols(r, cols)
	}
	return out
}

func GetMultipleColsSome(rs []io.Reader, cols []int, rs_to_subset []int) []io.Reader {
	out := make([]io.Reader, len(rs))
	for i, r := range rs {
		out[i] = r
	}
	for _, ridx := range rs_to_subset {
		out[ridx] = GetCols(rs[ridx], cols)
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
