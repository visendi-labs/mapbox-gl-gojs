package common

import (
	"math/rand/v2"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

var lines = geojson.NewFeatureCollection()

// / ### [demo]
func CreateLines() {
	for i := 0; i < 500; i++ {
		line := orb.LineString{orb.Point{rand.Float64() * 90, rand.Float64() * 90}, orb.Point{rand.Float64() * 90, rand.Float64() * 90}}
		lines.Append(geojson.NewFeature(line))
	}
}

func Example(token string) string {
	CreateLines()
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddSource("mySource", mapboxglgojs.MapSource{
				Type: "geojson", Data: lines,
			}),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{Id: "myLayer", Type: "line", Source: "mySource"}),
		),
	).MustRenderDefault()
}

func Filter(distance float64) string {
	filteredLines := geojson.NewFeatureCollection()
	for _, l := range lines.Features {
		points := l.Geometry.(orb.LineString)
		if geo.DistanceHaversine(points[0].Point(), points[1].Point()) <= distance {
			filteredLines.Append(l)
		}
	}
	return mapboxglgojs.NewScript(
		mapboxglgojs.NewMapSourceSetData("mySource", filteredLines),
	).MustRenderDefault()
}

/// [demo]
