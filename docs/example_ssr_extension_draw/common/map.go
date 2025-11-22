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
			comment = fmt.Sprintf("%0.1f m<sup>2</sup>", geo.Area(f.Geometry))
		case geojson.TypePoint:
			comment = fmt.Sprintf("(%0.5f,%0.5f)", f.Point()[0], f.Point()[1])
		}
		html += fmt.Sprintf(`<li class="list-group-item d-flex justify-content-between align-items-start">
			<div class="fw-bold">%s</div><span class="badge text-bg-primary rounded-pill">%s</span>
		</li>`, f.Geometry.GeoJSONType(), comment)
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
