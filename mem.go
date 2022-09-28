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

func FreeModuleCPtr(ptr unsafe.Pointer) {
	C.freeV8GoPtr(ptr)
}
