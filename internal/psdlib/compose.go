package psdlib

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/oov/psd"
)

func containInt(list []int, target int) bool {
	for i := range list {
		if list[i] == target {
			return true
		}
	}
	return false
}

func composeLayer(cImg *image.NRGBA, parentName string, layer *psd.Layer, idList []int) (*image.NRGBA, error) {
	currentName := parentName + "/" + strconv.Itoa(layer.SeqID) + ":" + layer.Name
	cim := cImg
	var err error
	if len(layer.Layer) > 0 {
		for _, l := range layer.Layer {
			if cim, err = composeLayer(cim, currentName, &l, idList); err != nil {
				return cim, err
			}
		}
	}
	if !layer.HasImage() {
		return cim, nil
	}
	if containInt(idList, layer.SeqID) {
		fmt.Printf("%d: %s added\n", layer.SeqID, currentName)
		cim = imaging.Overlay(cim, layer.Picker, layer.Rect.Bounds().Min, 1.0)
	}
	return cim, nil
}

func ReadPsd(psdname string) *psd.PSD {
	file, err := os.Open(psdname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	targetPsd, _, err := psd.Decode(file, &psd.DecodeOptions{SkipMergedImage: true})
	if err != nil {
		panic(err)
	}

	return targetPsd
}

func ComposeImage(targetPsd *psd.PSD, seqList []int) image.Image {
	compositeImage := image.NewNRGBA(targetPsd.Config.Rect)
	var err error
	for _, layer := range targetPsd.Layer {
		if compositeImage, err = composeLayer(compositeImage, "", &layer, seqList); err != nil {
			panic(err)
		}
	}

	return compositeImage
}

func Compose(psdname string, seqList []int, target string) {
	targetPsd := ReadPsd(psdname)
	compositeImage := ComposeImage(targetPsd, seqList)

	out, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, compositeImage)
	if err != nil {
		panic(err)
	}
}
