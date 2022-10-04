#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test15out

cat test15cfg.txt | ../cmd/all_singlebp_multiline -w 100000 -s 10000000 \
> test15out.txt 2>&1
