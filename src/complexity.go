package src

import "go/ast"

func getComplexity(fn *ast.FuncDecl) (branches int, loops int) {
	ast.Inspect(fn, func(node ast.Node) bool {
		if _, ok := node.(*ast.BranchStmt); ok {
			branches++
		} else if _, ok := node.(*ast.ForStmt); ok {
			loops++
		} else if _, ok := node.(*ast.RangeStmt); ok {
			loops++
		}
		return true
	})
	return branches, loops
}
