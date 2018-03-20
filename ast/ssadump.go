package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", os.Stdin, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
	}
	files := []*ast.File{f}

	pkg := types.NewPackage("hello", "")

	hello, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()},
		fset, pkg, files, ssa.SanityCheckFunctions,
	)

	hello.WriteTo(os.Stdout)

	hello.Func("main").WriteTo(os.Stdout)
}
