#!/bin/bash
set -e

(cd cmd && (
	ls *.go | while read i ; do
		go build $i
	done
))

(cd scripts && (
	ls *.go | while read i ; do
		go build $i
	done
))

cp scripts/full_single_cov_plot ~/mybin/full_single_cov_plot
chmod +x ~/mybin/full_single_cov_plot
cp scripts/plot_single_cov.R ~/mybin/plot_single_cov
chmod +x ~/mybin/plot_single_cov
cp scripts/plot_cov_helpers.R ~/rlibs
