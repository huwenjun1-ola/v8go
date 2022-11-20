// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//go:build windows
// +build windows

package v8go

// #cgo CFLAGS:-I${SRCDIR}/ -w
// #cgo LDFLAGS: -L${SRCDIR}/lib/win_x64 -lv8_export
import "C"

//windows 下使用动态库链接,需要直接链接,不需要编译v8_export
import _ "gitee.com/hasika/v8go/lib/win_x64"
