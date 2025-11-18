package common

import (
	"embed"
	_ "image/png" // Don't forget this for PNG
	"io/fs"
	"strconv"

	"github.com/paulmach/orb"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

/// ### [demo]

//go:embed weather
var weather embed.FS

var dir []fs.DirEntry

func ReadFiles() {
	dir, _ = weather.ReadDir("weather")
}

func UpdateUrl(t string) string {
	tInt, _ := strconv.Atoi(t)
	return mb.NewScript(mb.NewMapSourceUpdateImageUrl("weather", "../common/weather/"+dir[tInt].Name())).MustRenderDefault()
}

func Example(token string) string {
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", Zoom: 3.7, Center: orb.Point{16, 61.2}, AccessToken: token}),
		mb.NewMapOnLoad(
			mb.NewMapAddLayer(mb.MapLayer{
				Id: "weather", Type: "raster",
				Source: mb.MapSource{Type: "image", Url: "../common/weather/" + dir[0].Name(), Coordinates: []orb.Point{
					{9.3191647, 69.4197068}, {29.7990632, 69.4197068}, {29.7990632, 53.869605}, {9.3191647, 53.869605},
				}},
				Paint: mb.MapLayerPaint{RasterFadeDuration: 0, RasterOpacity: 0.7},
			}),
		),
	).MustRenderDefault()
}

/// [demo]
