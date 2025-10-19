//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example4/common"
)

func main() {
	example := func(this js.Value, args []js.Value) any {
		return common.Example()
	}
	js.Global().Set("example", js.FuncOf(example))
	<-make(chan struct{}) // keep running
}
