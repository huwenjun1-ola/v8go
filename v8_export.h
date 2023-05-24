// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef V8_EXPORT_H
#define V8_EXPORT_H
#include <stdbool.h>

#ifdef _WINDOWS
#define V8GO_EXPORT __declspec(dllexport)
#else
#define V8GO_EXPORT
#endif

#ifdef __cplusplus

#include "libplatform/libplatform.h"
#include "v8-profiler.h"
#include "v8.h"
#include "v8-inspector.h"

typedef v8::Isolate *IsolatePtr;
typedef v8::CpuProfiler *CpuProfilerPtr;
typedef v8::CpuProfile *CpuProfilePtr;
typedef const v8::CpuProfileNode *CpuProfileNodePtr;
typedef v8::ScriptCompiler::CachedData *ScriptCompilerCachedDataPtr;

extern "C" {
#else
// Opaque to cgo, but useful to treat it as a pointer to a distinct type
typedef struct v8Isolate v8Isolate;
typedef v8Isolate* IsolatePtr;

typedef struct v8CpuProfiler v8CpuProfiler;
typedef v8CpuProfiler* CpuProfilerPtr;

typedef struct v8CpuProfile v8CpuProfile;
typedef v8CpuProfile* CpuProfilePtr;

typedef struct v8CpuProfileNode v8CpuProfileNode;
typedef const v8CpuProfileNode* CpuProfileNodePtr;

typedef struct v8ScriptCompilerCachedData v8ScriptCompilerCachedData;
typedef const v8ScriptCompilerCachedData* ScriptCompilerCachedDataPtr;
#endif

#include <stddef.h>
#include <stdint.h>

// ScriptCompiler::CompileOptions values
extern V8GO_EXPORT const int ScriptCompilerNoCompileOptions;
extern V8GO_EXPORT const int ScriptCompilerConsumeCodeCache;
extern V8GO_EXPORT const int ScriptCompilerEagerCompile;

typedef struct m_ctx m_ctx;
typedef struct m_value m_value;
typedef struct m_template m_template;
typedef struct m_unboundScript m_unboundScript;

typedef m_ctx *ContextPtr;
typedef m_value *ValuePtr;
typedef m_template *TemplatePtr;
typedef m_unboundScript *UnboundScriptPtr;

typedef struct {
    const char *msg;
    const char *location;
    const char *stack;
} RtnError;

typedef struct {
    UnboundScriptPtr ptr;
    int cachedDataRejected;
    RtnError error;
} RtnUnboundScript;

typedef struct {
    ScriptCompilerCachedDataPtr ptr;
    const uint8_t *data;
    int length;
    int rejected;
} ScriptCompilerCachedData;

typedef struct {
    ScriptCompilerCachedData cachedData;
    int compileOption;
} CompileOptions;

typedef struct {
    CpuProfilerPtr ptr;
    IsolatePtr iso;
} CPUProfiler;

typedef struct CPUProfileNode {
    CpuProfileNodePtr ptr;
    const char *scriptResourceName;
    const char *functionName;
    int lineNumber;
    int columnNumber;
    int childrenCount;
    struct CPUProfileNode **children;
} CPUProfileNode;

typedef struct {
    CpuProfilePtr ptr;
    const char *title;
    CPUProfileNode *root;
    int64_t startTime;
    int64_t endTime;
} CPUProfile;

typedef struct {
    ValuePtr value;
    RtnError error;
} RtnValue;

typedef struct {
    const char *string;
    RtnError error;
} RtnString;

typedef struct {
    size_t total_heap_size;
    size_t total_heap_size_executable;
    size_t total_physical_size;
    size_t total_available_size;
    size_t used_heap_size;
    size_t heap_size_limit;
    size_t malloced_memory;
    size_t external_memory;
    size_t peak_malloced_memory;
    size_t number_of_native_contexts;
    size_t number_of_detached_contexts;
} IsolateHStatistics;

typedef struct {
    const uint64_t *word_array;
    int word_count;
    int sign_bit;
} ValueBigInt;


extern V8GO_EXPORT void
InitV8Go(m_ctx *(*getGoContextFuncEntry)(int), ValuePtr (*goFunctionCallbackEntry)(int, int, ValuePtr *, int));
extern V8GO_EXPORT void Init();

extern V8GO_EXPORT void CloseV8();

extern V8GO_EXPORT IsolatePtr NewIsolate(int ref);

extern V8GO_EXPORT void IsolatePerformMicrotaskCheckpoint(IsolatePtr ptr);

extern V8GO_EXPORT void IsolateDispose(IsolatePtr ptr);

extern V8GO_EXPORT void IsolateTerminateExecution(IsolatePtr ptr);

extern V8GO_EXPORT int IsolateIsExecutionTerminating(IsolatePtr ptr);

extern V8GO_EXPORT IsolateHStatistics IsolationGetHeapStatistics(IsolatePtr ptr);

extern V8GO_EXPORT ValuePtr IsolateThrowException(IsolatePtr iso, ValuePtr value);

extern V8GO_EXPORT RtnUnboundScript IsolateCompileUnboundScript(IsolatePtr iso_ptr,
                                                                const char *source,
                                                                const char *origin,
                                                                CompileOptions options);

extern V8GO_EXPORT ScriptCompilerCachedData *UnboundScriptCreateCodeCache(
        IsolatePtr iso_ptr,
        UnboundScriptPtr us_ptr);

extern V8GO_EXPORT void ScriptCompilerCachedDataDelete(
        ScriptCompilerCachedData *cached_data);

extern V8GO_EXPORT RtnValue UnboundScriptRun(ContextPtr ctx_ptr, UnboundScriptPtr us_ptr);

extern V8GO_EXPORT CPUProfiler *NewCPUProfiler(IsolatePtr iso_ptr);

extern V8GO_EXPORT void CPUProfilerDispose(CPUProfiler *ptr);

extern V8GO_EXPORT void CPUProfilerStartProfiling(CPUProfiler *ptr, const char *title);

extern V8GO_EXPORT CPUProfile *CPUProfilerStopProfiling(CPUProfiler *ptr,
                                                        const char *title);

extern V8GO_EXPORT void CPUProfileDelete(CPUProfile *ptr);

extern V8GO_EXPORT ContextPtr NewContext(IsolatePtr iso_ptr,
                                         TemplatePtr global_template_ptr,
                                         int ref);

extern V8GO_EXPORT void ContextFree(ContextPtr ptr);

extern V8GO_EXPORT RtnValue RunScript(ContextPtr ctx_ptr,
                                      const char *source,
                                      const char *origin);

extern V8GO_EXPORT RtnValue JSONParse(ContextPtr ctx_ptr, const char *str);

extern V8GO_EXPORT const char *JSONStringify(ContextPtr ctx_ptr, ValuePtr val_ptr);

extern V8GO_EXPORT ValuePtr ContextGlobal(ContextPtr ctx_ptr);

extern V8GO_EXPORT void TemplateFreeWrapper(TemplatePtr ptr);

extern V8GO_EXPORT void TemplateSetValue(TemplatePtr ptr,
                                         const char *name,
                                         ValuePtr val_ptr,
                                         int attributes);

extern V8GO_EXPORT void TemplateSetTemplate(TemplatePtr ptr,
                                            const char *name,
                                            TemplatePtr obj_ptr,
                                            int attributes);

extern V8GO_EXPORT TemplatePtr NewObjectTemplate(IsolatePtr iso_ptr);

extern V8GO_EXPORT RtnValue ObjectTemplateNewInstance(TemplatePtr ptr, ContextPtr ctx_ptr);

extern V8GO_EXPORT void ObjectTemplateSetInternalFieldCount(TemplatePtr ptr,
                                                            int field_count);

extern V8GO_EXPORT int ObjectTemplateInternalFieldCount(TemplatePtr ptr);

extern V8GO_EXPORT TemplatePtr NewFunctionTemplate(IsolatePtr iso_ptr, int callback_ref);

extern V8GO_EXPORT RtnValue FunctionTemplateGetFunction(TemplatePtr ptr,
                                                        ContextPtr ctx_ptr);

extern V8GO_EXPORT ValuePtr NewValueNull(IsolatePtr iso_ptr);

extern V8GO_EXPORT ValuePtr NewValueUndefined(IsolatePtr iso_ptr);

extern V8GO_EXPORT ValuePtr NewValueArray(IsolatePtr iso,ValuePtr elements[], int32_t size);

extern V8GO_EXPORT ValuePtr NewValueInteger(IsolatePtr iso_ptr, int32_t v);

extern V8GO_EXPORT ValuePtr NewValueIntegerFromUnsigned(IsolatePtr iso_ptr, uint32_t v);

extern V8GO_EXPORT RtnValue NewValueString(IsolatePtr iso_ptr, const char *v);

extern V8GO_EXPORT ValuePtr NewValueBoolean(IsolatePtr iso_ptr, int v);

extern V8GO_EXPORT ValuePtr NewValueNumber(IsolatePtr iso_ptr, double v);

extern V8GO_EXPORT ValuePtr NewValueBigInt(IsolatePtr iso_ptr, int64_t v);

extern V8GO_EXPORT ValuePtr NewValueBigIntFromUnsigned(IsolatePtr iso_ptr, uint64_t v);

extern V8GO_EXPORT RtnValue NewValueBigIntFromWords(IsolatePtr iso_ptr,
                                                    int sign_bit,
                                                    int word_count,
                                                    const uint64_t *words);

extern V8GO_EXPORT const char *ValueToString(ValuePtr ptr);

extern V8GO_EXPORT const uint32_t *ValueToArrayIndex(ValuePtr ptr);

extern V8GO_EXPORT int ValueToBoolean(ValuePtr ptr);

extern V8GO_EXPORT int32_t ValueToInt32(ValuePtr ptr);

extern V8GO_EXPORT int64_t ValueToInteger(ValuePtr ptr);

extern V8GO_EXPORT double ValueToNumber(ValuePtr ptr);

extern V8GO_EXPORT RtnString ValueToDetailString(ValuePtr ptr);

extern V8GO_EXPORT uint32_t ValueToUint32(ValuePtr ptr);

extern V8GO_EXPORT size_t GetArrayBufferViewByteLen(ValuePtr ptr);

extern V8GO_EXPORT size_t CopyArrayBufferViewContent(ValuePtr ptr, void *dest);

extern V8GO_EXPORT  ValueBigInt ValueToBigInt(ValuePtr ptr);

extern V8GO_EXPORT  RtnValue ValueToObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueSameValue(ValuePtr ptr, ValuePtr otherPtr);

extern V8GO_EXPORT int ValueIsUndefined(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsNull(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsNullOrUndefined(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsTrue(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsFalse(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsName(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsString(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsSymbol(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsFunction(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsBigInt(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsBoolean(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsNumber(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsExternal(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsInt32(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsUint32(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsDate(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsArgumentsObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsBigIntObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsNumberObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsStringObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsSymbolObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsNativeError(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsRegExp(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsAsyncFunction(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsGeneratorFunction(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsGeneratorObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsPromise(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsMap(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsSet(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsMapIterator(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsSetIterator(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsWeakMap(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsWeakSet(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsArray(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsArrayBuffer(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsArrayBufferView(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsTypedArray(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsUint8Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsUint8ClampedArray(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsInt8Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsUint16Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsInt16Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsUint32Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsInt32Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsFloat32Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsFloat64Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsBigInt64Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsBigUint64Array(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsDataView(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsSharedArrayBuffer(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsProxy(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsWasmModuleObject(ValuePtr ptr);

extern V8GO_EXPORT int ValueIsModuleNamespaceObject(ValuePtr ptr);

extern V8GO_EXPORT int GetArrayLen(ValuePtr ptr);

extern V8GO_EXPORT void ObjectSet(ValuePtr ptr, const char *key, ValuePtr val_ptr);

extern V8GO_EXPORT void ObjectSetIdx(ValuePtr ptr, uint32_t idx, ValuePtr val_ptr);

extern V8GO_EXPORT int ObjectSetInternalField(ValuePtr ptr, int idx, ValuePtr val_ptr);

extern V8GO_EXPORT int ObjectInternalFieldCount(ValuePtr ptr);

extern V8GO_EXPORT RtnValue ObjectGet(ValuePtr ptr, const char *key);

extern V8GO_EXPORT RtnValue ObjectGetIdx(ValuePtr ptr, uint32_t idx);

extern V8GO_EXPORT ValuePtr ObjectGetInternalField(ValuePtr ptr, int idx);

extern V8GO_EXPORT int ObjectHas(ValuePtr ptr, const char *key);

extern V8GO_EXPORT int ObjectHasIdx(ValuePtr ptr, uint32_t idx);

extern V8GO_EXPORT int ObjectDelete(ValuePtr ptr, const char *key);

extern V8GO_EXPORT int ObjectDeleteIdx(ValuePtr ptr, uint32_t idx);

extern V8GO_EXPORT RtnValue NewPromiseResolver(ContextPtr ctx_ptr);

extern V8GO_EXPORT ValuePtr PromiseResolverGetPromise(ValuePtr ptr);

extern V8GO_EXPORT int PromiseResolverResolve(ValuePtr ptr, ValuePtr val_ptr);

extern V8GO_EXPORT int PromiseResolverReject(ValuePtr ptr, ValuePtr val_ptr);

extern V8GO_EXPORT int PromiseState(ValuePtr ptr);

extern V8GO_EXPORT RtnValue PromiseThen(ValuePtr ptr, int callback_ref);

extern V8GO_EXPORT RtnValue PromiseThen2(ValuePtr ptr, int on_fulfilled_ref, int on_rejected_ref);

extern V8GO_EXPORT RtnValue PromiseCatch(ValuePtr ptr, int callback_ref);

extern V8GO_EXPORT ValuePtr PromiseResult(ValuePtr ptr);

extern V8GO_EXPORT RtnValue FunctionCall(ValuePtr ptr,
                                         ValuePtr recv,
                                         int argc,
                                         ValuePtr argv[]);

extern V8GO_EXPORT RtnValue FunctionNewInstance(ValuePtr ptr, int argc, ValuePtr args[]);

extern V8GO_EXPORT ValuePtr FunctionSourceMapUrl(ValuePtr ptr);

extern V8GO_EXPORT const char *Version();

extern V8GO_EXPORT void SetFlags(const char *flags);


typedef void *RawInspectorClientPtr;
typedef const char *RawCharPtr;
typedef int *RawIntPtr;

extern V8GO_EXPORT RawInspectorClientPtr NewInspectorClient(ContextPtr ctx, int32_t port);

extern V8GO_EXPORT void InspectorClose(RawInspectorClientPtr ptr);

extern V8GO_EXPORT bool InspectorTick(RawInspectorClientPtr ptr);

extern V8GO_EXPORT bool InspectorAlive(RawInspectorClientPtr ptr);

extern V8GO_EXPORT void freeV8GoPtr(void *p);

extern V8GO_EXPORT void deleteRecordValuePtr(ValuePtr p);

extern V8GO_EXPORT void batchDeleteRecordValuePtr(ValuePtr* p, int n);

extern V8GO_EXPORT void deleteRecordUnboundScriptPtr(UnboundScriptPtr p);

extern V8GO_EXPORT int getCtxRefByValuePtr(ValuePtr ctx);

extern V8GO_EXPORT ContextPtr getDefaultContext(IsolatePtr ctx);

#ifdef __cplusplus
}  // extern V8GO_EXPORT "C"
#endif
#endif  // V8_EXPORT_H
