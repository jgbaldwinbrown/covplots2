package covplots

import (
	"encoding/csv"
	"regexp"
	"io"
	"os"
	"fmt"

	"github.com/jgbaldwinbrown/csvh"
)

func ColSedSingle(r io.Reader, col int, re *regexp.Regexp, replace string) (io.Reader) {
	h := csvh.Handle0("ColSedSingle: %w")

	rout := PipeWrite(func(w io.Writer) {
		cr := csv.NewReader(r)
		cr.LazyQuotes = true
		cr.ReuseRecord = true
		cr.FieldsPerRecord = -1
		cr.Comma = rune('\t')

		cw := csv.NewWriter(w)
		cw.Comma = rune('\t')
		defer cw.Flush()

		i := 0
		j := 0

		for l, e := cr.Read() ; e != io.EOF; l, e = cr.Read() {
			if e != nil {
				panic(h(e))
			}
			if len(l) <= col {
				panic(h(fmt.Errorf("len(l) %v <= col %v", len(l), col)))
			}
			l[col] = re.ReplaceAllString(l[col], replace)
			cw.Write(l)
			j++
			i++
		}
		fmt.Fprintf(os.Stderr, "ColSed: printed %v of %v lines\n", j, i)
	})
	return rout
}

