//go:build !windows
// +build !windows

package v8go

// #cgo CXXFLAGS:-DCOMPILE_EXPORT -fpermissive -fno-rtti -fpic -std=c++14 -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH -I${SRCDIR}/deps/include -I${SRCDIR}/ -w
// #cgo CFLAGS:-DCOMPILE_EXPORT -I${SRCDIR}/deps/include -I${SRCDIR}/ -w
// #cgo LDFLAGS: -pthread -lv8
// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_x86_64
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64
// #cgo linux LDFLAGS: -L${SRCDIR}/lib/linux_x86_64 -ldl
import "C"

import (
	_ "code.flock-block.com/Zheng.Kaikai/v8go/lib/darwin_arm64"
	_ "code.flock-block.com/Zheng.Kaikai/v8go/lib/darwin_x86_64"
	_ "code.flock-block.com/Zheng.Kaikai/v8go/lib/linux_x86_64"
)
