//go:build linux
// +build linux

package v8go

// #cgo CFLAGS:-I${SRCDIR}/ -w
// #cgo CXXFLAGS:-I${SRCDIR}/ -w
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_x64 -lboost_system -lv8_export -lv8_monolith -lv8_libbase -lv8_libplatform -latomic -ldl -lstdc++ -lpthread -lrt -lm
import "C"
import (
	_ "gitee.com/hasika/v8go/lib/linux_x64"
)
