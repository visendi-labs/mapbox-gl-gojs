//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_raster/common"
)

var Token = "<token>"

func main() {
	common.ReadFiles()
	example := func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}
	timeline := func(this js.Value, args []js.Value) any {
		return common.UpdateUrl(args[0].String())
	}
	js.Global().Set("example", js.FuncOf(example))
	js.Global().Set("timeline", js.FuncOf(timeline))
	<-make(chan struct{}) // keep running
}
