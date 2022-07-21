// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_win

// #cgo CFLAGS:-I${SRCDIR}/include/ -I${SRCDIR}/ -w
// #cgo LDFLAGS: -L${SRCDIR}/lib/ -lv8_export
import "C"

// These imports forces `go mod vendor` to pull in all the folders that
// contain V8 libraries and headers which otherwise would be ignored.
// DO NOT REMOVE
import (
	_ "gitee..com/hasika/v8go/include"
)
