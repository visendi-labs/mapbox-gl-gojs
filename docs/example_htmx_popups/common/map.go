package common

import (
	"bytes"
	_ "embed"
	"fmt"
	_ "image/png" // Don't forget this for PNG
	"strconv"

	"image"
	"math/rand/v2"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

/// ### [demo]

//go:embed icon.png
var icon []byte

var points = geojson.NewFeatureCollection()

func Popup(index string) string {
	i, _ := strconv.Atoi(index)
	p := points.Features[i].Point()
	return mb.NewScript(mb.NewPopup(
		geojson.Point(p), mb.PopupConfig{CloseOnClick: true}, fmt.Sprintf("<b>SSR Popup %d</b>", i),
	), mb.NewMapFlyTo(mb.CameraOptions{Center: p}, mb.FlyToOptions{}),
	).MustRenderDefault()
}

func GeneratePoints() {
	for i := range 200 {
		points.Append(&geojson.Feature{ID: i, Geometry: orb.Point{rand.Float64() * 160, rand.Float64() * 80}})
	}
}

func Example(token string) string {
	heartIcon, _, _ := image.Decode(bytes.NewReader(icon))
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token}),
		mb.NewMapOnLoad(
			mb.NewMapAddImage("heart", heartIcon),
			mb.NewMapAddLayer(mb.MapLayer{
				Id: `points`, Type: "symbol",
				Source: mb.MapSource{Type: "geojson", Data: *points},
				Layout: mb.MapLayout{IconImage: "heart", IconAllowOverlap: true},
			}),
			mb.NewMapOnEventLayerCursor("mouseover", "points", "pointer"),
			mb.NewMapOnEventLayerCursor("mouseout", "points", ""),
			mb.NewMapOnEventLayer("click", "points", mb.NewHtmxAjax(
				mb.HtmxAjax{Path: "/popup", Verb: "GET",
					Context: mb.HtmxAjaxContext{
						// Call looks like this: /popup?eventType=click&featureId=184&lat=44.4321&layerId=points&layerType=symbol&lng=-0.2561129&sourceId=points&type=Feature&x=321.005&y=179.003
						Values: mb.HtmxAjaxContextEventValuesFull,
						Target: "body", Swap: "beforeend",
					},
				}),
			),
		),
	).MustRenderDefault()
}

/// [demo]
