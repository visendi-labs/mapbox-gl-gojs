//go:build js && wasm

package main

import (
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_colors/common"
)

var Token = "<token>"

func main() {
	js.Global().Set("example", js.FuncOf(func(this js.Value, args []js.Value) any {
		return common.Example(Token)
	}))
	js.Global().Set("thickness", js.FuncOf(func(this js.Value, args []js.Value) any {
		return common.PaintProperty("line-width", args[0].Int())
	}))
	js.Global().Set("color", js.FuncOf(func(this js.Value, args []js.Value) any {
		c, _ := url.QueryUnescape(args[0].String())
		return common.PaintProperty("line-color", fmt.Sprintf("\"%s\"", c))
	}))
	<-make(chan struct{})
}
