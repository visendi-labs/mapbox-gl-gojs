# Mapbox GL GOJS
Interactive map and geo visualization for Golang. Generate [Mapbox](https://github.com/mapbox/mapbox-gl-js) HTML.

## Docs + Examples https://visendi-labs.github.io/mapbox-gl-gojs/#/

#### Usage - Generate Mapbox HTML/JS from Go 

```go
mapbox := mbgojs.NewScript(
	mbgojs.NewMap(mbgojs.Map{
		Container:   "map",
		AccessToken: "<MAPBOX_ACCESS_TOKEN>",
		Config:      mbgojs.MapConfig{Basemap: mbgojs.BasemapConfig{Theme: "faded"}},
	}),
)

mapbox.MustRenderDefault() // <script>const map = new mapboxgl.Map({container: "map", ...});</script>
```

```go
points := geojson.NewFeatureCollection()
points.Append(geojson.NewFeature(orb.Point{rand.Float64() * 50, rand.Float64() * 50}))

event := mbgojs.NewMapOnLoad(
	mbgojs.NewMapAddLayer(mbgojs.MapLayer{
		Id: "layer", Type: "line",
		Source: mbgojs.MapSource{Type: "geojson", Data: *points},
	}),
)

event.MustRenderDefault() // map.on("load", (e) => {map.addLayer({ ... })});
```