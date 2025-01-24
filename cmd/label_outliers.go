package main

import (
	"github.com/jgbaldwinbrown/covplots2/pkg"
)

func main() {
	err := covplots.FullLabelOutliers()
	if err != nil { panic(err) }
}
