package types

import "go/token"

// PackageInfo represents the information about a package.
type PackageInfo struct {
	Name      string                   `json:"Name"`
	Imports   map[string]*ImportInfo   `json:"Imports,omitempty"`
	Functions map[string]*FunctionInfo `json:"Functions,omitempty"`
	Variables *VariableInfo            `json:"Variables,omitempty"`
}

type VariableInfo struct {
	ConstDeclarations int `json:"ConstDeclarations,omitempty"`
	VarDeclarations   int `json:"VarDeclarations,omitempty"`
	Initializations   int `json:"Initializations,omitempty"`
}

// ImportInfo represents the information about an imported package.
type ImportInfo struct {
	Name     string          `json:"Name"`
	Path     string          `json:"Path"`
	Size     int             `json:"Size,omitempty"`
	Funcs    map[string]bool `json:"Funcs,omitempty"`
	AllFuncs map[string]bool `json:"AllFuncs,omitempty"` // NEW
	Complete bool            `json:"Complete,omitempty"`
}

// FunctionInfo represents the information about a function.
type FunctionInfo struct {
	Name       string                   `json:"Name"`
	Doc        string                   `json:"Doc,omitempty"`
	Calls      map[string]*FunctionInfo `json:"Calls,omitempty"`
	Args       []string                 `json:"Args,omitempty"`    // NEW
	Returns    []string                 `json:"Returns,omitempty"` // NEW
	Complexity string                   `json:"Complexity,omitempty"`
	Pos        token.Pos                `json:"Pos,omitempty"`
	Where      WhereInfo                `json:"Where,omitempty"`
	ImportPath string                   `json:"ImportPath,omitempty"`
	Complete   bool                     `json:"Complete,omitempty"`
}

type WhereInfo struct {
	File string
	Line int
}
