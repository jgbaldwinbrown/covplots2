#!/bin/bash
set -e

(cd .. && (./install.sh))

mkdir -p test19out

cat test19cfg.txt | ../cmd/all_singlebp_multiline -w 1000000 -s 1000000 \
> test19out.txt 2>&1
