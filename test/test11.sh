#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test11out

cat test11cfg.txt | ../cmd/all_single_cov_plot -w 1000 -s 10000000 \
> test11out.txt 2>&1
