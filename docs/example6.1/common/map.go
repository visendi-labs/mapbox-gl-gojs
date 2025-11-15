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

// https://github.com/perliedman/svenska-landskap
//
//go:embed landskap.geojson
var zones []byte

func Example(token string) string {
	zones, _ := geojson.UnmarshalFeatureCollection(zones)
	return mb.NewGroup(
		mb.NewMap(mb.Map{Container: "map", AccessToken: token, Pitch: 44, Zoom: 4, Center: orb.Point{17, 61}}),
		mb.NewMapOnLoad(func() (layers []mb.EnclosedSnippetCollectionRenderable) {
			for _, z := range zones.Features {
				id := uuid.NewString()
				layers = append(layers, mb.NewMapAddLayer(mb.MapLayer{
					Id:   id,
					Type: "fill",
					Paint: mb.MapLayerPaint{
						FillColor:   fmt.Sprintf("rgb(%d,%d,%d)", rand.IntN(255), rand.IntN(255), rand.IntN(255)),
						FillOpacity: []any{"case", []any{"boolean", []any{"feature-state", "hover"}, false}, 0.5, 0.8},
					},
					Source: mb.MapSource{Type: "geojson", Data: *z, GenerateId: true},
				}), mb.NewMapOnEventLayerPairFeatureState("mouseover", "mouseout", id, id, "hover", "true", "false"))
			}
			return layers
		}()...),
	).MustRenderDefault()
}

/// [demo]
