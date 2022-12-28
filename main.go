package main

import (
	"flag"
	"fmt"
	"yoda/output"
	"yoda/src"
)

func main() {
	fileFlag := flag.String("file", "example/example.go", "the file to pass to GetPackageInfo")
	outputFileFlag := flag.String("output", "yoda.json", "the json file in which to output")

	flag.Usage = func() {
		fmt.Println("Usage: yoda [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	packageInfo, err := src.GetPackageInfo(*fileFlag)
	if err != nil {
		panic(err)
	}

	output.CreateJSON(packageInfo, *outputFileFlag)
}
