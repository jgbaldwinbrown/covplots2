#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test18out

cat test18cfg.txt | ../cmd/all_singlebp_multiline -w 100000 -s 10000000 \
> test18out.txt 2>&1
