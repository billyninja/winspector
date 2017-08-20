package main

import (
    "runtime"
    "time"
    "fmt"
    "reflect"
)

type Color struct {
    R uint8
    G uint8
    B uint8
    A uint8
}

type Block struct {
    Width       float32
    Height      float32
    X           float32
    Y           float32
    Label       string
    Color       Color
    Parent      *Block
    Children    []*Block
}

func (bc Block) Render(a, b, c int, d float64) (float64, error) {
    return 0.0, nil
}

var rootBlocks []*Block

func insp(T reflect.Type) {
    fmt.Printf("\n%+v (%s)\n%d-%d-%d", T, T.Kind(), T.Align(), T.FieldAlign(), T.Size())
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
    b := &Block{}

    for {

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
        //runtime.GC()
    }
}
