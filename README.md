# Mapbox GL JS Golang Wrapper
Generate Mapbox GL JS HTML templates from Golang. For robust, scalable and fast mapping via serverside rendering.

### Why
- Serverside rendering (SSR) + Mapbox does not go hand in hand. Mapbox GL JS forces critical map logic away from the server and down into JS code in the browser. This can cause a lot of confusion, context switching, split focus and even duplication of logic between browser code and server. 

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