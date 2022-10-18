#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test22out

cat test22cfg.txt | ../cmd/all_singlebp_multiline -w 100000 -s 10000000 \
> test22out.txt 2>&1
