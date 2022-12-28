package main

import (
	"yoda/output"
	"yoda/src"
)

func main() {
	packageInfo, err := src.GetPackageInfo("example/example.go")
	if err != nil {
		panic(err)
	}

	output.CreateJSON(packageInfo)
}
