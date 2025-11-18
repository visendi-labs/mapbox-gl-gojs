# Raster images 
[Full code](https://github.com/visendi-labs/mapbox-gl-gojs/tree/main/docs/example_htmx_raster) (run with `go run main.go` from the `golang` folder)

This page was used to get the weird weather coords to real ones (lmao look at the utm source of the url lmao): https://epsg.io/transform?utm_source=chatgpt.com#s_srs=3006&t_srs=4326&x=1075693.0000000&y=7771252.0000000

[](example_htmx_raster/wasm/index.html ':include :type=iframe width=100% height=500px')

`index.html` Parsed by Go's [html/template](https://pkg.go.dev/html/template) or another templating tool

[filename](/example_htmx_raster/golang/index.html ':include :type=code :fragment=demo')

`map.go`

[ ](/example_htmx_raster/common/map.go ':include :type=code :fragment=demo')


# Comments
- Shout-out SMHI
- Shout-out the storm "Babet" of 2023