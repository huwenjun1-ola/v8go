package gots

import (
	"fmt"
	"path"
	"strings"
)

type ScriptPreProcessor func(string, string) string

var entryWrapper ScriptPreProcessor = func(script string, fullFileName string) string {
	unixPathName := strings.ReplaceAll(fullFileName, "\\", "/")
	dirName, _ := path.Split(unixPathName)
	wrapped := fmt.Sprintf(`(function (){var __filename = '%s';var __dirname = '%s';var __exports = {};var module = {"exports": __exports, "filename": __filename};(function (exports, require, console,promise) { %s
})(__exports, genRequire('%s'), console);})()`, fullFileName, dirName, script, dirName)
	return wrapped
}
var logScript = `
var global = global || (function () {
    return this;
}());
(function (global) {
    "use strict";
    let tgjsLog = global.__log || function (level, msg) {
    }
    global.__log = undefined;

    const console_org = global.console;
    var console = {}

    function log(level, args) {
        tgjsLog(level, Array.prototype.map.call(args, x => {
            try {
                return x + '';
            } catch (err) {
                return err;
            }
        }).join(','));
    }

    console.log = function () {
        if (console_org) console_org.log.apply(null, Array.prototype.slice.call(arguments));
        log(0, arguments);
    }

    console.info = function () {
        if (console_org) console_org.info.apply(null, Array.prototype.slice.call(arguments));
        log(1, arguments);
    }

    console.warn = function () {
        if (console_org) console_org.warn.apply(null, Array.prototype.slice.call(arguments));
        log(2, arguments);
    }

    console.error = function () {
        if (console_org) console_org.error.apply(null, Array.prototype.slice.call(arguments));
        log(3, arguments);
    }
    global.console = console;
}(global));
`
var moduleScript = `
var global = global || (function () {
    return this;
}());
(function (global) {
    "use strict";
    function normalize(name) {
        if ('./' === name.substr(0, 2)) {
            name = name.substr(2);
        }
        return name;
    }

    let evalScript = global.__tgjsEvalScript || function (script, debugPath) {
        return eval(script);
    }
    global.__tgjsEvalScript = undefined;

    let loadModule = global.__tgjsLoadModule
    global.__tgjsLoadModule = undefined;

    let moduleCache = Object.create(null);

    function executeModule(fullPath, script, debugPath, module) {
        let fullPathInJs = fullPath.replace(/\\/g, '\\\\');
        let fullDirInJs = (fullPath.indexOf('/') != -1) ? fullPath.substring(0, fullPath.lastIndexOf("/")) : fullPath.substring(0, fullPath.lastIndexOf("\\")).replace(/\\/g, '\\\\');
        let exports = {};
        module.exports = exports;
        let wrapped = evalScript(
            // Wrap the script in the same way NodeJS does it. It is important since IDEs (VSCode) will use this wrapper pattern
            // to enable stepping through original source in-place.
            "(function (exports, require, module, __filename, __dirname) { " + script + "\n});",
            debugPath
        )
        wrapped(exports, genRequire(fullDirInJs), module, fullPathInJs, fullDirInJs)
        return module.exports;
    }

    function genRequire(requiringDir) {
		//console.log("genRequire for dir ",requiringDir)
        function require(moduleName) {
			//console.log("require modulename",moduleName,"at dir",requiringDir)
            moduleName = normalize(moduleName);
            let moduleInfo = loadModule(moduleName, requiringDir);
            let split = moduleInfo.indexOf('\n');
            let split2 = moduleInfo.indexOf('\n', split + 1);
            let fullPath = moduleInfo.substring(0, split);
            let debugPath = moduleInfo.substring(split + 1, split2);
            let script = moduleInfo.substring(split2 + 1);
			//console.log(fullPath,debugPath,script)
            let key = fullPath;
            if ((key in moduleCache)) {
                return moduleCache[key].exports;
            }
            let m = {"exports": {}};
            moduleCache[key] = m;
            if (fullPath.endsWith(".json")) {
                let packageConfigure = JSON.parse(script);
                if (fullPath.endsWith("package.json") && packageConfigure.main) {
                    let fullDirInJs = (fullPath.indexOf('/') != -1) ? fullPath.substring(0, fullPath.lastIndexOf("/")) : fullPath.substring(0, fullPath.lastIndexOf("\\")).replace(/\\/g, '\\\\');
                    let tmpRequire = genRequire(fullDirInJs);
                    let r = tmpRequire(packageConfigure.main);

                    m.exports = r;
                } else {

                    m.exports = packageConfigure;
                }
            } else {
                executeModule(fullPath, script, debugPath, m);
            }
            return m.exports;
        }

        return require;
    }
    global.genRequire = genRequire
    global.moduleCache = moduleCache
}(global));

`
