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

cp scripts/find_caph2_peaks ~/mybin/find_caph2_peaks
chmod +x ~/mybin/find_caph2_peaks
cp scripts/full_single_cov_plot ~/mybin/full_single_cov_plot
chmod +x ~/mybin/full_single_cov_plot
cp cmd/all_singlebp_multiline ~/mybin/all_singlebp_multiline
chmod +x ~/mybin/all_singlebp_multiline
cp cmd/filter_cov_outliers ~/mybin/filter_cov_outliers
chmod +x ~/mybin/filter_cov_outliers
cp cmd/label_outliers ~/mybin/label_outliers
chmod +x ~/mybin/label_outliers
cp scripts/plot_single_cov.R ~/mybin/plot_single_cov
chmod +x ~/mybin/plot_single_cov
cp scripts/plot_sub_single_cov.R ~/mybin/plot_sub_single_cov
chmod +x ~/mybin/plot_sub_single_cov

cp scripts/plot_singlebp_multiline_cov.R ~/mybin/plot_singlebp_multiline_cov
chmod +x ~/mybin/plot_singlebp_multiline_cov

cp scripts/plot_singlebp_multiline_cov_pretty.R ~/mybin/plot_singlebp_multiline_cov_pretty
chmod +x ~/mybin/plot_singlebp_multiline_cov_pretty

cp scripts/plot_singlebp_multiline_cov_facet.R ~/mybin/plot_singlebp_multiline_cov_facet
chmod +x ~/mybin/plot_singlebp_multiline_cov_facet

cp scripts/plot_singlebp_multiline_cov_facetscales.R ~/mybin/plot_singlebp_multiline_cov_facetscales
chmod +x ~/mybin/plot_singlebp_multiline_cov_facetscales

cp scripts/plot_cov_helpers.R ~/rlibs

cp scripts/plot_cov_vs_pair.R ~/mybin/plot_cov_vs_pair
chmod +x ~/mybin/plot_cov_vs_pair

cp scripts/plot_self_vs_pair.R ~/mybin/plot_self_vs_pair
chmod +x ~/mybin/plot_self_vs_pair

cp scripts/plot_self_vs_pair_lim.R ~/mybin/plot_self_vs_pair_lim
chmod +x ~/mybin/plot_self_vs_pair_lim

cp scripts/plot_self_vs_pair_pretty.R ~/mybin/plot_self_vs_pair_pretty
chmod +x ~/mybin/plot_self_vs_pair_pretty

cp scripts/plot_cov_vs_pair_minimal.R ~/mybin/plot_cov_vs_pair_minimal
chmod +x ~/mybin/plot_cov_vs_pair_minimal

cp scripts/plot_caph2_pairing.R ~/mybin/plot_caph2_pairing
chmod +x ~/mybin/plot_caph2_pairing

cp scripts/plot_repeat_structdiff.R ~/mybin/plot_repeat_structdiff
chmod +x ~/mybin/plot_repeat_structdiff

cp scripts/plot_repeat_pair.R ~/mybin/plot_repeat_pair
chmod +x ~/mybin/plot_repeat_pair

cp scripts/plot_cov_hist.R ~/mybin/plot_cov_hist
chmod +x ~/mybin/plot_cov_hist

cp scripts/plot_boxwhisker.R ~/mybin/plot_boxwhisker
chmod +x ~/mybin/plot_boxwhisker

cp scripts/plot_singlebp_multiline_cov_facetname_scales.R ~/mybin/plot_singlebp_multiline_cov_facetname_scales
chmod +x ~/mybin/plot_singlebp_multiline_cov_facetname_scales
