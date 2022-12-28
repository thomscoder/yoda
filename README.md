# Yoda 👽✨

Get to know your Golang files.

A Go tool to programmatically query `.go` files.

I just wanted something that helped me wrap my head around large Golang codebases, therefore I started building this tool.

## Usage

Yoda reads a Go file and outputs a JSON, that you can view with your favorite tool or library

- Clone a repository or create a Go application
- Install what you have to

```bash
yoda --file <FILE_NAME> --output <FILE_OUTPUT>
```

Example

```bash
yoda --file example/example.go --output yoda.json
```

## Output sample (based on `example/example.go`)

Given this input

```go
package example

import "fmt"

func sayHelloWorld() {
	fmt.Println("Hello world")
}
```

Yoda will generate the following output

```json
{
    "Name": "example",
    "NumberOfLines": 7,
    "Imports": {
        "fmt": {
            "Name": "fmt",
            "Path": "fmt",
            "Size": 212331
        }
    },
    "Functions": {
        "sayHelloWorld": {
            "Name": "sayHelloWorld",
            "Calls": {
                "Println": {
                    "Name": "Println",
                    "Where": {
                        "File": "example/small.go",
                        "Line": 6
                    },
                    "ImportPath": "fmt"
                }
            },
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 32,
            "Where": {
                "File": "example/small.go",
                "Line": 5
            }
        }
    }
}
```

It can also analyze more complex examples such as

```go
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
```

will translate to 

```json
{
    "Name": "example",
    "NumberOfLines": 45,
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
        },
        "str": {
            "Name": "str",
            "Type": "string",
            "Value": "",
            "Keyword": "var",
            "Where": {
                "File": "example/example.go",
                "Line": 36
            }
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
                        "Line": 20
                    },
                    "ImportPath": "fmt"
                },
                "text": {
                    "Name": "text",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 21
                    }
                }
            },
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 221,
            "Where": {
                "File": "example/example.go",
                "Line": 19
            },
            "CalledBy": [
                "text1"
            ]
        },
        "text": {
            "Name": "text",
            "Doc": "This is text function\n",
            "Calls": {
                "GetText": {
                    "Name": "GetText",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 31
                    },
                    "ImportPath": "texts"
                },
                "Println": {
                    "Name": "Println",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 28
                    },
                    "ImportPath": "fmt"
                },
                "len": {
                    "Name": "len",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 27
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
            "Pos": 311,
            "Where": {
                "File": "example/example.go",
                "Line": 25
            },
            "CalledBy": [
                "example"
            ]
        },
        "text1": {
            "Name": "text1",
            "Calls": {
                "example": {
                    "Name": "example",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 44
                    }
                }
            },
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 585,
            "Where": {
                "File": "example/example.go",
                "Line": 43
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
                        "Line": 40
                    },
                    "ImportPath": "texts"
                },
                "example2": {
                    "Name": "example2",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 40
                    }
                },
                "len": {
                    "Name": "len",
                    "Where": {
                        "File": "example/example.go",
                        "Line": 37
                    }
                }
            },
            "Returns": [
                "string"
            ],
            "Complexity": "1 (1 branches, 0 loops)",
            "Pos": 455,
            "Where": {
                "File": "example/example.go",
                "Line": 35
            }
        }
    },
    "Structs": {
        "mama": {
            "Name": "mama",
            "Fields": {
                "Name": {
                    "Name": "Name",
                    "Type": "string"
                },
                "Value": {
                    "Name": "Value",
                    "Type": "string"
                }
            },
            "Where": {
                "File": "example/example.go",
                "Line": 13
            }
        }
    }
}
```

and so on.

## Contributing

Any contribution is welcomed!!
Bug fix, features, docs... 
May the force be with you!