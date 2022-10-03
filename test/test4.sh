#!/bin/bash
set -e

(cd .. && (./install.sh))

cat /media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixwf_coverage.bg | \
grep '3L_ISO1' | \
awk -F "\t" -v OFS="\t" '$2 < 1000000' \
> 3l_iso1_mini.bg

cat 3l_iso1_mini.bg | ../scripts/full_single_cov_plot -o test4iso1 -C /home/jgbaldwinbrown/Documents/work_stuff/drosophila/homologous_hybrid_mispairing/refs/combos/ixw/chrlens.txt \
> test4out.txt 2>&1

cat /media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixwf_coverage.bg | \
grep '3L_W501' | \
awk -F "\t" -v OFS="\t" '$2 < 1000000' \
> 3l_w501_mini.bg

cat 3l_w501_mini.bg | ../scripts/full_single_cov_plot -o test4w501 -C /home/jgbaldwinbrown/Documents/work_stuff/drosophila/homologous_hybrid_mispairing/refs/combos/ixw/chrlens.txt \
>> test4out.txt 2>&1
