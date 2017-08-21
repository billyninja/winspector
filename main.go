package main

import (
	"github.com/billyninja/winspector/probe"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset
	f_ast, err := parser.ParseFile(fset, "../dstructures/main.go", nil, 0)
	if err != nil {
		panic(err)
	}

	probe.GenerateReport(f_ast)
}
