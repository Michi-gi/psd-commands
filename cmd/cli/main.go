package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	listFlags := flag.NewFlagSet("list", flag.ExitOnError)
	composeflags := flag.NewFlagSet("compose", flag.ExitOnError)

	fmt.Println((os.Args[1]))

	switch os.Args[1] {
	case "list":
		psdFilename := listFlags.String("p", "", "PSD file")
		listFlags.Parse(os.Args[2:])

		fmt.Printf("list the layer structure of: %s\n\n", *psdFilename)
		LayersList(*psdFilename)

	case "compose":
		psdFilename := composeflags.String("p", "", "PSD file")
		composeflags.Parse(os.Args[2:])

		restArgs := composeflags.Args()
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
	}
}
