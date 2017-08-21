package main

import (
	"go/parser"
    "go/token"
	"github.com/billyninja/winspector/probe"
)



func main() {
	fset := token.NewFileSet() // positions are relative to fset
	f_ast, err := parser.ParseFile(fset, "../dstructures/main.go", nil, 0)
	if err != nil {
		panic(err)
	}

	probe.GenerateReport(f_ast)
}
