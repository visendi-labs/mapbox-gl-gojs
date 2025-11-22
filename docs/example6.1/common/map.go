package common

import (
	_ "embed"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]

// https://github.com/perliedman/svenska-landskap
//
//go:embed landskap.geojson
var zones []byte

func Example(token string) string {
	lands, _ := geojson.UnmarshalFeatureCollection(zones)
	layers := []mb.EnclosedSnippetCollectionRenderable{}
	for _, z := range lands.Features {
		id := uuid.NewString()
		sw, ne := orb.Point{math.MaxFloat64, math.MaxFloat64}, orb.Point{0.0, 0.0}
		for _, p := range z.Geometry.(orb.MultiPolygon) {
			for _, r := range p {
				for _, c := range r {
					if c.Lat() < sw.Lat() {
						sw[1] = c.Lat()
					}
					if c.Lat() > ne.Lat() {
						ne[1] = c.Lat()
					}
					if c.Lon() < sw.Lon() {
						sw[0] = c.Lon()
					}
					if c.Lon() > ne.Lon() {
						ne[0] = c.Lon()
					}
				}
			}
		}
		layers = append(layers,
			mb.NewMapOnEventLayer("click", id, mb.NewMapFitBounds(sw, ne, mb.FitBoundsOptions{Padding: 50})),
			mb.NewMapAddLayer(mb.MapLayer{
				Id: id, Type: "fill",
				Paint: mb.MapLayerPaint{
					FillColor:   fmt.Sprintf("rgb(%d,%d,%d)", rand.IntN(255), rand.IntN(255), rand.IntN(255)),
					FillOpacity: []any{"case", []any{"boolean", []any{"feature-state", "hover"}, false}, 0.3, 0.6},
				},
				Source: mb.MapSource{Type: "geojson", Data: *z, GenerateId: true},
			}), mb.NewMapOnEventLayerPairFeatureState("mouseover", "mouseout", id, id, "hover", "true", "false"))
	}
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token, Zoom: 3, Center: orb.Point{17, 61}}), mb.NewMapOnLoad(layers...),
	).MustRenderDefault()
}

/// [demo]
