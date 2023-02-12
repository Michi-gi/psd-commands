package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/oov/psd"
)

func listLayers(parentName string, layer *psd.Layer) {
	currentName := parentName + "/" + strconv.Itoa(layer.SeqID) + ":" + layer.Name
	fmt.Println(currentName)
	if len(layer.Layer) > 0 {
		for _, l := range layer.Layer {
			listLayers(currentName, &l)
		}
	}

	return
}

func LayersList(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := psd.Decode(file, &psd.DecodeOptions{SkipMergedImage: true})
	if err != nil {
		panic(err)
	}

	for _, layer := range img.Layer {
		listLayers("", &layer)
	}
}
