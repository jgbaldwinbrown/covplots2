package covplots

import (
	"strings"
	"io"
	"bufio"
	"fmt"

	"github.com/jgbaldwinbrown/fastats/pkg"
)

func ParseEntry(line string) (fastats.BedEntry[float64], error) {
	var e fastats.BedEntry[float64]
	// fmt.Fprintf(os.Stderr, "line: |%v|\n", line)
	_, err := fmt.Sscanf(line, "%s	%d	%d	%f", &e.Chr, &e.Start, &e.End, &e.Fields)
	if err != nil {
		return fastats.BedEntry[float64]{}, fmt.Errorf("ParsePosVal: %w; line: |%v|", err, line)
	}
	return e, nil
}

func DumbSubtractInternal(r1, r2 io.Reader) (map[fastats.ChrSpan]SubVal, error) {
	out := map[fastats.ChrSpan]SubVal{}
	s1 := bufio.NewScanner(r1)
	s1.Buffer([]byte{}, 1e12)
	for s1.Scan() {
		if s1.Err() != nil {
			return nil, fmt.Errorf("DumbSubtractInternal: s1 error: %w", s1.Err())
		}
		if s1.Text() == "" {
			continue
		}
		e, err := ParseEntry(s1.Text())
		if err != nil {
			return nil, fmt.Errorf("DumbSubtractInternal: %w", err)
		}
		out[e.ChrSpan] = SubVal{e.Fields, false}
	}

	s2 := bufio.NewScanner(r2)
	s2.Buffer([]byte{}, 1e12)
	for s2.Scan() {
		if s2.Err() != nil {
			return nil, fmt.Errorf("DumbSubtractInternal: s2 error: %w", s2.Err())
		}
		if s2.Text() == "" {
			continue
		}
		e2, err := ParseEntry(s2.Text())
		if err != nil {
			return nil, fmt.Errorf("DumbSubtractInternal: %w", err)
		}
		if sv1, ok := out[e2.ChrSpan]; ok {
			out[e2.ChrSpan] = SubVal{sv1.Val - e2.Fields, true}
		}
	}
	return out, nil
}

func DumbSubtract(r1, r2 io.Reader) (*strings.Reader, error) {
	sub, err := DumbSubtractInternal(r1, r2)
	if err != nil {
		return nil, fmt.Errorf("Subtract: %w", err)
	}

	var out strings.Builder
	for entry, sval := range sub {
		if sval.Subtracted {
			fmt.Fprintf(&out, "%s\t%d\t%d\t%f\n", entry.Chr, entry.Start, entry.End, sval.Val)
		}
	}
	return strings.NewReader(out.String()), nil
}

func DumbSubtractTwo(rs []io.Reader, args any) ([]io.Reader, error) {
	newreader, err := DumbSubtract(rs[0], rs[1])
	if err != nil {
		return nil, fmt.Errorf("DumbSubtractTwo: %w", err)
	}
	return []io.Reader{newreader}, nil
}

