// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

// FunctionCallback is a callback that is executed in Go when a function is executed in JS.
type FunctionCallback func(info *FunctionCallbackInfo) *Value

// FunctionCallbackInfo is the argument that is passed to a FunctionCallback.
type FunctionCallbackInfo struct {
	ctx  *Context
	args []*Value
	this *Object
}

// Context is the current context that the callback is being executed in.
func (i *FunctionCallbackInfo) Context() *Context {
	return i.ctx
}

// This returns the receiver object "this".
func (i *FunctionCallbackInfo) This() *Object {
	return i.this
}

// Args returns a slice of the value arguments that are passed to the JS function.
func (i *FunctionCallbackInfo) Args() []*Value {
	return i.args
}

// FunctionTemplate is used to create functions at runtime.
// There can only be one function created from a FunctionTemplate in a context.
// The lifetime of the created function is equal to the lifetime of the context.
type FunctionTemplate struct {
	*template
}

// NewFunctionTemplate creates a FunctionTemplate for a given callback.
func NewFunctionTemplate(iso *Isolate, callback FunctionCallback) *FunctionTemplate {
	if iso == nil {
		panic("nil Isolate argument not supported")
	}
	if callback == nil {
		panic("nil FunctionCallback argument not supported")
	}

	cbref := iso.registerCallback(callback)

	tmpl := &template{
		ptr:  C.NewFunctionTemplate(iso.ptr, C.int(cbref)),
		iso:  iso,
		Name: time.Now().String(),
	}
	iso.templateLock.Lock()
	defer iso.templateLock.Unlock()
	iso.templates = append(iso.templates, tmpl)
	return &FunctionTemplate{tmpl}
}

// GetFunction returns an instance of this function template bound to the given context.
func (tmpl *FunctionTemplate) GetFunction(ctx *Context) *Function {
	rtn := C.FunctionTemplateGetFunction(tmpl.ptr, ctx.ptr)
	val, err := valueResult(ctx.iso, rtn)
	if err != nil {
		panic(err) // TODO: Consider returning the error
	}
	return &Function{val}
}

// Note that ideally `thisAndArgs` would be split into two separate arguments, but they were combined
// to workaround an ERROR_COMMITMENT_LIMIT error on windows that was detected in CI.
//export goFunctionCallback
func goFunctionCallback(ctxref int, cbref int, thisAndArgs *C.ValuePtr, argsCount int) C.ValuePtr {
	ctx := getContext(ctxref)

	this := *thisAndArgs
	info := &FunctionCallbackInfo{
		ctx:  ctx,
		this: &Object{Value: NewValueStruct(this, ctx.iso)},
		args: make([]*Value, argsCount),
	}
	defer func() {
		bt := make([]*Value, 0)
		bt = append(bt, info.this.Value)
		bt = append(bt, info.args...)
		if TraceMem {
			fmt.Println("Mark Can Be Released By Func Call This And Args,Len ", len(bt))
		}
		ctx.iso.BatchMarkCanReleaseInC(bt...)
	}()
	argv := (*[1 << 30]C.ValuePtr)(unsafe.Pointer(thisAndArgs))[1 : argsCount+1 : argsCount+1]
	for i, v := range argv {
		val := NewValueStruct(v, ctx.iso)
		info.args[i] = val
	}

	callbackFunc := ctx.iso.getCallback(cbref)
	if val := callbackFunc(info); val != nil {
		return val.ptr
	}
	return nil
}
