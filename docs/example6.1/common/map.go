package common

import (
	_ "embed"
	"fmt"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mb "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]

// https://github.com/perliedman/svenska-landskap (changed to all being MyltiPolygons)
//
//go:embed landskap.geojson
var zones []byte

func Example(token string) string {
	lands, _ := geojson.UnmarshalFeatureCollection(zones)
	layers := []mb.EnclosedSnippetCollectionRenderable{}
	for _, l := range lands.Features {
		id := uuid.NewString()
		b := l.Geometry.(orb.MultiPolygon).Bound()
		layers = append(layers,
			mb.NewMapOnEventLayer("click", id, mb.NewMapFitBounds(b.Min, b.Max, mb.FitBoundsOptions{Padding: 50})),
			mb.NewMapAddLayer(mb.MapLayer{
				Id: id, Type: "fill",
				Paint: mb.MapLayerPaint{
					FillColor:   fmt.Sprintf("rgb(%d,%d,%d)", rand.IntN(255), rand.IntN(255), rand.IntN(255)),
					FillOpacity: []any{"case", []any{"boolean", []any{"feature-state", "hover"}, false}, 0.3, 0.6},
				},
				Source: mb.MapSource{Type: "geojson", Data: *l, GenerateId: true},
			}), mb.NewMapOnEventLayerPairFeatureState("mouseover", "mouseout", id, id, "hover", "true", "false"))
	}
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token, Zoom: 3, Center: orb.Point{17, 61}}), mb.NewMapOnLoad(layers...),
	).MustRenderDefault()
}

/// [demo]
