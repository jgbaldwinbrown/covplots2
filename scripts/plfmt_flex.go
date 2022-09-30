package main

import (
	"os"
	"flag"
	"fmt"
	"strconv"
	"io"
	"bufio"
	"strings"
	"sort"
)

const (
	NAMED int = iota
	INORDER
	NUMERIC
)

type Flags struct {
	Header bool
	ChrCol int
	BpCol int
	BpCol2 int
	ChrBedPath string
	Named bool
}

type Entry struct {
	Line []string
	ChrStr string
	Chr int
	Bp int
	Bp2 int
}

type Entries []Entry
func (e Entries) Len() int { return len(e) }
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Entries) Less(i, j int) bool {
	if e[i].Chr == e[j].Chr {
		return e[i].Bp < e[j].Bp
	}
	return e[i].Chr < e[j].Chr
}

func GetEntry(line []string, chrcol, bpcol, bpcol2 int, chrflag int, chrname map[string]int) (e Entry, err error) {
	e.Line = line
	e.ChrStr = line[0]

	if chrflag == NUMERIC {
		var chrnum float64
		chrnum, err = strconv.ParseFloat(line[chrcol][3:], 64)
		if err != nil { return }
		e.Chr = int(chrnum)
	}
	if chrflag == NAMED {
		var ok bool
		e.Chr, ok = chrname[line[chrcol]]
		if !ok { return }
	}

	bp, err := strconv.ParseFloat(line[bpcol], 64)
	e.Bp = int(bp)

	e.Bp2 = -1
	if bpcol2 != -1 {
		bp, err = strconv.ParseFloat(line[bpcol2], 64)
		e.Bp2 = int(bp)
	}
	return
}

func GetData(r io.Reader, chrcol int, bpcol int, bpcol2 int, header bool, chrflag int, chrlens Entries) (es Entries, err error, header_string string) {
	s := bufio.NewScanner(r)

	if header {
		s.Scan()
		header_string = s.Text()
	}

	chrnames := make(map[string]int)
	if chrflag == NAMED {
		for _, e := range chrlens {
			chrnames[e.ChrStr] = e.Chr
		}
	}

	for i:=1; s.Scan(); i++ {
		line := strings.Split(s.Text(), "\t")
		var e Entry
		e, err = GetEntry(line, chrcol, bpcol, bpcol2, chrflag, chrnames)
		if err != nil { return }
		if chrflag == INORDER {
			e.Chr = i
		}
		es = append(es, e)
	}
	return
}

// func Flush(buf *Buf) {
// 	lines := strings.Join(buf.Lines, "\n")
// 	fmt.Fprintf(buf.Writer, "%s\n", lines)
// 	buf.Lines = buf.Lines[:0]
// }

func FprintEntry(w io.Writer, e Entry, cumsum, cumsum2 int) {
	fmt.Fprintf(w, "%s\t%v\t%v", strings.Join(e.Line, "\t"), e.Chr, cumsum)
	if (e.Bp2 != -1) {
		fmt.Fprintf(w, "\t%v", cumsum2)
	}
	fmt.Fprintf(w, "\n")
	// line := fmt.Sprintf("%s\t%v\t%v\n", strings.Join(e.Line, "\t"), e.Chr, cumsum)
	// buf.Lines = append(buf.Lines, line)
	// if len(buf.Lines) > buf.Max {
	// 	Flush(buf)
	// }
}

// type Buf struct {
// 	Writer io.Writer
// 	Lines []string
// 	Max int
// }

func append_pls(data Entries, w io.Writer) {
	sort.Sort(data)
	cumsum := 0
	cumsum2 := 0
	prevchr := -1
	prevbp := 0
	prevbp2 := 0

	// outbuf := Buf{}
	// outbuf.Max = 10000
	// outbuf.Writer = w
	// defer Flush(&outbuf)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for _, e := range data {
		if e.Chr != prevchr {
			// cumsum += 1000
			// cumsum2 += 1000
			prevbp = 0
			prevbp2 = 0
		}
		cumsum += e.Bp - prevbp
		if e.Bp2 != -1 {
			cumsum2 += e.Bp2 - prevbp2
		}
		FprintEntry(bw, e, cumsum, cumsum2)
		prevchr = e.Chr
		prevbp = e.Bp
		prevbp2 = e.Bp2
	}
}

func ReadChrBed(path string, chrflag int) (Entries, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	es, err, _ := GetData(r, 0, 1, 2, false, chrflag, nil)
	return es, err
}


func GetOffsets(es Entries) map[string]int {
	var sorted Entries
	for _, e := range es {
		e.Line = append([]string{}, e.Line...)
		sorted = append(sorted, e)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Chr < sorted[j].Chr })
	offs := make(map[string]int, len(sorted))
	cumsum := 0
	for _, e := range sorted {
		offs[e.ChrStr] = cumsum
		// cumsum += 1000 + e.Bp2 - e.Bp
		cumsum += e.Bp2 - e.Bp
	}
	return offs
}

func appendPlsChrBed(data Entries, offsets map[string]int, w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for _, e := range data {
		cumsum := e.Bp + offsets[e.ChrStr]

		cumsum2 := -1
		if e.Bp2 != -1 {
			cumsum2 = e.Bp2 + offsets[e.ChrStr]
		}

		FprintEntry(bw, e, cumsum, cumsum2)
	}
}

func Plfmt(flags Flags, r io.Reader, w io.Writer) {
	chrflag := NUMERIC
	if flags.Named {
		chrflag = INORDER
	}
	chrlens, err := ReadChrBed(flags.ChrBedPath, chrflag)
	if err != nil { panic(err) }

	chrflag = NUMERIC
	if flags.Named {
		chrflag = NAMED
	}
	data, err, _ := GetData(r, flags.ChrCol, flags.BpCol, flags.BpCol2, flags.Header, chrflag, chrlens)
	if err != nil { panic(err) }

	chroffsets := GetOffsets(chrlens)
	appendPlsChrBed(data, chroffsets, w)
}

func GetFlags() (f Flags) {
	flag.BoolVar(&f.Header, "H", false, "File includes a header line.")
	flag.IntVar(&f.ChrCol, "c", -1, "0-indexed column containing chromosome in format \"chr[0-9]*\"")
	flag.IntVar(&f.BpCol, "b", -1, "Column containing basepair position.")
	flag.IntVar(&f.BpCol2, "b2", -1, "Column containing end coordinates of spans (optional).")
	flag.StringVar(&f.ChrBedPath, "C", "", "Path to .bed file containing all chromosomes' lengths.")
	flag.BoolVar(&f.Named, "n", false, "Instead of having numbered chromosomes in format \"chr[0-9]*\", use the order of chromosomes in the -C bed file.")
	flag.Parse()
	if f.ChrCol == -1 || f.BpCol == -1 || f.ChrBedPath == "" {
		panic(fmt.Errorf("Missing argument"))
	}
	return
}

func main() {
	flags := GetFlags()
	Plfmt(flags, os.Stdin, os.Stdout)
}
