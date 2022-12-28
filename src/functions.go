package src

import (
	"fmt"
	"go/ast"
	"go/token"
	"yoda/y_types"
)

func getFunctions(p *y_types.PackageInfo, fset *token.FileSet, node ast.Node) {
	// Check if the node is a function declaration.
	fn, ok := node.(*ast.FuncDecl)
	if ok {
		// Add the function to the PackageInfo.
		p.Functions[fn.Name.Name] = &y_types.FunctionInfo{
			Name:       fn.Name.Name,
			Doc:        fn.Doc.Text(),
			Calls:      make(map[string]*y_types.FunctionInfo),
			Args:       []string{},
			Returns:    []string{},
			Pos:        fn.Pos(),
			Where:      y_types.WhereInfo{},
			Complexity: "",
			CalledBy:   []string{},
			ImportPath: "",
			Complete:   false,
		}

		pos := p.Functions[fn.Name.Name].Pos
		posInfo := fset.Position(pos)
		p.Functions[fn.Name.Name].Where.File = posInfo.Filename
		p.Functions[fn.Name.Name].Where.Line = posInfo.Line

		// Get the complexity of the function
		branches, loops := getComplexity(fn)
		p.Functions[fn.Name.Name].Complexity = fmt.Sprintf("%d (%d branches, %d loops)", branches+loops, branches, loops)

		// Extract the function's argument names.
		for _, field := range fn.Type.Params.List {
			for _, name := range field.Names {
				p.Functions[fn.Name.Name].Args = append(p.Functions[fn.Name.Name].Args, name.Name)
			}
		}

		// Extract the function's return values.
		if fn.Type.Results != nil {
			for _, field := range fn.Type.Results.List {
				returnType := ""
				switch t := field.Type.(type) {
				case *ast.Ident:
					returnType = t.Name
				case *ast.SelectorExpr:
					returnType = t.Sel.Name
				case *ast.StarExpr:
					if _type, ok := t.X.(*ast.Ident); ok {
						returnType = _type.Name
					}
				}

				p.Functions[fn.Name.Name].Returns = append(p.Functions[fn.Name.Name].Returns, returnType)
			}
		}

		ast.Inspect(fn, func(node ast.Node) bool {
			var currentFunctionName = fn.Name.Name
			switch node := node.(type) {
			case *ast.CallExpr:

				var funcName string
				var importPath string
				var pos token.Pos
				switch fun := node.Fun.(type) {
				case *ast.Ident:
					funcName = fun.Name
					importPath = ""
					pos = fun.Pos()
				case *ast.SelectorExpr:
					importPath = fun.X.(*ast.Ident).Name

					funcName = fun.Sel.Name
					pos = fun.Pos()
				}

				// Functions call graph
				// just to populate the CalledBy for each function. Gets deleted once done
				p.InvertedFunctionCallGraph[funcName] = append(p.InvertedFunctionCallGraph[funcName], currentFunctionName)

				posInfo := fset.Position(pos)
				// Add the function call to the y_types.FunctionInfo.
				p.Functions[fn.Name.Name].Calls[funcName] = &y_types.FunctionInfo{
					Name:       funcName,
					ImportPath: importPath,
					Where: y_types.WhereInfo{
						File: posInfo.Filename,
						Line: posInfo.Line,
					},
				}

			case *ast.SelectorExpr:
				// Check if the function call is using an imported package.
				if sel, ok := node.X.(*ast.Ident); ok {
					// Extract the name of the imported package and function.
					importPath := sel.Name
					funcName := node.Sel.Name
					pos := node.Sel.Pos()

					posInfo := fset.Position(pos)

					if _, ok := p.Imports[importPath]; ok {
						p.Functions[fn.Name.Name].Calls[funcName] = &y_types.FunctionInfo{
							Name:       funcName,
							ImportPath: importPath,
							Where: y_types.WhereInfo{
								File: posInfo.Filename,
								Line: posInfo.Line,
							},
						}
					}
				}
			}

			return true
		})

	}
}
