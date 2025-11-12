package common

import (
	_ "embed"
	"fmt"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]

// https://github.com/perliedman/svenska-landskap
//
//go:embed landskap.geojson
var zones []byte

func Example(token string) string {
	zones, _ := geojson.UnmarshalFeatureCollection(zones)
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token, Pitch: 44, Zoom: 4, Center: orb.Point{17, 61}}),
		mapboxglgojs.NewMapOnLoad(func() (layers []mapboxglgojs.EnclosedSnippetCollectionRenderable) {
			for _, z := range zones.Features {
				layers = append(layers, mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
					Id: uuid.NewString(), Type: "fill",
					Paint: mapboxglgojs.MapLayerPaint{
						FillColor:   fmt.Sprintf("rgb(%d,%d,%d)", rand.IntN(255), rand.IntN(255), rand.IntN(255)),
						FillOpacity: 0.6,
					},
					Source: mapboxglgojs.MapSource{Type: "geojson", Data: *z},
				}))
			}
			return layers
		}()...),
	).MustRenderDefault()
}

/// [demo]
