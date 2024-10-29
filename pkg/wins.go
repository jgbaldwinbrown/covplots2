package covplots

import (
	"io"
	"fmt"
	"github.com/jgbaldwinbrown/slide/pkg"
)

func SlidingMean(r io.Reader, winSize, winStep float64) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Printf("running Log10 internal func\n")
		slide.SlidingMeans(r, w, winSize, winStep)
	})
}
