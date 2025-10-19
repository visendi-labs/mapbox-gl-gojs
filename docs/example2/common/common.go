package common

import (
	"math/rand/v2"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

func Example(token string) string {
	/// [demo]
	lines := geojson.NewFeatureCollection()
	for i := 0; i < 20; i++ {
		line := orb.LineString{}
		for j := 0; j < 8; j++ {
			line = append(line, orb.Point{-45 + rand.Float64()*90, -45 + rand.Float64()*90})
		}
		lines.Append(geojson.NewFeature(line))
	}
	m := mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{
			Container:   "map",
			AccessToken: token,
		}),
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
	)
	/// [demo]
	s := m.MustRender(mapboxglgojs.RenderConfig{})
	return string(s)
}
