#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test12out

cat test12cfg.txt | ../cmd/all_single_cov_plot -w 1000000 -s 1000000 \
> test12out.txt 2>&1
