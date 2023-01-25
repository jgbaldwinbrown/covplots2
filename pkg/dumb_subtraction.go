package covplots

import (
	"strings"
	"io"
	"bufio"
	"fmt"
)

func ParseEntry(line string) (Entry, error) {
	var e Entry
	_, err := fmt.Sscanf(line, "%s	%d	%d	%f", &e.Chr, &e.Start, &e.End, &e.Val)
	if err != nil {
		return Entry{}, fmt.Errorf("ParsePosVal: %w", err)
	}
	return e, nil
}

func DumbSubtractInternal(r1, r2 io.Reader) (map[Span]SubVal, error) {
	out := map[Span]SubVal{}
	s1 := bufio.NewScanner(r1)
	s1.Buffer([]byte{}, 1e12)
	for s1.Scan() {
		e, err := ParseEntry(s1.Text())
		if err != nil {
			return nil, fmt.Errorf("DumbSubtractInternal: %w", err)
		}
		out[e.Span] = SubVal{e.Val, false}
	}

	s2 := bufio.NewScanner(r2)
	s2.Buffer([]byte{}, 1e12)
	for s2.Scan() {
		e2, err := ParseEntry(s1.Text())
		if err != nil {
			return nil, fmt.Errorf("DumbSubtractInternal: %w", err)
		}
		if sv1, ok := out[e2.Span]; ok {
			out[e2.Span] = SubVal{sv1.Val - e2.Val, true}
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

