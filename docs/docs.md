# Docs 
## Usage
Generate Mapbox HTML/JS from Go. Some use cases:
- Output a HTML file to vizualize your Golang geo work
- Serve HTML/JS from a webserver to support geo workflow for your web app users

See Examples for more details.

## Mapbox docs

- Mapbox-GL-JS Github: https://github.com/mapbox/mapbox-gl-js
- Mapbox-GL-JS Docs: https://docs.mapbox.com/mapbox-gl-js/api/

## Go templating
The Go core templating package `html/template` is being used. See https://pkg.go.dev/html/template. 

## Go geo
`orb` and `geojson` are used for geo works in Mapbox-GL-GOJS.
- https://github.com/paulmach/orb
- https://github.com/paulmach/orb/geojson

## Mapbox-GL-JS supported operations
Mapbox-GL-GOJS supports several of the key Mapbox-GL-JS operations. See the complete list of supported underlying operations below. Several of the operations below have more than one wrapper/helper, i.e. underlying `map.on` is used in both `func NewMapOnEventLayer(...)` and `func NewMapOnEvent(...)`. More to come. 

- `map.fitBounds()`
- `map.flyTo()`
- `map.on()`
- `map.getCanvas()`
- `map.setConfig()`
- `map.addImage()`
- `map.addSource()`
    - `<source>.setData()`
    - `<source>.updateImage()`
- `map.getSource()`
- `map.removeSource()`
- `map.removeLayer()`
- `map.setLayoutProperty()`
- `map.setFeatureState()`
- `map.getLayer()`
- `map.addLayer()`
- `map.addControl()`


- `new mapboxgl.Map()`
- `new mapboxgl.Popup()`
- `new MapboxDraw()`

## Notes on HTMX (https://htmx.org)
HTMX can be used to together with Mapbox-GL-GOJS to tie interactivity to Mapbox content/config changes from the server. Mapbox-GL-GOJS even comes with some built-in HTMX support, so that you can (serverside) generate a map setup that ties map events to HTMX calls (that will make further request to your server).

#
#
#
#
