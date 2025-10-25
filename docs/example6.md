# Points/symbols

Mapbox-GL-GOJS has built in support to add circle and square images as symbol layers. You can also import your own. The images are turned into base64 strings and added as input to a Mapbox load image function.

[Full code](http://example.com) (run with `go run main.go`  )

[](example6/wasm/index.html ':include :type=iframe width=100% height=500px')


`index.html` Parsed by Go's [html/template](https://pkg.go.dev/html/template) or another templating tool

[filename](/example6/golang/index.html ':include :type=code')



`main.go` 

[filename](/example6/common/common.go ':include :type=code :fragment=demo')

`main.go`
[filename](/example6/golang/main.go ':include :type=code')



