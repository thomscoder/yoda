package src

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"yoda/types"
)

func GetPackageInfo(filename string) (*types.PackageInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	p := &types.PackageInfo{
		Name:      f.Name.Name,
		Imports:   make(map[string]*types.ImportInfo),
		Functions: make(map[string]*types.FunctionInfo),
		Variables: &types.VariableInfo{},
	}

	// Add imported packages to PackageInfo.
	for _, imp := range f.Imports {
		importPath, _ := strconv.Unquote(imp.Path.Value)
		importName := "."
		if imp.Name != nil {
			importName = imp.Name.Name
		}
		if importName == "." {
			importName = path.Base(importPath)
		}

		pkg, err := build.Import(importPath, "", build.FindOnly)
		if err != nil {
			panic(err)
		}

		size, err := dirSize(pkg.Dir)
		if err != nil {
			fmt.Printf("Error getting size of package directory %s: %v\n", pkg.Dir, err)
			return nil, err
		}
		p.Imports[importName] = &types.ImportInfo{
			Name:     importName,
			Path:     importPath,
			Size:     size,
			AllFuncs: make(map[string]bool), // NEW
			Complete: false,
		}

		for funcName := range p.Imports[importName].AllFuncs {
			p.Imports[importName].Funcs[funcName] = true
		}

		if err != nil {
			fmt.Printf("Error importing package %s: %v\n", importPath, err)
			return nil, err
		}
	}

	ast.Inspect(f, func(node ast.Node) bool {
		// Find variables
		if genDecl, ok := node.(*ast.GenDecl); ok && (genDecl.Tok == token.CONST || genDecl.Tok == token.VAR) {
			// VAR
			if genDecl.Tok == token.VAR {

				if spec, ok := genDecl.Specs[0].(*ast.ValueSpec); ok {
					if spec.Values == nil {
						p.Variables.VarDeclarations = p.Variables.VarDeclarations + 1
					} else {
						p.Variables.Initializations = p.Variables.Initializations + 1
					}
				}
			} else {
				// CONST
				p.Variables.ConstDeclarations = p.Variables.ConstDeclarations + 1
			}
		}

		// Check if the node is a function declaration.
		fn, ok := node.(*ast.FuncDecl)
		if ok {
			// Add the function to the PackageInfo.
			p.Functions[fn.Name.Name] = &types.FunctionInfo{
				Name:       fn.Name.Name,
				Doc:        fn.Doc.Text(),
				Calls:      make(map[string]*types.FunctionInfo),
				Args:       []string{},
				Returns:    []string{},
				Pos:        fn.Pos(),
				Where:      types.WhereInfo{},
				Complexity: "",
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

					posInfo := fset.Position(pos)
					// Add the function call to the types.FunctionInfo.
					p.Functions[fn.Name.Name].Calls[funcName] = &types.FunctionInfo{
						Name:       funcName,
						ImportPath: importPath,
						Where: types.WhereInfo{
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
							p.Functions[fn.Name.Name].Calls[funcName] = &types.FunctionInfo{
								Name:       funcName,
								ImportPath: importPath,
								Where: types.WhereInfo{
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
		return true
	})

	return p, nil
}
