//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example8/common"
)

var Token = "<token>"

func main() {
	common.GeneratePoints()
	example := func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}
	popup := func(this js.Value, args []js.Value) any {
		return common.Popup(args[0].String())
	}
	js.Global().Set("example", js.FuncOf(example))
	js.Global().Set("popup", js.FuncOf(popup))
	<-make(chan struct{}) // keep running
}
