package common

import (
	"bytes"
	_ "embed"
	"fmt"
	_ "image/png"
	"strconv"

	"image"
	"math/rand/v2"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mbgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

/// ### [demo]

//go:embed icon.png
var icon []byte

var points = geojson.NewFeatureCollection()

func Popup(index string) string {
	i, _ := strconv.Atoi(index)
	return mbgojs.NewScript(
		mbgojs.NewPopup(
			geojson.Point(points.Features[i].Point()), mbgojs.PopupConfig{
				CloseOnClick: true,
				CloseOnMove:  true,
			}, fmt.Sprintf("<b>Serverside Popup</b><p>Id: %d</p>", i),
		),
	).MustRenderDefault()
}

func GeneratePoints() {
	for i := range 200 {
		points.Append(&geojson.Feature{
			ID:       i,
			Geometry: orb.Point{rand.Float64() * 80, rand.Float64() * 80},
		})
	}
}

func Example(token string) string {
	img, _, _ := image.Decode(bytes.NewReader(icon))
	return mbgojs.NewGroup(
		mbgojs.NewMap(mbgojs.Map{Container: "map", AccessToken: token}),
		mbgojs.NewMapOnLoad(
			mbgojs.NewMapAddImage("heart", img),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id: `points`, Type: "symbol",
				Source: mbgojs.MapSource{Type: "geojson", Data: *points},
				Layout: mbgojs.MapLayout{IconImage: "heart", IconAllowOverlap: true},
			}),
			mbgojs.NewMapOnEventLayerCursor("mouseover", "points", "pointer"),
			mbgojs.NewMapOnEventLayerCursor("mouseout", "points", ""),
			mbgojs.NewMapOnEventLayer("click", "points", mbgojs.NewHtmxAjax(
				mbgojs.HtmxAjax{
					Path: "/popup",
					Verb: "GET",
					Context: mbgojs.HtmxAjaxContext{
						Values: mbgojs.HtmxAjaxContextEventValuesFull,
						Target: "body",
						Swap:   "beforeend",
					},
				}),
			),
		),
	).MustRenderDefault()
}

/// [demo]
