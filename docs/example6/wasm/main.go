//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example6/common"
)

var Token = "<token>"

func main() {
	example := func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}
	filter := func(this js.Value, args []js.Value) any {
		return common.Filter(args[0].String())
	}
	js.Global().Set("example", js.FuncOf(example))
	js.Global().Set("filter", js.FuncOf(filter))
	<-make(chan struct{}) // keep running
}
