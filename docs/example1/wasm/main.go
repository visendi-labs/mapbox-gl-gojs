//go:build js && wasm

package main

import (
	"example/common"
	"syscall/js"
)

func main() {
	addWasm := func(this js.Value, args []js.Value) any {
		return common.Add(args[0].Int(), args[1].Int())
	}
	js.Global().Set("add", js.FuncOf(addWasm))
	<-make(chan struct{}) // keep running
}
