package main

import (
	"encoding/json"
	"io/ioutil"
)

type ComposeJsonFile struct {
	PsdFilename   string `json:"psd_filename"`
	AudioFilename string `json:"audio_filename"`
	ComposeList   []struct {
		Id        string `json:"id"`
		Threshold int    `json:"threshold"`
		Layers    []int  `json:"layers"`
	} `json:"compose"`
	Output struct {
		Dir  string `json:"outdir"`
		Base string `json:"base"`
		Fps  int    `json:"fps"`
	} `json:"output"`
}

func (c ComposeJsonFile) getComposeId(threshold int) string {
	maxThreshold := -1
	var curId string
	for _, co := range c.ComposeList {
		if co.Threshold <= threshold && maxThreshold < co.Threshold {
			curId = co.Id
			maxThreshold = co.Threshold
		}
	}
	return curId
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
