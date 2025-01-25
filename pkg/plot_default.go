package covplots

import (
	"os"
	"os/exec"
	"io"
	"strings"
)

const defaultPlot = `#!/usr/bin/env Rscript
library(data.table)
library(ggplot2)

main = function() {
	args = commandArgs(trailingOnly = TRUE)
	data = as.data.frame(fread(args[1], sep = "\t"))
	colnames(data) = c("Chr", "Start", "End", "Value", "ChrNum", "StartOffset", "EndOffset")
	data$Position = (data$StartOffset + data$EndOffset) / 2
	g = ggplot(data, aes(Position, Value)) +
		geom_point() +
		theme_bw()
	pdf(args[2], height = 3, width = 8)
	print(g)
	dev.off()
}

main()`

func PlotDefault(outpre string) (err error) {
	tmp, e := os.CreateTemp("", "PlotDefault")
	if e != nil {
		return e
	}
	defer func() {
		if e := os.Remove(tmp.Name()); err == nil {
			err = e
		}
	}()
	if _, e := io.Copy(tmp, strings.NewReader(defaultPlot)); e != nil {
		return e
	}
	if e := tmp.Close(); e != nil {
		return e
	}
	cmd := exec.Command("Rscript", tmp.Name(), outpre + "_plfmt.bed", outpre + ".pdf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
