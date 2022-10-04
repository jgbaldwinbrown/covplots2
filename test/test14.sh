#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test14out

cat test14cfg.txt | ../cmd/all_subtract_single_cov_plot -w 100000 -s 10000000 \
> test14out.txt 2>&1
