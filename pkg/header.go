package covplots

import (
	"bufio"
	"io"
	"fmt"
)

func StripOneHeader(r io.Reader) io.Reader {
	b := bufio.NewScanner(r)
	b.Buffer([]byte{}, 1e12)

	return PipeWrite(func(w io.Writer) {
		if !b.Scan() { return }

		for b.Scan() {
			fmt.Fprintln(w, b.Text())
		}
	})
}
