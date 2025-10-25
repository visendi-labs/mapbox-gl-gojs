package common

import (
	"bytes"
	"image"
	"math/rand/v2"

	_ "embed"
	_ "image/png"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]

//go:embed icon.png
var imgData []byte

func Example(token string) string {
	img, _, _ := image.Decode(bytes.NewReader(imgData))

	points1, points2 := geojson.NewFeatureCollection(), geojson.NewFeatureCollection()
	for range 100 {
		points1.Append(geojson.NewFeature(orb.Point{rand.Float64() * 40, rand.Float64() * 30}))
		points2.Append(geojson.NewFeature(orb.Point{rand.Float64() * 40, rand.Float64() * -30}))
	}
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddImageCircle("myCircleImage", 5, 2),
			mapboxglgojs.NewMapAddImage("food", img),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id:     "points1",
				Type:   "symbol",
				Source: mapboxglgojs.MapSource{Type: "geojson", Data: *points1},
				Layout: mapboxglgojs.MapLayout{
					IconImage:        "myCircleImage",
					IconAllowOverlap: true,
				},
			}),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id:     "points2",
				Type:   "symbol",
				Source: mapboxglgojs.MapSource{Type: "geojson", Data: *points2},
				Layout: mapboxglgojs.MapLayout{
					IconImage:        "food",
					IconAllowOverlap: true,
				},
			}),
		),
	).MustRenderDefault()
}

/// [demo]
