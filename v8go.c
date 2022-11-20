#include "v8go.h"
#include "_cgo_export.h"

void InitV8GoCallBack() {
    InitV8Go(goContext, goFunctionCallback);
}
