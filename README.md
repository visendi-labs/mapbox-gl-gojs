# Mapbox GL JS Golang Wrapper
Generate Mapbox GL JS from Golang. For robust, scalable and fast mapping operations via serverside templates.

### Why
- Serverside rendering (SSR) + Mapbox don't go hand in hand. The use of Mapbox GL JS in a SSR web app forces critical map logic away from the server and down into JS code in the browser. It could also require the SSR webserver to support non-hypermedia endpoints, to e.g. return a GeoJSON data response, which goes against the SSR concept.

The aim of `mapboxglgojs` is to enable you to keep all map logic in the webserver, where it belongs.

### Details
Uses `"github.com/paulmach/orb"` & `"github.com/paulmach/orb/geojson"` for Golang GeoJSON data, which supports fast JSON marshaling via `"github.com/json-iterator/go"`. Uses `"html/template"` for templating.

### TODO
- HTMX

```go
package main

import (
    "fmt"
	"math/rand/v2" 
    "github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mbgojs "github.com/visendi-labs/mapbox-gl-gojs"
)
    
func main() {
    lines := geojson.NewFeatureCollection()

	for i := 0; i < 10; i++ {
		line := orb.LineString{}
		for i := 0; i < 10; i++ {
			line = append(line, orb.Point{-30 + rand.Float64()*60, -30 + rand.Float64()*60})
		}
		lines.Append(geojson.NewFeature(line))
	}

	map := mbgojs.NewMapScript(
		mbgojs.NewMap(mbgojs.Map{
			Container:   "map",
			AccessToken: "<MAPBOX_ACCESS_TOKEN>",
		}),
		mbgojs.NewConsoleLog("map"),
		mbgojs.NewMapOnLoad(
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id:   "layer",
				Type: "line",
				Source: mbgojs.MapSource{
					Type: "geojson",
					Data: *lines,
				},
			}),
        )
	)
    mapStr, err := map.Render()
	if err != nil {
		panic(err)
	}
    fmt.Println(mapStr)
}
```