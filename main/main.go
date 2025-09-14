package main

import (
	"fmt"

	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

func main() {
	s, err := mapboxglgojs.NewMapScript(
		mapboxglgojs.NewMap(mapboxglgojs.MapConfig{
			Container:   "123",
			AccessToken: "token",
		}),
		mapboxglgojs.NewMapOnLoad(
			mapboxglgojs.NewMapAddLayer(mapboxglgojs.MapLayer{
				Id:     "aaa",
				Type:   "",
				Source: nil,
				Layout: mapboxglgojs.MapLayout{},
				Paint:  "",
			}),
		),
	).Render(
		mapboxglgojs.RenderConfig{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
