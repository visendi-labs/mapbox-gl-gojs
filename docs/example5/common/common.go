package common

import (
	"math/rand"

	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

// / ### [demo]
func Example(token string) string {
	return mapboxglgojs.NewGroup(
		mapboxglgojs.NewMap(mapboxglgojs.Map{Container: "map", AccessToken: token}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewPopup(geojson.Point{rand.Float64() * 30, rand.Float64() * 30}, "<h3>HTML</h3><br><button>Button</button>"),
			mapboxglgojs.NewPopup(geojson.Point{rand.Float64() * -30, rand.Float64() * -30}, "<h3>HTML</h3><br><button>Button</button>"),
		),
	).MustRenderDefault()
}

/// [demo]
