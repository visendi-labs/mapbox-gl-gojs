package common

import (
	"math/rand"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]
func Example(token string) string {
	fc := geojson.NewFeatureCollection()
	for i := 0; i < 100; i++ {
		fc = fc.Append(&geojson.Feature{
			Type:       "Feature",
			Geometry:   orb.Point{-45.0 + rand.Float64()*90, -45.0 + rand.Float64()*90},
			Properties: geojson.Properties{"val": rand.Float64() * 1000},
		})
	}
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddSource("mySource", mapboxglgojs.MapSource{Type: "geojson", Data: fc, GenerateId: true}),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id:     "myLayer",
				Type:   "heatmap",
				Source: "mySource",
				Paint: mapboxglgojs.MapLayerPaint{
					HeatmapWeight: []any{
						"interpolate",
						[]string{"linear"},
						[]string{"get", "val"},
						0, 0, 6, 1,
					},
				},
			}),
		),
	).MustRenderDefault()
}

/// [demo]
