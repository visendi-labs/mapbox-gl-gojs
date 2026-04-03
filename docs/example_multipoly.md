# Fill/polygon

[Full code](https://github.com/visendi-labs/mapbox-gl-gojs/tree/main/docs/example_multipoly) (run with `go run main.go` from the `golang` folder)

[](example_multipoly/wasm/index.html ':include :type=iframe width=100% height=500px')


`index.html` Parsed by Go's [html/template](https://pkg.go.dev/html/template) or similar templating tool

[filename](/example_multipoly/golang/index.html ':include :type=code')



`map.go` 

[filename](/example_multipoly/common/map.go ':include :type=code :fragment=demo')

# Comments

- `GenerateId: true` is needed here in `map.go`, since the hover effect require feature IDs (not present in input data geojson).
- `mb.NewMapOnEventLayerPairFeatureState("mouseover", "mouseout", id, id, "hover"...` the double `id` here is because the source gets the same id as the layer when added inline.


