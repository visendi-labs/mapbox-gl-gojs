//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_kmeans/common"
)

var Token = "<token>"

func main() {
	js.Global().Set("example", js.FuncOf(func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}))
	js.Global().Set("kmeanLines", js.FuncOf(func(this js.Value, args []js.Value) any {
		return common.KmeanClusterLines(args[0].String())
	}))
	js.Global().Set("kmeanPoints", js.FuncOf(func(this js.Value, args []js.Value) any {
		return common.KmeanClusterPoints(args[0].String())
	}))
	<-make(chan struct{})
}
