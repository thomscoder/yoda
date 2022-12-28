package y_types

import "go/token"

// PackageInfo represents the information about a package.
type PackageInfo struct {
	Name                      string                    `json:"Name"`
	NumberOfLines             int                       `json:"NumberOfLines"`
	Imports                   map[string]*ImportInfo    `json:"Imports,omitempty"`
	Variables                 map[string]*VariableInfo  `json:"Variables,omitempty"`
	Functions                 map[string]*FunctionInfo  `json:"Functions,omitempty"`
	Structs                   map[string]*StructInfo    `json:"Structs,omitempty"`
	Interfaces                map[string]*InterfaceInfo `json:"Interfaces,omitempty"`
	InvertedFunctionCallGraph map[string][]string       `json:"InvertedFunctionCallGraph,omitempty"`
}

type VariableInfo struct {
	Name    string    `json:"Name"`
	Type    string    `json:"Type"`
	Value   any       `json:"Value,omitempty"`
	Keyword string    `json:"Keyword,omitempty"`
	Where   WhereInfo `json:"Where,omitempty"`
}

type StructInfo struct {
	Name   string                `json:"Name"`
	Fields map[string]*FieldInfo `json:"Fields,omitempty"`
	Where  WhereInfo             `json:"Where,omitempty"`
}

type FieldInfo struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
}

type InterfaceInfo struct {
	Name    string                 `json:"Name"`
	Methods map[string]*MethodInfo `json:"Fields,omitempty"`
	Where   WhereInfo              `json:"Where,omitempty"`
}

type MethodInfo struct {
	Function string
}

// ImportInfo represents the information about an imported package.
type ImportInfo struct {
	Name     string          `json:"Name"`
	Path     string          `json:"Path"`
	Size     int             `json:"Size,omitempty"`
	Funcs    map[string]bool `json:"Funcs,omitempty"`
	AllFuncs map[string]bool `json:"AllFuncs,omitempty"`
	Complete bool            `json:"Complete,omitempty"`
}

// FunctionInfo represents the information about a function.
type FunctionInfo struct {
	Name       string                   `json:"Name"`
	Doc        string                   `json:"Doc,omitempty"`
	Calls      map[string]*FunctionInfo `json:"Calls,omitempty"`
	Args       []string                 `json:"Args,omitempty"`
	Returns    []string                 `json:"Returns,omitempty"`
	Complexity string                   `json:"Complexity,omitempty"`
	Pos        token.Pos                `json:"Pos,omitempty"`
	Where      WhereInfo                `json:"Where,omitempty"`
	ImportPath string                   `json:"ImportPath,omitempty"`
	Complete   bool                     `json:"Complete,omitempty"`
	CalledBy   []string                 `json:"CalledBy,omitempty"`
}

type WhereInfo struct {
	File string
	Line int
}
