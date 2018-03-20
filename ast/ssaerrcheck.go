// Code inspired by https://github.com/dominikh/go-tools/blob/master/errcheck/errcheck.go

package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "ast/helloworld.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
	}
	files := []*ast.File{f}

	pkg := types.NewPackage("hello", "")

	hello, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()},
		fset, pkg, files, ssa.SanityCheckFunctions,
	)

	mainFunc := hello.Func("main")
	// start-errcheck OMIT
	for _, b := range mainFunc.Blocks {
		for _, ins := range b.Instrs {
			ssacall, ok := ins.(ssa.CallInstruction)
			if !ok {
				continue
			}

			results := ssacall.Common().Signature().Results()
			if errIndex := results.Len() - 1; errIndex >= 0 {
				if last := results.At(errIndex); last.Type().String() != "error" {
					continue
				}

				ssav, ok := ins.(ssa.Value)
				if !ok {
					continue
				}

				if refs := *ssav.Referrers(); len(refs) <= errIndex {
					fmt.Printf("Missing error check at %s: %s\n", ins.String(), fset.Position(ins.Pos()).String())
				}
			}
		}
	}
	// end-errcheck OMIT
}
