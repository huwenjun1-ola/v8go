package main

import (
	"gitee.com/hasika/v8go"
	"runtime"
)

func main() {

	for i := 0; i < 10000; i++ {
		iso := v8go.NewIsolate()
		ctx := v8go.NewContext(iso)
		ctx.Close()
		iso.Dispose()
	}
	runtime.GC()
}

func getInUse() uint64 {
	mem2 := &runtime.MemStats{}
	runtime.ReadMemStats(mem2)
	return mem2.Alloc - mem2.Frees
}
