#include "_cgo_export.h"
#include "v8go.h"
void InitV8GoCallBack() {
    InitV8Go(goContext,goFunctionCallback,goTick,goSendMessage);
}

