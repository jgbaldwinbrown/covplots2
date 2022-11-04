#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p cov_vs_pair

cat cov_vs_pair_cfg.json | ../cmd/all_singlebp_multiline -w 1000000000 -s 10000000000 \
> test_cov_vs_pairout.txt 2>&1
