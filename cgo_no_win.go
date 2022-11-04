//go:build !windows
// +build !windows

package v8go

//// #cgo CXXFLAGS:-DCOMPILE_EXPORT -fpermissive -fno-rtti -fpic -std=c++14 -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH -I${SRCDIR}/deps/include -I${SRCDIR}/ -w
//// #cgo CFLAGS:-DCOMPILE_EXPORT -I${SRCDIR}/deps/include -I${SRCDIR}/ -w
//// #cgo LDFLAGS: -pthread -lv8
//// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_x86_64
//// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64
//// #cgo linux LDFLAGS: -L${SRCDIR}/lib/linux_x86_64 -ldl
//import "C"

// #cgo CFLAGS:-I${SRCDIR}/ -w
// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lv8_export -lv8 -pthread -ldl -lc++
import "C"
import (
	_ "gitee.com/hasika/v8go/lib/darwin_arm64"
	_ "gitee.com/hasika/v8go/lib/darwin_x86_64"
	_ "gitee.com/hasika/v8go/lib/linux_x86_64"
)
