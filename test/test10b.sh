#!/bin/bash
set -e

(cd .. && (./install.sh))

cat /media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixaf_coverage.txt | \
../cmd/full_single_cov_plot \
	-o test10biso1 \
	-C /home/jgbaldwinbrown/Documents/work_stuff/drosophila/homologous_hybrid_mispairing/refs/combos/ixw/chrlens.txt \
	-c 3L \
	-s 0 \
	-e 1000000 \
> test10bout.txt 2>&1
