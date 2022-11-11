//go:build darwin
// +build darwin

package v8go

// #cgo CFLAGS:-I${SRCDIR}/ -I${SRCDIR}/v8_export_darwin/ -w -stdlib=libc++
// #cgo CXXFLAGS:-I${SRCDIR}/ -I${SRCDIR}/v8_export_darwin/ -w -stdlib=libc++
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lv8_export -lv8_monolith -lv8_libbase -lv8_libplatform -ldl -lpthread -lc++ -lc++abi -lm
// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_amd64 -lv8_export -lv8_monolith -lv8_libbase -lv8_libplatform -ldl -lpthread -lc++ -lc++abi -lm
import "C"
import (
	_ "gitee.com/hasika/v8go/lib/darwin_amd64"
	_ "gitee.com/hasika/v8go/lib/darwin_arm64"
)
