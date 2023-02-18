package main

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"image/png"
	"math"
	"os"

	"github.com/go-audio/wav"
	"github.com/oov/psd"

	"internal/psdlib"
)

func main() {

	jsonFile := os.Args[1]
	cj, err := readComposeJson(jsonFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fail to json file")
		os.Exit(1)
	}

	// 音声ファイル
	file, err := os.Open(cj.AudioFilename)
	if err != nil {
		fmt.Println("Error opening gofile:", err)
		os.Exit(1)
	}
	defer file.Close()

	// 音声波形の読み込み
	decoder := wav.NewDecoder(file)
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		fmt.Println("Error reading audio data:", err)
		os.Exit(1)
	}

	bytesMap := make(map[string][]byte)
	targetPsd := psdlib.ReadPsd(cj.PsdFilename)
	for _, co := range cj.ComposeList {
		bytesMap[co.Id] = makeImageBytes(targetPsd, co.Layers)
		fmt.Println(co.Id)
	}

	// 指定された間隔（ここでは33msec）ごとに音量を取得する
	fps := cj.Output.Fps
	interval := int(decoder.SampleRate / uint32(fps))
	j := 0
	sum := 0

	framePos := 0
	for i := 0; i < buf.NumFrames(); i++ {
		sum += int(math.Abs(float64(buf.Data[i])))
		j++
		if j >= interval {
			volume := sum / interval
			outFilename := fmt.Sprintf("%s/%s%010d.png", cj.Output.Dir, cj.Output.Base, framePos)
			file, _ := os.Create(outFilename)
			defer file.Close()
			id := cj.getComposeId(volume)
			fmt.Printf("%s: %d\n", id, volume)
			file.Write(bytesMap[id])
			j = 0
			sum = 0
			framePos++
		}
	}
}

func makeImageBytes(targetPsd *psd.PSD, layers []int) []byte {
	img := psdlib.ComposeImage(targetPsd, layers)
	buf := bytes.NewBuffer([]byte{})
	png.Encode(buf, img)

	return buf.Bytes()
}
