#!/bin/bash
set -e

cat ixw_hits_1kb_named.bed | \
awk -F "\t" -v OFS="\t" '{a=$1; $1=sprintf("%s_ISO1", a); print $0; $1=sprintf("%s_W501", a); print $0}' \
> ixw_hits_1kb_named_rechr.bed


cat ixw_hits_1kb_named.txt | \
awk -F "\t" -v OFS="\t" '{a=$1; $1=sprintf("%s_ISO1", a); print $0; $1=sprintf("%s_W501", a); print $0}' \
> ixw_hits_1kb_named_rechr.txt
