package covplots

import (
	"strconv"
	"bufio"
	"github.com/jgbaldwinbrown/lscan/pkg"
	"io"
	"fmt"
	"math"
)

func Log10(rs []io.Reader, args any) ([]io.Reader, error) {
	return OneArgArithMulti(rs, math.Log10), nil
}

func Abs(rs []io.Reader, args any) ([]io.Reader, error) {
	return OneArgArithMulti(rs, math.Abs), nil
}

func OneArgArithMulti(rs []io.Reader, f func(float64) float64) []io.Reader {
	var out []io.Reader
	for _, r := range rs {
		outr := OneArgArith(r, f)
		out = append(out, outr)
	}
	return out
}

func OneArgArith(r io.Reader, f func(float64) float64) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Printf("running Log10 internal func\n")
		var line []string
		split := lscan.ByByte('\t')

		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		i := 0
		for s.Scan() {
			line = lscan.SplitByFunc(line, s.Text(), split)
			var val float64
			if len(line) < 3 {
				continue
			}
			if len(line) < 4 {
				val = math.NaN()
			} else {
				var err error
				val, err = strconv.ParseFloat(line[3], 64)
				if err != nil {
					val = math.NaN()
				} else {
					val = f(val)
				}
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", line[0], line[1], line[2], val)
			i++
		}
		fmt.Printf("OneArgArith lines: %v\n", i)
	})
}
