# Yoda

A Go/WASI tool to get a quick summary of `.go` files.

I something that helped me wrap my head around large Golang codebases, therefore I started building this tool.

## Usage

- Clone a repository, or create a Go application. 
- Install what you have to

```bash
yoda <FILE_NAME>
```

## Output sample (based on `example/example.go`)

Given this input

```go
package example

import (
	"fmt"
	"yoda/example/texts"
)

var EXAMPLE string = "mama"

const EXAMPLE_1 = "example 1"
const EXAMPLE_2 = "example 2"

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
	return example2(texts.GetText())
}

```
Yoda will generate the following JSON
(Json output filters unnecessary output)

```json
{
    "Name": "example",
    "Imports": {
        "fmt": {
            "Name": "fmt",
            "Path": "fmt",
            "Size": 212331
        },
        "texts": {
            "Name": "texts",
            "Path": "yoda/example/texts",
            "Size": 63
        }
    },
    "Functions": {
        "example": {
            "Name": "example",
            "Doc": "This is an example\n",
            "Calls": {
                "Println": {
                    "Name": "Println",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 15
                    },
                    "ImportPath": "fmt"
                },
                "text": {
                    "Name": "text",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 16
                    }
                }
            },
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 171,
            "Where": {
                "File": "example/example.go",
                "Line": 14
            }
        },
        "text": {
            "Name": "text",
            "Doc": "This is text function\n",
            "Calls": {
                "GetText": {
                    "Name": "GetText",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 26
                    },
                    "ImportPath": "texts"
                },
                "Println": {
                    "Name": "Println",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 23
                    },
                    "ImportPath": "fmt"
                },
                "len": {
                    "Name": "len",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 22
                    }
                }
            },
            "Args": [
                "str"
            ],
            "Returns": [
                "string"
            ],
            "Complexity": "1 (0 branches, 1 loops)",
            "Pos": 261,
            "Where": {
                "File": "example/example.go",
                "Line": 20
            }
        },
        "text2": {
            "Name": "text2",
            "Doc": "this is text2 function\n",
            "Calls": {
                "GetText": {
                    "Name": "GetText",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 31
                    },
                    "ImportPath": "texts"
                },
                "example2": {
                    "Name": "example2",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 31
                    }
                }
            },
            "Returns": [
                "string"
            ],
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 405,
            "Where": {
                "File": "example/example.go",
                "Line": 30
            }
        }
    },
    "Variables": {
        "EXAMPLE": {
            "Name": "EXAMPLE",
            "Type": "string",
            "Value": "\"mama\"",
            "Keyword": "var",
            "Where": {
                "File": "example/example.go",
                "Line": 8
            }
        },
        "EXAMPLE_1": {
            "Name": "EXAMPLE_1",
            "Type": "",
            "Value": "\"example 1\"",
            "Keyword": "const",
            "Where": {
                "File": "example/example.go",
                "Line": 10
            }
        },
        "EXAMPLE_2": {
            "Name": "EXAMPLE_2",
            "Type": "",
            "Value": "\"example 2\"",
            "Keyword": "const",
            "Where": {
                "File": "example/example.go",
                "Line": 11
            }
        }
    }
}
```

You can query pretty much anything about a `.go` and the features are expanding.