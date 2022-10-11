package covplots

import (
	"os"
	"io"
	"sync"
)

const bufsize = 8192

func Tee(r io.Reader, npipes int) []io.ReadCloser {
	rs := make([]io.ReadCloser, npipes)
	ws := make([]io.Writer, npipes)
	for i:=0; i<npipes; i++ {
		rs[i], ws[i] = io.Pipe()
	}

	go func() {
		buf := make([]byte, bufsize)
		n, err := r.Read(buf)
		for n != 0 {

			for _, w := range ws {
				_, errw := w.Write(buf)
				if errw != nil {
					break
				}
			}
			if err != nil {
				break
			}
			n, err = r.Read(buf)
		}
		CloseAny(ws...)
	}()
	return rs
}

func TeeToPath(r io.Reader, path string) (io.Reader, *sync.WaitGroup, error) {
	fp, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	trs := Tee(r, 2)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		io.Copy(fp, trs[1])
		fp.Close()
		trs[1].Close()
		wg.Done()
	}()
	return trs[0], wg, nil
}
