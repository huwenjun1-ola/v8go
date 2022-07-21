#include "_cgo_export.h"
#include "v8go.h"
void InitV8Go0() {
    InitV8Go(goContext,goFunctionCallback,goTick,goSendMessage);
}