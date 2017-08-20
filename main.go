package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	//"go/types"
	"reflect"
	"runtime"
	"time"
	//"os"
	//"bufio"
)

var TMAP map[string]int

func intro(){
	fset := token.NewFileSet()  // positions are relative to fset
	f_ast, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		panic(err)
	}

	for _, dc := range f_ast.Decls {
		switch tt := dc.(type) {
			case *ast.GenDecl:
				if tt.Tok.String() == "type" {
					for _, sp := range tt.Specs {
						switch ti := sp.(type) {
							case *ast.TypeSpec:
								fmt.Printf("%+v\n", ti.Name)
								switch tj := ti.Type.(type) {
									case *ast.StructType:
										stSize := 0
										complete := true

										fmt.Printf("\n\nit's a struct...\n")
										fmt.Printf("Fields\n")
										for _, fld := range tj.Fields.List {
											fmt.Printf(" > %+v", fld.Type)
											fmt.Printf("\t%+v\n", fld.Names[0])
											switch id := fld.Type.(type) {
												case *ast.Ident:
													if val, ok := TMAP[id.String()]; ok {
														stSize += val
													} else {
														fmt.Printf("\t\t\n >>>%s not found \n", id.String())
														complete = false
													}
											}
										}
										fmt.Printf("\n\n\nTotal Size: %d, %t\n\n\n", stSize, complete)
										if complete {
											TMAP[ti.Name.String()] = stSize
										}
								}
						}
					}
				}

				if tt.Tok.String() == "var" {
					for _, sp := range tt.Specs {
						switch ti := sp.(type) {
							case *ast.ValueSpec:
								fmt.Printf("%+v\n", ti)
								fmt.Printf("%+v\n", ti.Names[0])
								fmt.Printf("%+v\n", sp)
						}
					}
				}
		}
	}
}


type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Block struct {
	Width    float32
	Height   float32
	X        float32
	Y        float32
	Label    string
	Color    Color
	Parent   *Block
	Children []*Block
}

func (bc Block) Render(a, b, c int, d float64) (float64, error) {
	return 0.0, nil
}

var rootBlocks []*Block

func insp(T reflect.Type) {
	kind := T.Kind().String()
	fmt.Printf("\n%+v (%s)\n%d-%d-%d", T, kind, T.Align(), T.FieldAlign(), T.Size())
	if kind != "struct" {
		return
	}
	nf := T.NumField()
	fmt.Printf("\nFIELDS: (%d)", nf)
	for i := 0; i < nf; i++ {
		f := T.Field(i)
		Ti := f.Type
		fmt.Printf("\n\t%s: %+v (%s) -- %d-%d-%d", f.Name, Ti, Ti.Kind(), Ti.Align(), Ti.FieldAlign(), Ti.Size())
	}
}

func main() {
	TMAP = make(map[string]int)
	TMAP["uint8"] = 1
	TMAP["bool"] = 1
	TMAP["float32"] = 4
	TMAP["int32"] = 4
	TMAP["uint32"] = 4
	TMAP["uint64"] = 8
	TMAP["int64"] = 8
	TMAP["float64"] = 8
	TMAP["ptr"] = 8
	TMAP["string"] = 16
	TMAP["slice"] = 24

	m := &runtime.MemStats{}
	b := &Block{
		Label: "sadwqeqwewqewqewq",
	}

	for {
		intro()
		return

		T := reflect.TypeOf(*m)
		insp(T)

		T2 := reflect.TypeOf(*b)
		insp(T2)

		V := reflect.ValueOf(*b)
		fmt.Printf("\nM> %s %+v", V.Kind(), V.NumMethod())
		return

		runtime.ReadMemStats(m)
		fmt.Printf("\nMem: %dMB // %dGC", m.Sys/(1024*1024), m.NumGC)
		fmt.Printf("\nOps: %dMc %dL %dF", m.Mallocs, m.Lookups, m.Frees)
		fmt.Printf("\nHeap: %dKB/%dKB/%dob/%d", m.HeapInuse/1024, m.HeapIdle/1024, m.HeapObjects, m.HeapReleased)
		time.Sleep(1000 * time.Millisecond)
		return
		//runtime.GC()
	}
}
