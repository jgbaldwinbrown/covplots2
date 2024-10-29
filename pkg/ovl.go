package covplots

import (
	"io"
)

func HicOvlColumns(rs io.Reader) (io.Reader, error) {
	return GetCols(rs, []int{0,1,2,18}), nil
}

func HicNonovlColumns(rs io.Reader) (io.Reader, error) {
	return GetCols(rs, []int{0,1,2,19}), nil
}

func HicNonovlPropFpkmColumns(rs io.Reader) (io.Reader, error) {
	return GetCols(rs, []int{0,1,2,25}), nil
}

func HicNonovlPropColumns(rs io.Reader) (io.Reader, error) {
	return GetCols(rs, []int{0,1,2,21}), nil
}

func HicOvlPropFpkmColumns(rs io.Reader) (io.Reader, error) {
	return GetCols(rs, []int{0,1,2,24}), nil
}

func HicOvlPropColumns(rs io.Reader) (io.Reader, error) {
	return GetCols(rs, []int{0,1,2,20}), nil
}
