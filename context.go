// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"sync"
	"unsafe"
)

// Due to the limitations of passing pointers to C from Go we need to create
// a registry so that we can lookup the Context from any given callback from V8.
// This is similar to what is described here: https://github.com/golang/go/wiki/cgo#function-variables
type ctxRef struct {
	ctx      *Context
	refCount int
}

var ctxMutex sync.RWMutex
var ctxRegistry = make(map[int]*ctxRef)
var ctxSeq = 0

// Context is a GlobalObject root execution environment that allows separate,
// unrelated, JavaScript applications to run in a single instance of V8.
type Context struct {
	ref          int
	ptr          C.ContextPtr
	iso          *Isolate
	GlobalObject *Object
}

type contextOptions struct {
	iso   *Isolate
	gTmpl *ObjectTemplate
}

// ContextOption sets options such as Isolate and Global Template to the NewContext
type ContextOption interface {
	apply(*contextOptions)
}

// NewContextWithOptions creates a new JavaScript context; if no Isolate is passed as a
// ContextOption than a new Isolate will be created.
func NewContextWithOptions(opt ...ContextOption) *Context {
	opts := contextOptions{}
	for _, o := range opt {
		if o != nil {
			o.apply(&opts)
		}
	}

	if opts.iso == nil {
		opts.iso = NewIsolate()
	}

	if opts.gTmpl == nil {
		opts.gTmpl = NewObjectTemplate(opts.iso)
	}
	ctxMutex.Lock()
	ctxSeq++
	ref := ctxSeq
	ctxMutex.Unlock()
	ptr := C.NewContext(opts.iso.ptr, opts.gTmpl.ptr, C.int(ref))
	ctx := NewContext(ref, ptr, opts.iso)
	return ctx
}

func NewContext(ref int, ptr C.ContextPtr, iso *Isolate) *Context {
	ctx := &Context{
		ref: ref,
		ptr: ptr,
		iso: iso,
	}
	valPtr := C.ContextGlobal(ctx.ptr)
	v := NewValueStruct(valPtr, ctx.iso)
	ctx.GlobalObject = &Object{v}
	ctx.register()
	return ctx
}

// Isolate gets the current context's parent isolate.An  error is returned
// if the isolate has been terninated.
func (c *Context) Isolate() *Isolate {
	return c.iso
}

// RunScript executes the source JavaScript; origin (a.k.a. filename) provides a
// reference for the script and used in the stack trace if there is an error.
// error will be of type `JSError` if not nil.
func (c *Context) RunScript(source string, origin string) (*Value, error) {
	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer FreeCPtr(unsafe.Pointer(cSource))
	defer FreeCPtr(unsafe.Pointer(cOrigin))

	rtn := C.RunScript(c.ptr, cSource, cOrigin)
	return valueResult(c.iso, rtn)
}

// Global returns the GlobalObject proxy object.
// Global proxy object is a thin wrapper whose prototype points to actual
// context's GlobalObject object with the properties like Object, etc. This is
// done that way for security reasons.
// Please note that changes to GlobalObject proxy object prototype most probably
// would break the VM — V8 expects only GlobalObject object as a prototype of
// GlobalObject proxy object.
func (c *Context) Global() *Object {
	return c.GlobalObject
}

// PerformMicrotaskCheckpoint runs the default MicrotaskQueue until empty.
// This is used to make progress on Promises.
func (c *Context) PerformMicrotaskCheckpoint() {
	C.IsolatePerformMicrotaskCheckpoint(c.iso.ptr)
}

// Close will dispose the context and free the memory.
// Access to any values associated with the context after calling Close may panic.
func (c *Context) Close() {
	c.deregister()
	C.ContextFree(c.ptr)
	c.ptr = nil
}

func (c *Context) register() {
	ctxMutex.Lock()
	r := ctxRegistry[c.ref]
	if r == nil {
		r = &ctxRef{ctx: c}
		ctxRegistry[c.ref] = r
	}
	r.refCount++
	ctxMutex.Unlock()
}

func (c *Context) deregister() {
	ctxMutex.Lock()
	defer ctxMutex.Unlock()
	r := ctxRegistry[c.ref]
	if r == nil {
		return
	}
	r.refCount--
	if r.refCount <= 0 {
		delete(ctxRegistry, c.ref)
	}
}

func getContext(ref int) *Context {
	ctxMutex.RLock()
	defer ctxMutex.RUnlock()
	r := ctxRegistry[ref]
	if r == nil {
		return nil
	}
	return r.ctx
}

//export goContext
func goContext(ref int) C.ContextPtr {
	ctx := getContext(ref)
	return ctx.ptr
}

func valueResult(ctx *Isolate, rtn C.RtnValue) (*Value, error) {
	if rtn.value == nil {
		return nil, newJSError(rtn.error)
	}
	return NewValueStruct(rtn.value, ctx), nil
}

func objectResult(ctx *Isolate, rtn C.RtnValue) (*Object, error) {
	if rtn.value == nil {
		return nil, newJSError(rtn.error)
	}
	return &Object{NewValueStruct(rtn.value, ctx)}, nil
}
