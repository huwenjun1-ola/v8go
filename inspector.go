package v8go

//#include <stdlib.h>
//#include "v8go.h"
import "C"
import "fmt"

type InspectorServer struct {
	InspectorClientPtr C.RawInspectorClientPtr
}

func (i *InspectorServer) WaitDebugger() {
	for {
		ok := bool(C.InspectorTick(i.InspectorClientPtr))
		if ok {
			fmt.Println("Inspector Connected Now")
			break
		}
	}
}

func (i *InspectorServer) Run() {
	for {
		C.InspectorTick(i.InspectorClientPtr)
		alive := bool(C.InspectorAlive(i.InspectorClientPtr))
		if !alive {
			fmt.Println("Inspector Not Alive,Close Now")
			return
		}
	}
}

func NewInspectorServer(iso *Isolate, ctx *Context, port uint32) *InspectorServer {
	return &InspectorServer{InspectorClientPtr: C.NewInspectorClient(ctx.ptr, C.int32_t(port))}
}
