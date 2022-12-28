package output

import (
	"fmt"
	"yoda/y_types"
)

func PrintPackageInfo(p *y_types.PackageInfo) {
	fmt.Printf("📦 Package %s:\n", p.Name)
	fmt.Println("  Imports:")
	for _, imp := range p.Imports {
		fmt.Printf("    - %s (%s)\n", imp.Name, imp.Path)
		fmt.Printf("      Size: %dKB\n", imp.Size/1024)
	}
	fmt.Println("  Functions:")
	printFunctions(p.Functions)
	fmt.Println("  Variables:")

}

func printFunctions(functions map[string]*y_types.FunctionInfo) {
	for _, fn := range functions {
		fmt.Printf("    %s %s (%v)\n", "\U0001f1eb", fn.Name, fn.Where)
		fmt.Printf("      Doc: %q\n", fn.Doc)
		fmt.Println("      Args:")
		for _, arg := range fn.Args {
			fmt.Printf("        - %s\n", arg)
		}
		fmt.Println("      Returns:")
		for _, res := range fn.Returns {
			fmt.Printf("        - %s\n", res)
		}
		fmt.Println("      Calls:")
		for _, f := range fn.Calls {
			fmt.Printf("        - %s (%s)\n", f.Name, f.ImportPath)
		}
		fmt.Printf("      Complexity: %s\n", fn.Complexity)

	}
}
