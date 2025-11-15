# Fill/polygon

[Full code](https://github.com/visendi-labs/mapbox-gl-gojs/tree/main/docs/example6.1) (run with `go run main.go` from the `golang` folder)

[](example6.1/wasm/index.html ':include :type=iframe width=100% height=500px')


`index.html` Parsed by Go's [html/template](https://pkg.go.dev/html/template) or another templating tool

[filename](/example6.1/golang/index.html ':include :type=code')



`main.go` 

[filename](/example6.1/common/common.go ':include :type=code :fragment=demo')

`main.go`
[filename](/example6.1/golang/main.go ':include :type=code :fragment=demo')


# Comments

`GenerateId: true` is needed here, since the hover effect require IDs (not present in input data geojson).


