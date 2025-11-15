package common

import (
	_ "embed"
	"math/rand"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]
func Example(token string) string {
	features := []geojson.Feature{}
	for range 20 {
		p := orb.Point{rand.Float64() * 50, rand.Float64() * 50}
		s := rand.Float64() * 20
		features = append(
			features,
			*geojson.NewFeature(orb.Point{rand.Float64() * 50, rand.Float64() * 50}),
			*geojson.NewFeature(orb.LineString{
				orb.Point{rand.Float64() * 50, rand.Float64() * 50},
				orb.Point{rand.Float64() * 50, rand.Float64() * 50},
			}),
			*geojson.NewFeature(orb.Polygon{orb.Ring{
				p,
				orb.Point{p[0], p[1] + s},
				orb.Point{p[0] + s, p[1] + s},
				orb.Point{p[0] + s, p[1]},
				p,
			}}),
		)
	}
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token, Zoom: 2, Center: orb.Point{30, 30}}),
		mb.NewMapboxDraw(mb.MapboxDrawConfig{}),
		mb.NewMapOnLoad(mb.NewMapboxDrawAddFeatures(features...)),
	).MustRenderDefault()
}

/// [demo]
