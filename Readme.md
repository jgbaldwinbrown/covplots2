# Covplots

A framework for quick plotting of data from multiple genomes and chromosomes

## Introduction

Covplots is a go pipeline that is designed to flexibly take in multiple
datasets in different formats and plot them as a Manhattan plot. You specify
your input files and the transformations to do on the data using a .json-formatted
config file. Here's an example:

```json
[
	{
		"inputsets": [
			{
				"paths":[ "coverage_bedgraph.bed" ],
				"name": "hxwf",
				"functions": ["cov_win_cols", "per_bp", "normalize"]
			},
			{
				"paths":[ "100kb_sliding_windows_self.txt" ],
				"name": "hxw_hic_self",
				"functions": ["hic_self_cols", "normalize"]
			},
			{
				"paths":[ "100kb_sliding_windows_paired.txt" ],
				"name": "hxw_hic_pair",
				"functions": ["hic_pair_cols", "normalize"]
			},
			{
				"paths":[ "100kb_sliding_window_pair_proportion.txt" ],
				"name": "hxw_hic_pair_prop",
				"functions": ["hic_pair_prop_cols", "normalize"]
			}
		],
		"chrlens": "chrlens.txt",
		"outpre": "outdir/plots",
		"ylim": [-8.0, 8.0]
	}

```

To run:

```sh
mkdir -p outdir
cat cfg.json | all_singlebp_multiline -w 1000000 -s 100000 \
```

This code will produce one set of plots with four lines each. These lines
correspond to the "inputsets" portion of the config file. The program will
produce plots in sliding windows. Each plot will cover 1Mb of sequence (the -w
option), and they will be staggered by 100kb (the -s option) such that they
overlap. There are two different input formats -- one for the first input set,
and another for the other three input sets. This is possible because different
functions are used on each input set.

The first input set is a multi-column tab-delimited bed file that contains
counts of sequencing reads in spans specified by the first three columns of the
bed file (chrom, start, end, half-open 0-based intervals, as is standard for
bed files). This is the core format used by the program. You'll see that it
uses three functions: "cov_win_cols", "per_bp", and "normalize". These
functions are run in order, and transform the data in a stream. "cov_win_cols"
extracts the appropriate columns for counting coverage from the output of
bedtools' coverage tool. It produces data in this format:

```
chrom	start	end	val
```

The next tool, per_bp, divides the value by the length of the span (end -
start), normalizing the coverage per basepair. The final function, "normalize",
normalizes the data by subtracting the mean of the data, then dividing by the
standard deviation of the data.

Note that this ends in a four-column .bed-format file. Any set of defined
functions can be used, but the last function must leave the data in this
four-column format. This is the format used for plotting.

Here are all of the currently available functions:

- subtract_two
	- operates on exactly two 4-column bed files, subtracting the values in the 2nd one from the values in the 1st one
- unchanged
	- does nothing -- a placeholder
- normalize
	- Works on any number 4-column bed files. Subtracts the mean of the value column and divides by the standard deviation.
- columns
	- Works on any number of tab-separated files. Extracts the specified 0-indexed columns (using the "functionargs" variable).
	- example:
```json
{
	...
	"functions": ["columns"],
	"functionargs": [[3,4,5]],
	...
}
```
			
- columns_some
	- Works on any number of tab-separated files. Extracts the specified 0-indexed columns (using the "functionargs" variable). Only operates on the files indexed by the 2nd input argument list. All other files are unchanged
	- example:
```json
{
	...
	"functions": ["columns"],
	"functionargs": [[[3,4,5], [0,1]]],
	...
}
```
- hic_self_cols
	- Extracts the self-interacting read counts from a pairviz output file into a 4-column bed file. Works on any number of files. The alternative function "hic_self_cols" takes a list of indices as an argument, and only changes the indexed files.
- hic_pair_cols
	- Extracts the pair-interacting read counts from a pairviz output file into a 4-column bed file. Works on any number of files. The alternative function "hic_pair_cols" takes a list of indices as an argument, and only changes the indexed files.
- hic_pair_prop_cols
	- Extracts the pair_prop-interacting read counts from a pairviz output file into a 4-column bed file. Works on any number of files. The alternative function "hic_pair_prop_cols" takes a list of indices as an argument, and only changes the indexed files.
- cov_win_cols
	- Extracts the windowed coverage counts from a bedtools coverage output file into a 4-column bed file. Works on any number of files. The alternative function "cov_win_cols" takes a list of indices as an argument, and only changes the indexed files.
- per_bp
	- Takes any number of 4-column bed files. Divides the value column by the length of the span (end - start).
- combine_to_one_line
	- Takes any number of 4-column bed files. Combines them on a per-bp basis so that each column after the first three represents the value from a file, i.e.:

starting file 1:

```
chr	start	end	val1
```

starting file 2:

```
chr	start	end	val2
```

final file:

```
chr	start	end	val1	val2
```
