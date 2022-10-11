package covplots

import (
	"io"
)

func PipeWrite(f func(io.Writer)) io.ReadCloser {
	r, w := io.Pipe()
	go func() {
		f(w)
		w.Close()
	}()
	return r
}
