package main

import (
	"fmt"
	"gitee.com/hasika/v8go"
	"runtime"
)

func main() {

	//for i := 0; i < 10000; i++ {
	iso := v8go.NewIsolate()
	ctx := v8go.NewContextWithOptions(iso)
	_, err := ctx.RunScript("var x=new Uint8Array([1,3,4,5]);globalThis.x=x;", "test.js")
	if err != nil {
		panic(err)
	}
	x, xerr := ctx.Global().Get("x")
	if xerr != nil {
		panic(err)
	}
	copied := x.GetCopiedArrayBufferViewContents()
	fmt.Println(len(copied))
	ctx.Close()
	iso.Dispose()
	v8go.CloseAllV8()
	//}
	runtime.GC()
}

func getInUse() uint64 {
	mem2 := &runtime.MemStats{}
	runtime.ReadMemStats(mem2)
	return mem2.Alloc - mem2.Frees
}
