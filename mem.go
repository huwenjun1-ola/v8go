package v8go

//#include <stdlib.h>
//#include "v8go.h"
import "C"
import (
	"unsafe"
)

func FreeCPtr(ptr unsafe.Pointer) {
	C.free(ptr)
}
func FreeAnyCPtr(ptr unsafe.Pointer) {
	C.freeAny(ptr)
}
