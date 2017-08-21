package winspector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	//"go/types"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"reflect"
	"runtime"
	"time"
	//	"github.com/veandco/go-sdl2/sdl_ttf"
	//"os"
	//"bufio"
)

var TMAP map[string]uint32

func Init() {
	TMAP = make(map[string]uint32)
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
}

func Probe() {
	Init()

	_, file, _, ok := runtime.Caller(1)

	if !ok {
		panic("not possible to get to the caller")
	}

	fset := token.NewFileSet() // positions are relative to fset
	f_ast, err := parser.ParseFile(fset, file, nil, 0)
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
							complete := true

							fmt.Printf("\n\nit's a struct...\n")
							fmt.Printf("Fields\n")

							stBlk := &Block{
								Height: 64,
								Label:  ti.Name.String(),
							}

							for _, fld := range tj.Fields.List {

								fldBlk := &Block{
									Height: 64,
									Label:  fld.Names[0].String(),
								}

								fmt.Printf("\t%s %s\n", fld.Names[0], fld.Type)

								switch id := fld.Type.(type) {
								case *ast.Ident:
									if val, ok := TMAP[id.String()]; ok {
										fldBlk.Size = val
										stBlk.Size += fldBlk.Size
									} else {
										fmt.Printf("\t\t\n >>>%s not found \n", id.String())
										complete = false
									}
								case *ast.StarExpr:
									fldBlk.Size = TMAP["ptr"]
									stBlk.Size += fldBlk.Size
									fldBlk.Label += " (ptr)"
								case *ast.SliceExpr:
									fldBlk.Size = TMAP["slice"]
									stBlk.Size += fldBlk.Size
									fldBlk.Label += " (slice)"
								case *ast.ArrayType:
									fldBlk.Size = TMAP["slice"]
									stBlk.Size += fldBlk.Size
									fldBlk.Label += " (slice)"
								default:
									complete = false
									fmt.Printf("\nunknown expr %s\n\n", id)
								}

								stBlk.Children = append(stBlk.Children, fldBlk)
							}
							fmt.Printf("\n\n\nTotal Size: %d, %t\n\n\n", stBlk.Size, complete)
							if complete {
								TMAP[ti.Name.String()] = stBlk.Size
							}
							if len(rootBlocks) == 0 {
								stBlk.Color = Color{255, 0, 0, 255}
							}
							if len(rootBlocks) == 1 {
								stBlk.Color = Color{0, 255, 0, 255}
							}
							if len(rootBlocks) == 2 {
								stBlk.Color = Color{0, 0, 255, 255}
							}

							rootBlocks = append(rootBlocks, stBlk)
						}
					}
				}
			}

			surf, _ := sdl.CreateRGBSurface(0, 1024, 768, 32, 0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
			rend, _ := sdl.CreateSoftwareRenderer(surf)

			rend.SetDrawColor(198, 40, 140, 255)
			rend.FillRect(&sdl.Rect{0, 0, 1024, 768})

			var lastX int32
			for _, rb := range rootBlocks {
				stWidth := int32(rb.Size * 8)

				rect := &sdl.Rect{rb.X + lastX + 30, 10, stWidth, rb.Height}
				rend.SetDrawColor(rb.Color.R, rb.Color.G, rb.Color.B, rb.Color.A)
				rend.FillRect(rect)
				rend.DrawRect(rect)

				var lastXCh int32 = rect.X
				for _, cd := range rb.Children {
					fwidth := int32(cd.Size * 8)
					rend.SetDrawColor(255, 255, 255, 255)
					rend.DrawRect(&sdl.Rect{lastXCh, 10, fwidth, rb.Height})
					lastXCh += fwidth
				}

				lastX = rect.X + stWidth
			}
			img.SavePNG(surf, "test.png")

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
	Size     uint32
	Label    string
	Color    Color
	Parent   *Block
	Children []*Block

	Width  int32
	Height int32
	X      int32
	Y      int32
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

	m := &runtime.MemStats{}
	b := &Block{
		Label: "sadwqeqwewqewqewq",
	}

	for {
		Probe()

		T := reflect.TypeOf(*m)
		insp(T)

		T2 := reflect.TypeOf(*b)
		insp(T2)

		V := reflect.ValueOf(*b)
		fmt.Printf("\nM> %s %+v", V.Kind(), V.NumMethod())

		runtime.ReadMemStats(m)
		fmt.Printf("\nMem: %dMB // %dGC", m.Sys/(1024*1024), m.NumGC)
		fmt.Printf("\nOps: %dMc %dL %dF", m.Mallocs, m.Lookups, m.Frees)
		fmt.Printf("\nHeap: %dKB/%dKB/%dob/%d", m.HeapInuse/1024, m.HeapIdle/1024, m.HeapObjects, m.HeapReleased)
		time.Sleep(1000 * time.Millisecond)
		return
		//runtime.GC()
	}
}
