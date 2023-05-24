package v8go

//#include <stdlib.h>
//#include "v8go.h"
import "C"
import (
	context2 "context"
	"fmt"
	"time"
)

type InspectorServer struct {
	InspectorClientPtr C.RawInspectorClientPtr
	ctx                context2.Context
	cancel             context2.CancelFunc
}

func (i *InspectorServer) WaitDebugger() {
	ticker := time.NewTicker(time.Millisecond * 16)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ticker.C:
			ok := bool(C.InspectorTick(i.InspectorClientPtr))
			if ok {
				fmt.Println("Inspector Connected Now")
				break loop
			}
		case <-i.ctx.Done():
			return
		}
	}

}

func (i *InspectorServer) Run() {
	ticker := time.NewTicker(time.Millisecond * 16)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			C.InspectorTick(i.InspectorClientPtr)
			alive := bool(C.InspectorAlive(i.InspectorClientPtr))
			if !alive {
				fmt.Println("Inspector Not Alive,Close Now")
				return
			}
		case <-i.ctx.Done():
			return
		}
	}
}
func (i *InspectorServer) Destroy() {
	i.cancel()
	C.InspectorClose(i.InspectorClientPtr)
}
func NewInspectorServer(iso *Isolate, ctx *Context, port uint32) *InspectorServer {
	cancel, cancelFunc := context2.WithCancel(context2.Background())
	return &InspectorServer{InspectorClientPtr: C.NewInspectorClient(ctx.ptr, C.int32_t(port)), ctx: cancel, cancel: cancelFunc}
}
