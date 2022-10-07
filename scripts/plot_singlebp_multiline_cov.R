#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_cov_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	cov_path = args[1]
	out_path = args[2]
	ymin = as.numeric(args[3])
	ymax = as.numeric(args[4])

	cov = read_bed_cov_named(cov_path, FALSE)

	plot_cov_multi(cov, out_path, 20, 8, 300, calc_chrom_labels_string(cov), ymin, ymax)
}

main()
