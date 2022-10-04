#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test17out

cat test17cfg.txt | ../cmd/all_singlebp_multiline -w 1000000 -s 1000000 \
> test17out.txt 2>&1
