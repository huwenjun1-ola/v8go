package gots

import (
	"fmt"
	"gitee.com/hasika/v8go"
)

type TsEnv struct {
	Iso           *v8go.Isolate
	Ctx           *v8go.Context
	TsModuleDir   string
	packageConfig *PackageConfig
	globalClass   *v8go.ObjectTemplate
	emptyClass    *v8go.ObjectTemplate
	PreLoadFiles  []string
	Trace         bool
}

func NewTsEnv(tsModuleDir string) *TsEnv {
	pkgConfig, absDirPath, err := extractPackageConfig(tsModuleDir)
	if err != nil {
		return nil
	}
	env := &TsEnv{
		Iso:           nil,
		Ctx:           nil,
		TsModuleDir:   absDirPath,
		packageConfig: pkgConfig,
	}
	return env
}

func (t *TsEnv) Init(waitDebugger bool, debugEnable bool, debugPort uint32) error {
	t.Iso = v8go.NewIsolate()
	t.globalClass = v8go.NewObjectTemplate(t.Iso)
	t.initGlobalClass()
	t.emptyClass = v8go.NewObjectTemplate(t.Iso)
	t.Ctx = v8go.NewContextWithOptions(t.Iso, t.globalClass)
	if debugEnable {
		ins := v8go.NewInspectorServer(t.Iso, t.Ctx, debugPort)
		if waitDebugger {
			ins.WaitDebugger()
		}
	}
	//初始化为node
	err := t.Ctx.Global().Set("global", t.Ctx.Global())
	process, err := t.CreateEmptyObject()
	defer func() {
		if process != nil {
			process.MarkValuePtrCanReleaseInC()
		}
	}()
	if err != nil {
		panic(err)
	}
	err = t.Ctx.Global().Set("process", process)
	if err != nil {
		panic(err)
	}
	version, err := t.CreateEmptyObject()
	defer func() {
		if version != nil {
			version.MarkValuePtrCanReleaseInC()
		}
	}()
	if err != nil {
		panic(err)
	}
	err = process.Set("versions", version)
	if err != nil {
		panic(err)
	}
	node, err := t.CreateEmptyObject()
	defer func() {
		if node != nil {
			node.MarkValuePtrCanReleaseInC()
		}
	}()
	if err != nil {
		panic(err)
	}
	err = version.Set("node", node)
	if err != nil {
		panic(err)
	}
	//log modular 最先运行
	_, err = t.RunScriptWithWrapperByScript("____log.js", nil, logScript, false)
	if err != nil {
		panic(err)
	}
	_, err = t.RunScriptWithWrapperByScript("____modular.js", nil, moduleScript, false)
	if err != nil {
		panic(err)
	}
	for _, f := range t.PreLoadFiles {
		_, err = t.RunScriptWithWrapperByPath(f, nil)
		if err != nil {
			panic(err)
		}
	}
	//进入js main,开始加载运行整个项目
	_, err = t.RunScriptWithWrapperByPath(t.TsModuleDir+"/"+t.packageConfig.Main, entryWrapper)
	if err != nil {
		panic(err)
	}
	t.Iso.TryReleaseValuePtrInC(true)
	return nil
}

func (t *TsEnv) CreateEmptyObject() (*v8go.Object, error) {
	return t.emptyClass.NewInstance(t.Ctx)
}

func (t *TsEnv) Destroy() {
	t.Ctx.Close()
	t.Iso.Dispose()
}

func (t *TsEnv) initGlobalClass() {
	globalClass := t.globalClass
	globalClass.Set("exportToGlobal", v8go.NewFunctionTemplate(t.Iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		n := info.Args()[0]
		ob := info.Args()[1]
		t.Ctx.Global().Set(n.String(), ob)
		return nil
	}))
	_ = globalClass.Set("__log", v8go.NewFunctionTemplate(t.Iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Println(info.Args())
		return nil
	}))
	_ = globalClass.Set("__tgjsEvalScript", v8go.NewFunctionTemplate(t.Iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		source := info.Args()[0].String()
		orig := info.Args()[1].String()
		ret, err := t.RunScriptWithWrapperByScript(orig, nil, source, true)
		if err != nil {
			panic(err)
		}
		defer func() {
			if ret != nil {
				ret.MarkValuePtrCanReleaseInC()
			}
		}()
		return ret
	}))
	_ = globalClass.Set("__tgjsLoadModule", v8go.NewFunctionTemplate(t.Iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		moduleName := info.Args()[0].String()
		requireDir := info.Args()[1].String()
		path := requireDir + "/" + moduleName + ".js"
		path, err := GetAbsPath(path)
		if err != nil {
			panic(err)
		}
		bs, err := t.loadFile(path)
		if err != nil {
			panic(err)
		}
		ret := fmt.Sprintf("%s\n%s\n%s", path, path, string(bs))
		retV, _ := v8go.NewValue(t.Iso, ret)
		defer func() {
			if retV != nil {
				retV.MarkValuePtrCanReleaseInC()
			}
		}()
		return retV
	}))
}

func (t *TsEnv) RunScriptWithWrapperByPath(fileName string, preprocessor ScriptPreProcessor) (*v8go.Value, error) {
	fullPath, err := GetAbsPath(fileName)
	if err != nil {
		return nil, err
	}
	data, err := t.loadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return t.RunScriptWithWrapperByScript(fullPath, preprocessor, string(data), false)
}

func (t *TsEnv) loadFile(moduleName string) ([]byte, error) {
	data, err := ReadFile(moduleName)
	if err != nil {
		t.Print("read file err %s", moduleName)
		return nil, err
	}
	return data, nil
}

func (t *TsEnv) RunScriptWithWrapperByScript(fullPath string, preprocessor ScriptPreProcessor, script string, needResult bool) (*v8go.Value, error) {
	if preprocessor != nil {
		script = preprocessor(script, fullPath)
	}
	v, err := t.Ctx.RunScript(script, fullPath)
	if err != nil {
		t.Print("run script error for %s,error %s", fullPath, err)
		return nil, err
	}
	t.Print("eval script success:%s", fullPath)
	defer func() {
		if !needResult && v != nil {
			if v8go.TraceMem {
				t.Print("Mark Can Be Released By Eval Script")
			}
			v.MarkValuePtrCanReleaseInC()
		}
	}()
	return v, nil
}

func (t *TsEnv) Print(in string, args ...interface{}) {
	if t.Trace {
		fmt.Println(fmt.Sprintf(in, args...))
	}
}
