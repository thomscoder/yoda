package src

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"yoda/y_types"
)

func getVariables(p *y_types.PackageInfo, fset *token.FileSet, node ast.Node) {
	// Find variables
	if genDecl, ok := node.(*ast.GenDecl); ok && (genDecl.Tok == token.CONST || genDecl.Tok == token.VAR) {
		if spec, ok := genDecl.Specs[0].(*ast.ValueSpec); ok {

			for idx, name := range spec.Names {
				// Extract the value of the variable
				var value, typeName string
				if spec.Values != nil && len(spec.Values) > idx {
					if lit, ok := spec.Values[idx].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						value = lit.Value
					}
				}

				if spec.Type != nil {
					var buf bytes.Buffer
					if err := printer.Fprint(&buf, token.NewFileSet(), spec.Type); err != nil {
						panic(err)
					} else {
						typeName = buf.String()
					}
				} else {
					typeName = ""
				}

				pos := genDecl.Pos()
				posInfo := fset.Position(pos)

				p.Variables[name.Name] = &y_types.VariableInfo{
					Name:    name.Name,
					Value:   value,
					Type:    typeName,
					Keyword: genDecl.Tok.String(),
					Where: y_types.WhereInfo{
						File: posInfo.Filename,
						Line: posInfo.Line,
					},
				}
			}

		}
	}
}
