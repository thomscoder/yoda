package src

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"

	"path"
	"strconv"
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
		Name:          f.Name.Name,
		Imports:       make(map[string]*y_types.ImportInfo),
		Functions:     make(map[string]*y_types.FunctionInfo),
		Variables:     make(map[string]*y_types.VariableInfo),
		NumberOfLines: numLines,
		Structs:       make(map[string]*y_types.StructInfo),
		Interfaces:    make(map[string]*y_types.InterfaceInfo),
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
		p.Imports[importName] = &y_types.ImportInfo{
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
		// Find structs
		if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if tspec, ok := spec.(*ast.TypeSpec); ok {
					pos := genDecl.Pos()
					posInfo := fset.Position(pos)
					if st, ok := tspec.Type.(*ast.StructType); ok {
						fields := make(map[string]*y_types.FieldInfo)

						for _, field := range st.Fields.List {
							var typeName string

							var buf bytes.Buffer
							if err := printer.Fprint(&buf, token.NewFileSet(), field.Type); err != nil {
								panic(err)
							}
							typeName = buf.String()

							// field.Names is a list of *Ident values representing the names of the fields
							// field.Type is an expression representing the type of the field
							// You can use these values to create a FieldInfo and add it to the fields list
							for _, name := range field.Names {
								fields[name.Name] = &y_types.FieldInfo{
									Name: name.Name,
									Type: typeName,
								}
							}
						}

						// tspec.Name is the name of the struct
						// tspec.Type is the *ast.StructType for the struct
						// You can use this information to create a StructInfo
						// and add it to the package info.
						p.Structs[tspec.Name.Name] = &y_types.StructInfo{
							Name:   tspec.Name.Name,
							Fields: fields,
							Where: y_types.WhereInfo{
								File: posInfo.Filename,
								Line: posInfo.Line,
							},
						}
					}
					if in, ok := tspec.Type.(*ast.InterfaceType); ok {
						methods := make(map[string]*y_types.MethodInfo)

						for _, method := range in.Methods.List {
							var typeName string

							var buf bytes.Buffer
							if err := printer.Fprint(&buf, token.NewFileSet(), method.Type); err != nil {
								panic(err)
							}
							typeName = buf.String()

							// field.Names is a list of *Ident values representing the names of the fields
							// field.Type is an expression representing the type of the field
							// You can use these values to create a FieldInfo and add it to the fields list
							for _, name := range method.Names {
								methods[name.Name] = &y_types.MethodInfo{
									Function: typeName,
								}
							}
						}
						// tspec.Name is the name of the interface
						// tspec.Type is the *ast.InterfaceType for the interface
						// You can use this information to create a InterfaceInfo
						// and add it to the package info.
						p.Interfaces[tspec.Name.Name] = &y_types.InterfaceInfo{
							Name:    tspec.Name.Name,
							Methods: methods,
							Where: y_types.WhereInfo{
								File: posInfo.Filename,
								Line: posInfo.Line,
							},
						}
					}
				}
			}
		}
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
		return true
	})

	return p, nil
}
