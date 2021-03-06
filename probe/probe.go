package probe

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"strings"
)

var TMAP map[string]uint32
var rootBlocks []*Block

const FSIZE = 14

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

func GenerateReport(f_ast *ast.File) {

	ttf.Init()

	_, cfile, _, _ := runtime.Caller(0)
	path := strings.Replace(cfile, "probe.go", "../assets/fonts/Go-Regular.ttf", 1)
	font, err := ttf.OpenFont(path, FSIZE)
	if err != nil {
		fmt.Printf(">>>>>>>> %v\n", err)
		return
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
										fldBlk.Label = fmt.Sprintf("%s %s", fldBlk.Label, id.String())
									} else {
										fmt.Printf("\t\t\n >>>%s not found \n", id.String())
										complete = false
									}
								case *ast.StarExpr:
									fldBlk.Size = TMAP["ptr"]
									stBlk.Size += fldBlk.Size
									fldBlk.Label = fmt.Sprintf("%s ptr", fldBlk.Label)
								case *ast.SliceExpr:
									fldBlk.Size = TMAP["slice"]
									stBlk.Size += fldBlk.Size
									fldBlk.Label = fmt.Sprintf("%s slice", fldBlk.Label)
								case *ast.ArrayType:
									fldBlk.Size = TMAP["slice"]
									stBlk.Size += fldBlk.Size
									fldBlk.Label = fmt.Sprintf("%s slice", fldBlk.Label)
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

							stBlk.Color = Color{0, 255, 255, 255}
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
				stWidth := int32(rb.Size * FSIZE)

				rect := &sdl.Rect{rb.X + lastX + 30, 25, stWidth, rb.Height}
				rend.SetDrawColor(rb.Color.R, rb.Color.G, rb.Color.B, rb.Color.A)
				rend.FillRect(rect)
				rend.DrawRect(rect)

				labelSurf, _ := font.RenderUTF8_Blended_Wrapped(rb.Label, sdl.Color{255, 255, 255, 255}, int(stWidth))

				labelText, _ := rend.CreateTextureFromSurface(labelSurf)
				rend.Copy(
					labelText,
					&sdl.Rect{0, 0, labelSurf.W, labelSurf.H},
					&sdl.Rect{rect.X - 10, rect.Y - labelSurf.H, labelSurf.W, labelSurf.H},
				)

				var lastXCh int32 = rect.X
				for _, cd := range rb.Children {
					fwidth := int32(cd.Size * FSIZE)
					rend.SetDrawColor(255, 255, 255, 255)
					rend.DrawRect(&sdl.Rect{lastXCh, 25, fwidth, rb.Height})

					labelSurf, _ := font.RenderUTF8_Blended_Wrapped(cd.Label, sdl.Color{255, 255, 255, 255}, int(fwidth))
					labelText, _ := rend.CreateTextureFromSurface(labelSurf)
					rend.Copy(
						labelText,
						&sdl.Rect{0, 0, labelSurf.W, labelSurf.H},
						&sdl.Rect{lastXCh + 5, rect.Y, labelSurf.W, labelSurf.H},
					)

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

	GenerateReport(f_ast)
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
