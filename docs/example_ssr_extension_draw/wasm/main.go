//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_ssr_extension_draw/common"
)

var Token = "<token>"

func main() {
	example := func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}
	addFeatures := func(this js.Value, args []js.Value) any {
		return common.AddFeatures(args[0].String())
	}
	js.Global().Set("example", js.FuncOf(example))
	js.Global().Set("addFeatures", js.FuncOf(addFeatures))
	<-make(chan struct{}) // keep running
}
