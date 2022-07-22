package main

import (
	"code.flock-block.com/Zheng.Kaikai/v8go"
	"fmt"
)

func main() {
	iso := v8go.NewIsolate()
	defer func() {
		iso.Dispose()
	}()
	fmt.Println("abc")
}
