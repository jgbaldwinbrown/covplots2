package covplots

import (
	"testing"
	"strings"
	"fmt"
)

var injson = `[
	{
		"inputsets": [
			{
				"paths":[
					"/media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixwf_coverage.bg"
				],
				"name": "ixwf",
				"functions": ["normalize"]
			},
			{
				"paths":[
					"/media/jgbaldwinbrown/jim_work1/melements/cheetah_output/single_coverage/ixaf_coverage.bg"
				],
				"name": "ixaf",
				"functions": ["normalize"]
			},
			{
				"paths":[
					"ixw_hits_1kb_named_rechr.bed"
				],
				"name": "ixw_hic",
				"functions": ["normalize"]
			},
			{
				"paths":[
					"ixw_hits_1kb_named.txt"
				],
				"name": "ixw_hic_self",
				"functions": ["hic_self_cols", "normalize"]
			},
			{
				"paths":[
					"ixw_hits_1kb_named.txt"
				],
				"name": "ixw_hic_pair",
				"functions": ["hic_pair_cols", "normalize"]
			}
		],
		"chrlens": "chrlens.txt",
		"outpre": "test20out/ixw_hic_out",
		"ylim": [-8.0, 8.0]
	}
]`

func TestReadUltimateConfig(t *testing.T) {
	in := strings.NewReader(injson)
	cfg, err := ReadUltimateConfig(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
