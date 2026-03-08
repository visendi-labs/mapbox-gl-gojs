package common

import (
	"bytes"
	"image"
	"math/rand/v2"

	_ "embed"
	"image/color"
	_ "image/png"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]

//go:embed icon.png
var takeawayFood []byte

func Example(token string) string {
	img, _, _ := image.Decode(bytes.NewReader(takeawayFood))

	p1, p2 := geojson.NewFeatureCollection(), geojson.NewFeatureCollection()
	for range 100 {
		p1.Append(geojson.NewFeature(orb.Point{rand.Float64() * 120, rand.Float64() * 80}))
		p2.Append(geojson.NewFeature(orb.Point{rand.Float64() * 120, rand.Float64() * 80}))
	}
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token, Config: mb.MapConfig{
			Basemap: mb.BasemapConfig{Theme: "monochrome", LightPreset: "dawn"},
		}}),
		mb.NewMapOnLoad(
			mb.NewMapAddImageCircle("circle", 5, 2, color.RGBA{10, 10, 10, 255}, color.RGBA{80, 150, 200, 255}),
			mb.NewMapAddImage("food", img),
			mb.NewMapAddLayer(mb.MapLayer{
				Id: "points1", Type: "symbol",
				Source: mb.MapSource{Type: "geojson", Data: *p1},
				Layout: mb.MapLayout{IconImage: "circle", IconAllowOverlap: true},
			}),
			mb.NewMapAddLayer(mb.MapLayer{
				Id: "points2", Type: "symbol",
				Source: mb.MapSource{Type: "geojson", Data: *p2},
				Layout: mb.MapLayout{IconImage: "food", IconAllowOverlap: true},
			}),
		),
	).MustRenderDefault()
}

/// [demo]
