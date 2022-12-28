package example

import (
	"fmt"
	"yoda/example/texts"
)

var EXAMPLE string = "mama"

const EXAMPLE_1 = "example 1"
const EXAMPLE_2 = "example 2"

type mama struct {
	Name  string
	Value string
}

// This is an example
func example() {
	fmt.Println("Hello world")
	text("example")
}

// This is text function
func text(str string) string {

	for i := 0; i < len(str); i++ {
		fmt.Println(str, i)
	}

	return texts.GetText()
}

// this is text2 function
func text2() string {
	var str string
	if len(texts.GetText()) > 0 {
		str = "yoda"
	}
	return example2(texts.GetText()) + str
}

func text1() {
	example()
}
