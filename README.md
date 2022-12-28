# Yoda

A Go/WASI tool to get a quick summary of `.go` files.

I something that helped me wrap my head around large Golang codebases, therefore I started building this tool.

## Usage
```bash
yoda <FILE_NAME>
```

## Output sample (based on `example/example.go`)

Bash output (useful for quick testing)

```bash
📦 Package example:
  Imports:
    - fmt (fmt)
      Size: 207KB
    - texts (yoda/example/texts)
      Size: 0KB
  Functions:
    🇫 example (File: example/example.go, line: 14)
      Doc: "This is an example\n"
      Args:
      Returns:
      Calls:
        - Println (fmt)
        - text ()
      Complexity: 0 (0 branches, 0 loops)
    🇫 text (File: example/example.go, line: 19)
      Doc: ""
      Args:
        - str
      Returns:
        - string
      Calls:
        - len ()
        - Println (fmt)
        - GetText (texts)
      Complexity: 1 (0 branches, 1 loops)
    🇫 text2 (File: example/example.go, line: 28)
      Doc: ""
      Args:
      Returns:
        - string
      Calls:
        - example2 ()
        - GetText (texts)
      Complexity: 0 (0 branches, 0 loops)
  Variables:
    - Const:
       - Declarations: 2
    - Var:
       - Declarations: 1
    - Initializations: 0
```
JSON output
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
                    "Where": "File: example/example.go, line: 15",
                    "ImportPath": "fmt"
                },
                "text": {
                    "Name": "text",
                    "Where": "File: example/example.go, line: 16"
                }
            },
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 162,
            "Where": "File: example/example.go, line: 14"
        },
        "text": {
            "Name": "text",
            "Calls": {
                "GetText": {
                    "Name": "GetText",
                    "Where": "File: example/example.go, line: 25",
                    "ImportPath": "texts"
                },
                "Println": {
                    "Name": "Println",
                    "Where": "File: example/example.go, line: 22",
                    "ImportPath": "fmt"
                },
                "len": {
                    "Name": "len",
                    "Where": "File: example/example.go, line: 21"
                }
            },
            "Args": [
                "str"
            ],
            "Returns": [
                "string"
            ],
            "Complexity": "1 (0 branches, 1 loops)",
            "Pos": 227,
            "Where": "File: example/example.go, line: 19"
        },
        "text2": {
            "Name": "text2",
            "Calls": {
                "GetText": {
                    "Name": "GetText",
                    "Where": "File: example/example.go, line: 29",
                    "ImportPath": "texts"
                },
                "example2": {
                    "Name": "example2",
                    "Where": "File: example/example.go, line: 29"
                }
            },
            "Returns": [
                "string"
            ],
            "Complexity": "0 (0 branches, 0 loops)",
            "Pos": 345,
            "Where": "File: example/example.go, line: 28"
        }
    },
    "Variables": {
        "ConstDeclarations": 2,
        "VarDeclarations": 1
    }
}
```
