package common

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]

//go:embed features.json
var island []byte

func AddFeatures(f string) (html string) {
	feats := []geojson.Feature{}
	json.Unmarshal([]byte(f), &feats)
	for _, f := range feats {
		comment := ""
		switch f.Geometry.GeoJSONType() {
		case geojson.TypeLineString:
			comment = fmt.Sprintf("%0.1fm", geo.LengthHaversine(f.Geometry))
		case geojson.TypePolygon:
			comment = fmt.Sprintf("%0.1f m<sup>2</sup> / %0.1fm", geo.Area(f.Geometry), geo.LengthHaversine(f.Geometry))
		case geojson.TypePoint:
			comment = fmt.Sprintf("(%0.5f,%0.5f)", f.Point()[0], f.Point()[1])
		}
		html += fmt.Sprintf(`<li class="list-row flex justify-between">
			<div>
				<div>%s</div>
				<div class="font-semibold text-xs opacity-60">%s</div>
			</div>
			<button class="btn btn-square btn-ghost" onclick="%s">
				<svg class="size-[0.8rem]" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-zoom-in-icon lucide-zoom-in"><circle cx="11" cy="11" r="8"/><line x1="21" x2="16.65" y1="21" y2="16.65"/><line x1="11" x2="11" y1="8" y2="14"/><line x1="8" x2="14" y1="11" y2="11"/></svg>
			</button>
		</li>`, f.Geometry.GeoJSONType(), comment, mb.NewMapFitBounds(f.Geometry.Bound().Min, f.Geometry.Bound().Max, mb.FitBoundsOptions{}).MustRenderDefault())
	}
	return html
}

func Example(token string) string {
	feats := []geojson.Feature{}
	json.Unmarshal(island, &feats)
	return mb.NewScript(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token, Zoom: 12.5, Center: orb.Point{18.71822, 59.30057}}),
		mb.NewMapboxDraw(mb.MapboxDrawConfig{}),
		mb.NewMapOnEvent("draw.create", mb.NewHtmxAjax(mb.HtmxAjax{
			Path: "/create", Verb: "POST", Context: mb.HtmxAjaxContext{
				Swap: "afterbegin", Target: "#features",
				Values: map[string]string{"features": "JSON.stringify(e.features)"},
			},
		})),
		mb.NewMapOnLoad(mb.NewMapboxDrawAddFeatures(feats...)),
	).MustRenderDefault() + AddFeatures(string(island))
}

/// [demo]
