package v8go

//#include <stdlib.h>
//#include "v8go.h"
import "C"

type InspectorServer struct {
	InspectorClientPtr C.RawInspectorClientPtr
}

func (i *InspectorServer) WaitDebugger() {
	for {
		ok := bool(C.InspectorTick(i.InspectorClientPtr))
		if ok {
			break
		}
	}
	go i.run()
}

func (i *InspectorServer) run() {
	for {
		C.InspectorTick(i.InspectorClientPtr)
		closed := bool(C.InspectorAlive(i.InspectorClientPtr))
		if closed {
			return
		}
	}
}

func NewInspectorServer(iso *Isolate, ctx *Context, port uint32) *InspectorServer {
	return &InspectorServer{InspectorClientPtr: C.NewInspectorClient(ctx.ptr, C.int32_t(port))}
}
