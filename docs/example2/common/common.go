package common

import (
	"math/rand/v2"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]
func Example(token string) string {
	lines := geojson.NewFeatureCollection()
	for range 20 {
		line := orb.LineString{}
		for range 8 {
			line = append(line, orb.Point{-45 + rand.Float64()*90, -45 + rand.Float64()*90})
		}
		lines.Append(geojson.NewFeature(line))
	}
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddSource("mySource", mapboxglgojs.MapSource{
				Type: "geojson",
				Data: lines,
			}),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id:     "myLayer",
				Type:   "line",
				Source: "mySource",
			}),
		),
	).MustRenderDefault()
}

/// [demo]
