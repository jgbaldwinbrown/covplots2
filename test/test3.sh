#!/bin/bash
set -e

(cd .. && (./install.sh))

cat /media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixwf_coverage.bg | \
grep '3L' > 3l_only.bg

cat 3l_only.bg | ../scripts/full_single_cov_plot -o test3 -C /home/jgbaldwinbrown/Documents/work_stuff/drosophila/homologous_hybrid_mispairing/refs/combos/ixw/chrlens.txt \
> test3out.txt 2>&1
