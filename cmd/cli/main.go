package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
		Compose(*psdFilename, layerList, targetFilename)

	case "json":
		jsonFilename := jsonComposeFlags.String("j", "", "JSON file")
		jsonComposeFlags.Parse(os.Args[2:])

		fmt.Printf("compose picture with Json: %s\n", *jsonFilename)
		data, err := readComposeJson(*jsonFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail unmarshalling json: %s\n", *jsonFilename)
			os.Exit(1)
		}

		for _, composeData := range data.ComposeList {
			fmt.Println("compose picture from PSD")
			fmt.Printf("from: %s\n", data.PsdFilename)
			fmt.Print("layer: ")
			fmt.Println(composeData.Layers)
			fmt.Printf("target: %s\n", composeData.OutputFilename)
			Compose(data.PsdFilename, composeData.Layers, composeData.OutputFilename)
		}
	}
}
