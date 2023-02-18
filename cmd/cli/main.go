package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"strconv"

	"internal/psdlib"
)

func main() {
	listFlags := flag.NewFlagSet("list", flag.ExitOnError)
	composeFlags := flag.NewFlagSet("compose", flag.ExitOnError)
	jsonComposeFlags := flag.NewFlagSet("json", flag.ExitOnError)

	fmt.Println((os.Args[1]))

	switch os.Args[1] {
	case "list":
		psdFilename := listFlags.String("p", "", "PSD file")
		listFlags.Parse(os.Args[2:])

		fmt.Printf("list the layer structure of: %s\n\n", *psdFilename)
		LayersList(*psdFilename)

	case "compose":
		psdFilename := composeFlags.String("p", "", "PSD file")
		composeFlags.Parse(os.Args[2:])

		restArgs := composeFlags.Args()
		lastArgPos := len(restArgs) - 1
		if lastArgPos <= 0 {
			fmt.Fprintln(os.Stderr, "not specify layer.")
			os.Exit(1)
		}
		fmt.Println(lastArgPos)
		layerList := make([]int, lastArgPos)
		for i, v := range restArgs[:lastArgPos] {
			layerList[i], _ = strconv.Atoi(v)
		}
		targetFilename := restArgs[lastArgPos]

		fmt.Println("compose picture from PSD")
		fmt.Printf("from: %s\n", *psdFilename)
		fmt.Print("layer: ")
		fmt.Println(layerList)
		fmt.Printf("target: %s\n", targetFilename)
		psdlib.Compose(*psdFilename, layerList, targetFilename)

	case "json":
		jsonFilename := jsonComposeFlags.String("j", "", "JSON file")
		jsonComposeFlags.Parse(os.Args[2:])

		fmt.Printf("compose picture with Json: %s\n", *jsonFilename)
		data, err := readComposeJson(*jsonFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail unmarshalling json: %s\n", *jsonFilename)
			os.Exit(1)
		}

		targetPsd := psdlib.ReadPsd(data.PsdFilename)
		for _, composeData := range data.ComposeList {
			fmt.Println("compose picture from PSD")
			fmt.Printf("from: %s\n", data.PsdFilename)
			fmt.Print("layer: ")
			fmt.Println(composeData.Layers)
			fmt.Printf("target: %s\n", composeData.OutputFilename)
			compositeImage := psdlib.ComposeImage(targetPsd, composeData.Layers)
			file, err := os.Create(composeData.OutputFilename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fail to create file: %s\n", composeData.OutputFilename)
				os.Exit(1)
			}
			defer file.Close()
			err = png.Encode(file, compositeImage)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fail to encode image into PNG: %s\n", composeData.OutputFilename)
				os.Exit(1)
			}
		}
	}
}
