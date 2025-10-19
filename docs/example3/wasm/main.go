//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example3/common"
)

var Token = "<token>"

func main() {
	example := func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}
	js.Global().Set("example", js.FuncOf(example))
	<-make(chan struct{}) // keep running
}
