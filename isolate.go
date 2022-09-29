// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

var v8once sync.Once

// Isolate is a JavaScript VM instance with its own heap and
// garbage collector. Most applications will create one ISO
// with many V8 contexts for execution.
type Isolate struct {
	ptr C.IsolatePtr

	cbMutex sync.RWMutex
	cbSeq   int
	cbs     map[int]FunctionCallback

	null      *Value
	undefined *Value

	templateLock sync.Mutex
	templates    []ITemplate
	InternalCtx  *Context

	canReleasedValuePtrLock sync.Mutex
	canReleasedValuePtrMap  map[C.ValuePtr]interface{}

	tracedValuePtrLock sync.Mutex
	tracedValuePtrMap  map[C.ValuePtr]interface{}

	traceUnboundScriptPtrLock sync.Mutex
	tracedUnboundScriptPtrMap map[C.UnboundScriptPtr]interface{}

	stopLock sync.Mutex
	stopped  bool
}

func (i *Isolate) TraceScriptPtr(ptr C.UnboundScriptPtr) {
	i.traceUnboundScriptPtrLock.Lock()
	defer i.traceUnboundScriptPtrLock.Unlock()
	i.tracedUnboundScriptPtrMap[ptr] = struct{}{}
}

func (i *Isolate) TraceValuePtr(ptr C.ValuePtr) {
	i.tracedValuePtrLock.Lock()
	defer i.tracedValuePtrLock.Unlock()
	i.tracedValuePtrMap[ptr] = struct{}{}
}

func (i *Isolate) MoveTracedPtrToCanReleaseMap(ptr C.ValuePtr, lock bool) {
	if lock {
		i.stopLock.Lock()
		defer i.stopLock.Unlock()
		if i.stopped {
			return
		}
		i.tracedValuePtrLock.Lock()
		defer i.tracedValuePtrLock.Unlock()
		i.canReleasedValuePtrLock.Lock()
		defer i.canReleasedValuePtrLock.Unlock()
	}
	delete(i.tracedValuePtrMap, ptr)
	i.canReleasedValuePtrMap[ptr] = struct{}{}
	if TraceMem {
		fmt.Println("MoveTracedPtrToCanReleaseMap,Can Release Map Size", len(i.canReleasedValuePtrMap), "Still Tracing Size", len(i.tracedValuePtrMap))
	}
}

func (i *Isolate) TryReleaseValuePtrInC(lock bool) {
	if lock {
		i.canReleasedValuePtrLock.Lock()
		defer i.canReleasedValuePtrLock.Unlock()
	}
	l := len(i.canReleasedValuePtrMap)
	if l <= 0 {
		return
	}
	valuePointers := make([]C.ValuePtr, l)
	index := 0
	for ptr := range i.canReleasedValuePtrMap {
		valuePointers[index] = ptr
		index++
	}
	C.batchDeleteRecordValuePtr(&valuePointers[0], C.int(len(valuePointers)))
	runtime.KeepAlive(valuePointers)
	i.canReleasedValuePtrMap = map[C.ValuePtr]interface{}{}
	if TraceMem {
		fmt.Println("Real Released Cnt ", l, ",Traced Not Released", i.GetTracedValueCnt())
	}
}

func (i *Isolate) releaseTracedValuePtrInC(lock bool) {
	if lock {
		i.tracedValuePtrLock.Lock()
		defer i.tracedValuePtrLock.Unlock()
	}
	l := len(i.tracedValuePtrMap)
	if l <= 0 {
		return
	}
	valuePointers := make([]C.ValuePtr, l)
	index := 0
	for ptr := range i.tracedValuePtrMap {
		valuePointers[index] = ptr
		index++
	}
	C.batchDeleteRecordValuePtr(&valuePointers[0], C.int(len(valuePointers)))
	runtime.KeepAlive(valuePointers)
	i.tracedValuePtrMap = map[C.ValuePtr]interface{}{}
}

func (i *Isolate) GetTracedValueCnt() int {
	return len(i.tracedValuePtrMap) + len(i.canReleasedValuePtrMap)
}

func (i *Isolate) releaseAllValuePtrInC() {
	i.tracedValuePtrLock.Lock()
	defer i.tracedValuePtrLock.Unlock()
	i.canReleasedValuePtrLock.Lock()
	defer i.canReleasedValuePtrLock.Unlock()
	i.TryReleaseValuePtrInC(false)
	i.releaseTracedValuePtrInC(false)
}

func (i *Isolate) releaseScriptsInC() {
	i.traceUnboundScriptPtrLock.Lock()
	defer i.traceUnboundScriptPtrLock.Unlock()
	for ptr, _ := range i.tracedUnboundScriptPtrMap {
		C.deleteRecordUnboundScriptPtr(ptr)
		delete(i.tracedUnboundScriptPtrMap, ptr)
	}
	i.tracedUnboundScriptPtrMap = map[C.UnboundScriptPtr]interface{}{}
}

// HeapStatistics represents V8 ISO heap statistics
type HeapStatistics struct {
	TotalHeapSize            uint64
	TotalHeapSizeExecutable  uint64
	TotalPhysicalSize        uint64
	TotalAvailableSize       uint64
	UsedHeapSize             uint64
	HeapSizeLimit            uint64
	MallocedMemory           uint64
	ExternalMemory           uint64
	PeakMallocedMemory       uint64
	NumberOfNativeContexts   uint64
	NumberOfDetachedContexts uint64
}

// NewIsolate creates a new V8 ISO. Only one thread may access
// a given ISO at a time, but different threads may access
// different isolates simultaneously.
// When a ISO is no longer used its resources should be freed
// by calling iso.Dispose().
// An *Isolate can be used as a v8go.ContextOption to create a new
// Context, rather than creating a new default Isolate.
func NewIsolate() *Isolate {
	v8once.Do(func() {
		C.Init()
		C.InitV8GoCallBack()
	})
	var ref int
	ctxMutex.Lock()
	ctxSeq++
	ref = ctxSeq
	ctxMutex.Unlock()
	iso := &Isolate{
		ptr:                       C.NewIsolate(C.int(ref)),
		cbs:                       make(map[int]FunctionCallback),
		tracedValuePtrMap:         map[C.ValuePtr]interface{}{},
		canReleasedValuePtrMap:    map[C.ValuePtr]interface{}{},
		tracedUnboundScriptPtrMap: map[C.UnboundScriptPtr]interface{}{},
	}
	contextPtr := C.getDefaultContext(iso.ptr)
	ctx := NewContext(ref, contextPtr, iso)
	iso.InternalCtx = ctx
	iso.null = newValueNull(iso)
	iso.undefined = newValueUndefined(iso)
	return iso
}

// TerminateExecution terminates forcefully the current thread
// of JavaScript execution in the given ISO.
func (i *Isolate) TerminateExecution() {
	C.IsolateTerminateExecution(i.ptr)
}

// IsExecutionTerminating returns whether V8 is currently terminating
// Javascript execution. If true, there are still JavaScript frames
// on the stack and the termination exception is still active.
func (i *Isolate) IsExecutionTerminating() bool {
	return C.IsolateIsExecutionTerminating(i.ptr) == 1
}

type CompileOptions struct {
	CachedData *CompilerCachedData

	Mode CompileMode
}

// CompileUnboundScript will create an UnboundScript (i.e. context-indepdent)
// using the provided source JavaScript, origin (a.k.a. filename), and options.
// If options contain a non-null CachedData, compilation of the script will use
// that code cache.
// error will be of type `JSError` if not nil.
func (i *Isolate) CompileUnboundScript(source, origin string, opts CompileOptions) (ret *UnboundScript, err error) {
	defer func() {
		if ret != nil {
			i.TraceScriptPtr(ret.ptr)
		}
	}()
	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer FreeCPtr(unsafe.Pointer(cSource))
	defer FreeCPtr(unsafe.Pointer(cOrigin))
	defer FreeCPtr(unsafe.Pointer(cOrigin))

	var cOptions C.CompileOptions
	if opts.CachedData != nil {
		if opts.Mode != 0 {
			panic("On CompileOptions, Mode and CachedData can't both be set")
		}
		cOptions.compileOption = C.ScriptCompilerConsumeCodeCache
		cOptions.cachedData = C.ScriptCompilerCachedData{
			data:   (*C.uchar)(unsafe.Pointer(&opts.CachedData.Bytes[0])),
			length: C.int(len(opts.CachedData.Bytes)),
		}
	} else {
		cOptions.compileOption = C.int(opts.Mode)
	}

	rtn := C.IsolateCompileUnboundScript(i.ptr, cSource, cOrigin, cOptions)
	if rtn.ptr == nil {
		return nil, newJSError(rtn.error)
	}
	if opts.CachedData != nil {
		opts.CachedData.Rejected = int(rtn.cachedDataRejected) == 1
	}

	return &UnboundScript{
		ptr: rtn.ptr,
		iso: i,
	}, nil
}

// GetHeapStatistics returns heap statistics for an ISO.
func (i *Isolate) GetHeapStatistics() HeapStatistics {
	hs := C.IsolationGetHeapStatistics(i.ptr)

	return HeapStatistics{
		TotalHeapSize:            uint64(hs.total_heap_size),
		TotalHeapSizeExecutable:  uint64(hs.total_heap_size_executable),
		TotalPhysicalSize:        uint64(hs.total_physical_size),
		TotalAvailableSize:       uint64(hs.total_available_size),
		UsedHeapSize:             uint64(hs.used_heap_size),
		HeapSizeLimit:            uint64(hs.heap_size_limit),
		MallocedMemory:           uint64(hs.malloced_memory),
		ExternalMemory:           uint64(hs.external_memory),
		PeakMallocedMemory:       uint64(hs.peak_malloced_memory),
		NumberOfNativeContexts:   uint64(hs.number_of_native_contexts),
		NumberOfDetachedContexts: uint64(hs.number_of_detached_contexts),
	}
}

// Dispose will dispose the Isolate VM; subsequent calls will panic.
func (i *Isolate) Dispose() {
	i.stopLock.Lock()
	defer i.stopLock.Unlock()
	i.stopped = true
	i.releaseAllValuePtrInC()
	i.releaseScriptsInC()
	if i.ptr == nil {
		return
	}
	i.templateLock.Lock()
	defer i.templateLock.Unlock()
	for _, tpl := range i.templates {
		tpl.finalizer()
	}
	i.templates = nil
	C.IsolateDispose(i.ptr)
	i.ptr = nil
}

// ThrowException schedules an exception to be thrown when returning to
// JavaScript. When an exception has been scheduled it is illegal to invoke
// any JavaScript operation; the caller must return immediately and only after
// the exception has been handled does it become legal to invoke JavaScript operations.
func (i *Isolate) ThrowException(value *Value) *Value {
	if i.ptr == nil {
		panic("Isolate has been disposed")
	}
	return NewValueStruct(C.IsolateThrowException(i.ptr, value.ptr), i)
}

// Deprecated: use `iso.Dispose()`.
func (i *Isolate) Close() {
	i.Dispose()
}

func (i *Isolate) apply(opts *contextOptions) {
	opts.iso = i
}

func (i *Isolate) registerCallback(cb FunctionCallback) int {
	i.cbMutex.Lock()
	i.cbSeq++
	ref := i.cbSeq
	i.cbs[ref] = cb
	i.cbMutex.Unlock()
	return ref
}

func (i *Isolate) getCallback(ref int) FunctionCallback {
	i.cbMutex.RLock()
	defer i.cbMutex.RUnlock()
	return i.cbs[ref]
}

func (i *Isolate) BatchMarkCanReleaseInC(values ...*Value) {
	i.stopLock.Lock()
	defer i.stopLock.Unlock()
	if i.stopped {
		return
	}
	i.tracedValuePtrLock.Lock()
	defer i.tracedValuePtrLock.Unlock()
	i.canReleasedValuePtrLock.Lock()
	defer i.canReleasedValuePtrLock.Unlock()
	for _, v := range values {
		runtime.SetFinalizer(v, nil)
		i.MoveTracedPtrToCanReleaseMap(v.ptr, false)
	}
}
