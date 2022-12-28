package src

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"yoda/y_types"
)

func getStructsAndInterfaces(p *y_types.PackageInfo, f *ast.File, fset *token.FileSet, node ast.Node) {

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
}
