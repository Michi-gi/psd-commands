package main

import (
	"encoding/json"
	"io/ioutil"
)

type ComposeJsonFile struct {
	PsdFilename string `json:"psd_filename"`
	ComposeList []struct {
		Layers         []int  `json:"layers"`
		OutputFilename string `json:"output_filename"`
	} `json:"compose"`
}

func readComposeJson(filename string) (ComposeJsonFile, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return ComposeJsonFile{}, err
	}

	var result ComposeJsonFile
	if err := json.Unmarshal(raw, &result); err != nil {
		return ComposeJsonFile{}, err
	}

	return result, nil
}
