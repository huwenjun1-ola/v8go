package main

import (
	"fmt"
	"gitee.com/hasika/v8go"
)

func main() {
	iso := v8go.NewIsolate()
	defer func() {
		iso.Dispose()
	}()
	fmt.Println("abc")
}
