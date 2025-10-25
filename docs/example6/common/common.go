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
var foodIcon []byte

func Example(token string) string {
	img, _, _ := image.Decode(bytes.NewReader(foodIcon))

	p1, p2 := geojson.NewFeatureCollection(), geojson.NewFeatureCollection()
	for range 100 {
		p1.Append(geojson.NewFeature(orb.Point{rand.Float64() * 40, rand.Float64() * 30}))
		p2.Append(geojson.NewFeature(orb.Point{rand.Float64() * 40, rand.Float64() * -30}))
	}
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token, Zoom: 2.5, Config: mapboxglgojs.MapConfig{
			Basemap: mapboxglgojs.BasemapConfig{Theme: "monochrome", LightPreset: "dawn"},
		}}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddImageCircle("circle", 5, 2),
			mapboxglgojs.NewMapAddImage("food", img),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id: "points1", Type: "symbol",
				Source: mapboxglgojs.MapSource{Type: "geojson", Data: *p1},
				Layout: mapboxglgojs.MapLayout{IconImage: "circle", IconAllowOverlap: true},
			}),
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id: "points2", Type: "symbol",
				Source: mapboxglgojs.MapSource{Type: "geojson", Data: *p2},
				Layout: mapboxglgojs.MapLayout{IconImage: "food", IconAllowOverlap: true},
			}),
		),
	).MustRenderDefault()
}

/// [demo]
