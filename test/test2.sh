#!/bin/bash
set -e

(cd .. && (./install.sh))

plot_single_cov testpre_plfmt.bed testpre_plot.png
