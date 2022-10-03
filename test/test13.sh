#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test13out

cat test13cfg.txt | ../cmd/all_subtract_single_cov_plot -w 1000000 -s 1000000 \
> test13out.txt 2>&1
