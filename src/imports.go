package src

import (
	"fmt"
	"go/ast"
	"go/build"
	"path"
	"strconv"
	"yoda/y_types"
)

func getImports(p *y_types.PackageInfo, f *ast.File) error {
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
			return err
		}
		p.Imports[importName] = &y_types.ImportInfo{
			Name:     importName,
			Path:     importPath,
			Size:     size,
			AllFuncs: make(map[string]bool),
			Complete: false,
		}

		for funcName := range p.Imports[importName].AllFuncs {
			p.Imports[importName].Funcs[funcName] = true
		}

		if err != nil {
			fmt.Printf("Error importing package %s: %v\n", importPath, err)
			return err
		}
	}

	return nil
}
