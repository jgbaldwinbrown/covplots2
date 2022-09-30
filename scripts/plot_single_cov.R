#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	cov_path = args[1]
	out_path = args[2]

	cov = read_bed_cov(cov_path, FALSE)
	cov$color = nothreshcolor(cov, "CHR")

	plot_cov(cov, out_path, 20, 8, 300, calc_chrom_labels(cov))
}

main()
