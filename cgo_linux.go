//go:build linux
// +build linux

package v8go

// #cgo CFLAGS:-I${SRCDIR}/ -I${SRCDIR}/v8_export/ -w -stdlib=libc++
// #cgo CXXFLAGS:-I${SRCDIR}/ -I${SRCDIR}/v8_export/ -w -stdlib=libc++
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_x86_64 -lv8_export -lv8_monolith -lv8_libbase -lv8_libplatform -latomic -ldl -lpthread -lrt -lc++ -lc++abi -lm
import "C"
import (
	_ "gitee.com/hasika/v8go/lib/linux_x86_64"
)
