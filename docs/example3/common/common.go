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
	for i := 0; i < 20; i++ {
		line := orb.LineString{}
		for j := 0; j < 2; j++ {
			line = append(line, orb.Point{-45 + rand.Float64()*90, -45 + rand.Float64()*90})
		}
		lines.Append(geojson.NewFeature(line))
	}
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{
			Container:   "map",
			AccessToken: token,
		}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddSource("mySource", mapboxglgojs.MapSource{
				Type:       "geojson",
				Data:       lines,
				GenerateId: true,
			}),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id:     "myLayer",
				Type:   "line",
				Source: "mySource",
				Paint: mapboxglgojs.MapLayerPaint{
					LineColor: "#44f",
					LineWidth: []any{
						"case",
						[]any{"boolean", []any{"feature-state", "hover"}, false},
						12, 6,
					},
				},
			}),
			mapboxglgojs.NewMapOnEventLayerPairFeatureState("mouseover", "mouseout", "myLayer", "mySource", "hover", "true", "false"),
		),
	).MustRenderDefault()
}

/// [demo]
