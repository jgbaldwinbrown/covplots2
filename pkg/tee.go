package covplots

import (
	"os"
	"io"
)

func TeeToPath(r io.Reader, path string) (tr io.Reader, closef func() error, err error) {
	fp, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}
	closef = func() error { return fp.Close() }
	return io.TeeReader(r, fp), closef, nil
}
