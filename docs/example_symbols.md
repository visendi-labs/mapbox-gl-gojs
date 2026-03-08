# Points/symbols

Mapbox-GL-GOJS has built in support to add circle and square images as symbol layers. You can also import your own. The images are turned into base64 strings and added as input to a Mapbox load image function.

[Full code](https://github.com/visendi-labs/mapbox-gl-gojs/tree/main/docs/example_symbols) (run with `go run main.go` from the `golang` folder)

[](example_symbols/wasm/index.html ':include :type=iframe width=100% height=500px')


`index.html` Parsed by Go's [html/template](https://pkg.go.dev/html/template) or another templating tool

[filename](/example_symbols/golang/index.html ':include :type=code')



`map.go` 

[filename](/example_symbols/common/map.go ':include :type=code :fragment=demo')



