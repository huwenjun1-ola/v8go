package wrapper

import (
	. "gitee.com/hasika/ts-in-go"
	"gitee.com/hasika/v8go"
)

type V8JsEnv struct {
	Iso *v8go.Isolate
	Ctx *v8go.Context
}

func (v *V8JsEnv) Init(debugPort uint32, class IObjectTemplate, waitDebugger bool, debugEnable bool) {
	globalClass := class.(*V8GoObjectTemplateWrapper).template
	v.Ctx = v8go.NewContext(v.Iso, globalClass)
	if debugEnable {
		ins := v8go.NewClient(v.Iso, v.Ctx, debugPort)
		if waitDebugger {
			ins.WaitDebugger()
		}
	}
}

func (v *V8JsEnv) Dispose() {
	v.Ctx.Close()
	v.Iso.Dispose()
}

func (v *V8JsEnv) GetGlobal() IObject {
	return v.Ctx.Global()
}

type V8GoFunction struct {
	fun *v8go.Function
}

func (v V8GoFunction) Call(recv IValue, args ...IValue) (IValue, error) {
	v8args := make([]v8go.Valuer, len(args))
	for i, arg := range args {
		v8args[i] = arg.(v8go.Valuer)
	}
	return v.fun.Call(recv.(v8go.Valuer), v8args...)
}

type V8GoFunctionTemplateWrapper struct {
	template *v8go.FunctionTemplate
	env      *V8JsEnv
}

func (v *V8GoFunctionTemplateWrapper) Set(key string, value interface{}, constraints ...PropertyConstraint) error {
	attributes := make([]v8go.PropertyAttribute, len(constraints))
	for i, constraint := range constraints {
		attributes[i] = v8go.PropertyAttribute(constraint)
	}
	return v.template.Set(key, value, attributes...)
}

func (v *V8GoFunctionTemplateWrapper) GetFunction() IFunction {
	return &V8GoFunction{fun: v.template.GetFunction(v.env.Ctx)}
}

type V8GoObjectTemplateWrapper struct {
	env      *V8JsEnv
	template *v8go.ObjectTemplate
}

func (v *V8GoObjectTemplateWrapper) Set(key string, value interface{}, constraints ...PropertyConstraint) error {
	attributes := make([]v8go.PropertyAttribute, len(constraints))
	for i, constraint := range constraints {
		attributes[i] = v8go.PropertyAttribute(constraint)
	}
	return v.template.Set(key, value, attributes...)
}

func (v *V8GoObjectTemplateWrapper) NewInstance() (IObject, error) {
	return v.template.NewInstance(v.env.Ctx)
}

func (v *V8JsEnv) NewFunctionTemplate(call GoFunCallBack) IFunctionTemplate {
	v8Template := v8go.NewFunctionTemplate(v.Iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := make([]IObject, len(info.Args()))
		for i, value := range info.Args() {
			args[i] = value.Object()
		}
		ret := call(&CommonFunctionCallbackInfo{
			This: info.This(),
			Args: args,
		})
		return ret.(*v8go.Value)
	})
	return &V8GoFunctionTemplateWrapper{
		template: v8Template,
	}
}

func (v *V8JsEnv) NewObjectTemplate() IObjectTemplate {
	v8GoObjectTemplate := v8go.NewObjectTemplate(v.Iso)
	return &V8GoObjectTemplateWrapper{
		template: v8GoObjectTemplate,
	}
}

func (v *V8JsEnv) NewValue(goValue interface{}) (IValue, error) {
	return v8go.NewValue(v.Iso, goValue)
}

func (v *V8JsEnv) RunScript(script string, fullPath string) (IValue, error) {
	return v.Ctx.RunScript(script, fullPath)
}

func (v *V8JsEnv) ConvertUint8ArrayToGoSlice(object IObject) []byte {
	uint8ArrObj := object.(*v8go.Object)
	arrayLenV, err := uint8ArrObj.Get("length")
	if err != nil {
		panic(err)
	}
	arrayLen := arrayLenV.Integer()
	var byteArray []byte = nil
	var index int64
	for index = 0; index < arrayLen; index++ {
		v, _ := uint8ArrObj.GetIdx(uint32(index))
		b := v.Integer()
		byteArray = append(byteArray, byte(b))
	}
	return byteArray
}
