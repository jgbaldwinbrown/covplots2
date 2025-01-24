package main

import (
	"github.com/jgbaldwinbrown/covplots2/pkg"
)

func main() {
	err := covplots.FullFilterCov()
	if err != nil { panic(err) }
}
