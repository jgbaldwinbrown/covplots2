#!/bin/bash
set -e

(cd .. && (./install.sh))

cat /media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixwf_coverage.bg | \
../scripts/full_single_cov_plot -o testpre -C /home/jgbaldwinbrown/Documents/work_stuff/drosophila/homologous_hybrid_mispairing/refs/combos/ixw/chrlens.txt \
> testout.txt 2>&1
