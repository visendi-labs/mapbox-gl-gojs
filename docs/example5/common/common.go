package common

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]
func Example(token string) string {
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{
			Container:   "map",
			AccessToken: token,
			Center:      orb.Point{18.07, 59.326},
			Zoom:        15.2,
			Bearing:     -10,
			Pitch:       65,
			Hash:        true,
			Config: mapboxglgojs.MapConfig{
				Basemap: mapboxglgojs.BasemapConfig{
					ShowTerrain:            true,
					LightPreset:            "night",
					ShowBuildingExtrusions: true,
					Show3DObjects:          true,
				},
			},
		}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewPopup(geojson.Point{18.071531, 59.326609}, 50, "<h3>Palace</h3>HTML<br><button>Button</button>"),
			mapboxglgojs.NewPopup(geojson.Point{18.0647542, 59.3245777}, 50, "<h3>Church</h3>HTML<br><button>Button</button>"),
		),
	).MustRenderDefault()
}

/// [demo]
