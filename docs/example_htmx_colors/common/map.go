package common

import (
	"math/rand/v2"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

func PaintProperty(key string, value any) string {
	return mb.NewScript(
		mb.NewMapSetPaintProperty("lines", key, value),
	).MustRenderDefault()
}

func Example(token string) string {
	fc := geojson.NewFeatureCollection()
	for range 100 {
		line := orb.LineString{
			orb.Point{rand.Float64() * 50, rand.Float64() * 50},
			orb.Point{rand.Float64() * 50, rand.Float64() * 50},
		}
		fc.Append(geojson.NewFeature(line))
	}
	return mb.NewGroup(
		mb.NewMap(mb.Map{
			Container:   "map",
			AccessToken: token,
		}),
		mb.NewMapOnLoad(
			mb.NewMapAddLayer(mb.MapLayer{
				Id:   "lines",
				Type: "line",
				Source: mb.MapSource{
					Type: "geojson",
					Data: *fc,
				},
				Paint: mb.MapLayerPaint{
					LineWidth: 5,
				},
			}),
		),
	).MustRenderDefault()
}
