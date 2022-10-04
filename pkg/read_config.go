package covplots

import (
	"io"
	"encoding/json"
)

type InputSet struct {
	Paths []string `json:"paths"`
	Names []string `json:"names"`
	Function string `json:"function"`
	FunctionArgs any `json: "functionargs"`
	Extra any `json: "extra"`
}

type UltimateConfig struct {
	InputSets []InputSet
	Chrlens string `json:"chrlens"`
	Outpre string `json:"outpre"`
}

func ReadUltimateConfig(r io.Reader) ([]UltimateConfig, error) {
	cfgbytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var cfg []UltimateConfig
	err = json.Unmarshal(cfgbytes, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetUltimateConfig(path string) ([]UltimateConfig, error) {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	return ReadUltimateConfig(r)
}
