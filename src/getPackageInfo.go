package src

import (
	"go/ast"
	"go/parser"
	"go/token"

	"yoda/y_types"
)

func GetPackageInfo(filename string) (*y_types.PackageInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	numLines := fset.File(f.Pos()).LineCount()
	if err != nil {
		return nil, err
	}

	p := &y_types.PackageInfo{
		Name:                      f.Name.Name,
		Imports:                   make(map[string]*y_types.ImportInfo),
		Functions:                 make(map[string]*y_types.FunctionInfo),
		Variables:                 make(map[string]*y_types.VariableInfo),
		NumberOfLines:             numLines,
		Structs:                   make(map[string]*y_types.StructInfo),
		Interfaces:                make(map[string]*y_types.InterfaceInfo),
		InvertedFunctionCallGraph: make(map[string][]string),
	}

	getImports(p, f)

	ast.Inspect(f, func(node ast.Node) bool {
		getVariables(p, fset, node)
		getFunctions(p, fset, node)
		getStructsAndInterfaces(p, f, fset, node)
		return true
	})

	// Populate the calledby
	for _, fun := range p.Functions {
		calledBy, ok := p.InvertedFunctionCallGraph[fun.Name]
		if !ok {
			calledBy = []string{}
		}
		p.Functions[fun.Name].CalledBy = calledBy
	}

	p.InvertedFunctionCallGraph = nil

	return p, nil
}
