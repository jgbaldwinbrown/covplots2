#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_pretty_multiple_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	full_pfst_path = args[1]
	full_pfst_sig_path = args[2]

	black_pfst_path = args[3]
	black_pfst_sig_path = args[4]

	white_pfst_path = args[5]
	white_pfst_sig_path = args[6]

	runt_pfst_path = args[7]
	runt_pfst_sig_path = args[8]

	figurita_pfst_path = args[9]
	figurita_pfst_sig_path = args[10]

	out_path = args[11]

	full_pfst = read_pvals_nowin(full_pfst_path)
	black_pfst = read_pvals_nowin(black_pfst_path)
	white_pfst = read_pvals_nowin(white_pfst_path)
	runt_pfst = read_pvals_nowin(runt_pfst_path)
	figurita_pfst = read_pvals_nowin(figurita_pfst_path)

	full_pfst_hithresh = calc_thresh(full_pfst, "VAL", .999, TRUE)
	full_pfst_lowthresh = -log10(0.05)
	full_pfst$pass = pass_thresh(full_pfst, "VAL", full_pfst_hithresh)
	full_pfst$color = threshcolor(full_pfst, "CHR", "pass")
	full_pfst_threshes = data.frame(THRESH=c(full_pfst_hithresh, full_pfst_lowthresh))

	black_pfst_hithresh = calc_thresh(black_pfst, "VAL", .999, TRUE)
	black_pfst_lowthresh = -log10(0.05)
	black_pfst$pass = pass_thresh(black_pfst, "VAL", black_pfst_hithresh)
	black_pfst$color = threshcolor(black_pfst, "CHR", "pass")
	black_pfst_threshes = data.frame(THRESH=c(black_pfst_hithresh, black_pfst_lowthresh))

	white_pfst_hithresh = calc_thresh(white_pfst, "VAL", .999, TRUE)
	white_pfst_lowthresh = -log10(0.05)
	white_pfst$pass = pass_thresh(white_pfst, "VAL", white_pfst_hithresh)
	white_pfst$color = threshcolor(white_pfst, "CHR", "pass")
	white_pfst_threshes = data.frame(THRESH=c(white_pfst_hithresh, white_pfst_lowthresh))

	runt_pfst_hithresh = calc_thresh(runt_pfst, "VAL", .999, TRUE)
	runt_pfst_lowthresh = -log10(0.05)
	runt_pfst$pass = pass_thresh(runt_pfst, "VAL", runt_pfst_hithresh)
	runt_pfst$color = threshcolor(runt_pfst, "CHR", "pass")
	runt_pfst_threshes = data.frame(THRESH=c(runt_pfst_hithresh, runt_pfst_lowthresh))

	figurita_pfst_hithresh = calc_thresh(figurita_pfst, "VAL", .999, TRUE)
	figurita_pfst_lowthresh = -log10(0.05)
	figurita_pfst$pass = pass_thresh(figurita_pfst, "VAL", figurita_pfst_hithresh)
	figurita_pfst$color = threshcolor(figurita_pfst, "CHR", "pass")
	figurita_pfst_threshes = data.frame(THRESH=c(figurita_pfst_hithresh, figurita_pfst_lowthresh))

	joinlist = join(
		list(full_pfst, black_pfst, white_pfst, runt_pfst, figurita_pfst),
		list(full_pfst_threshes, black_pfst_threshes, white_pfst_threshes, runt_pfst_threshes, figurita_pfst_threshes),
		c("All", "R1", "R2", "R3", "R4"),
		c("All", "R1", "R2", "R3", "R4")
	)
	data = joinlist[[1]]
	thresholds = joinlist[[2]]

	full_pfst_sig_rect = bed2rect(full_pfst_sig_path)
	if (nrow(full_pfst_sig_rect) > 0) {
		full_pfst_sig_rect$NAME = "All"
	}

	black_pfst_sig_rect = bed2rect(black_pfst_sig_path)
	if (nrow(black_pfst_sig_rect) > 0) {
		black_pfst_sig_rect$NAME = "R1"
	}
	white_pfst_sig_rect = bed2rect(white_pfst_sig_path)
	if (nrow(white_pfst_sig_rect) > 0) {
		white_pfst_sig_rect$NAME = "R2"
	}
	runt_pfst_sig_rect = bed2rect(runt_pfst_sig_path)
	if (nrow(runt_pfst_sig_rect) > 0) {
		runt_pfst_sig_rect$NAME = "R3"
	}
	figurita_pfst_sig_rect = bed2rect(figurita_pfst_sig_path)
	if (nrow(figurita_pfst_sig_rect) > 0) {
		figurita_pfst_sig_rect$NAME = "R4"
	}

	significant_boxes = as.data.frame(rbind(
		full_pfst_sig_rect,
		black_pfst_sig_rect,
		white_pfst_sig_rect,
		runt_pfst_sig_rect,
		figurita_pfst_sig_rect
	))

	scales_y = list (
		`All` = scale_y_continuous(limits = c(0, 30)),
		`R1` = scale_y_continuous(limits = c(0, 30)),
		`R2` = scale_y_continuous(limits = c(0, 30)),
		`R3` = scale_y_continuous(limits = c(0, 30)),
		`R4` = scale_y_continuous(limits = c(0, 30))
	)

	plot_scaled_y_boxed(data, VAL, out_path, 20, 8, 300, thresholds, calc_chrom_labels(black_pfst), scales_y, significant_boxes)
}

main()
