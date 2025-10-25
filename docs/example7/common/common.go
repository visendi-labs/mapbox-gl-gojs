package common

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"html"
	"math/rand/v2"
	"strconv"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

var lines = geojson.NewFeatureCollection()

// / ### [demo]
func CreateLines() {
	for i := 0; i < 1000; i++ {
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

func Filter(distance string) string {
	d, _ := strconv.ParseFloat(distance, 64)
	filteredLines := geojson.NewFeatureCollection()
	for _, l := range lines.Features {
		points := l.Geometry.(orb.LineString)
		if geo.DistanceHaversine(points[0].Point(), points[1].Point()) <= d {
			filteredLines.Append(l)
		}
	}
	script := mapboxglgojs.NewScript(mapboxglgojs.NewMapSourceSetData("mySource", filteredLines)).MustRenderDefault()

	randomId := uuid.NewString()
	buf := bytes.Buffer{}
	w := gzip.NewWriter(&buf)
	w.Write([]byte(script))
	w.Close()
	respStatus := fmt.Sprintf(`
		<li hx-swap-oob="afterbegin:#list-group" class="list-group-item">
			<div class="d-flex justify-content-between">
				<div>/filter?distance=%s</div>
				<button type="button" class="btn btn-sm btn-outline-success" disabled>200 OK</button>
			</div>
			<p class="d-inline-flex gap-1">
				<a data-bs-toggle="collapse" href="#%s" role="button" aria-expanded="false">Response (%.3f Kb, %.3f Kb gzip) ></a>
			</p>
			<div class="collapse m-1" id="%s">
				<div style="font-size:0.7rem" class="card font-monospace text-muted card-body">%s</div>
			</div>
		</li>`,
		distance, randomId, float64(len(script))/1024.0, float64(len(buf.Bytes()))/1024.0, randomId, html.EscapeString(script))
	return script + respStatus
}

/// [demo]
