package mapboxglgojs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	htmltemplate "html/template"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/paulmach/orb/geojson"
)

type RenderConfigOption int

const (
	UseCustomJsonMarshal RenderConfigOption = iota
)

// TODO: Add option to say what the map object should be called (mapboxgljs.Map())
type RenderConfig struct {
	UseCustomJsonMarshal bool
	KeepWhitespace       bool
}

type EnclosedSnippetCollectionRenderable func(RenderConfig) *EnclosedSnippetCollection

type EnclosedSnippetCollection struct {
	Template string
	Data     any
	Children []EnclosedSnippetCollectionRenderable
}

func NewEnclosedSnippetCollection(template string, data any, c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	a := EnclosedSnippetCollection{
		Template: template,
		Data:     data,
		Children: c,
	}
	return EnclosedSnippetCollectionRenderable(func(rc RenderConfig) *EnclosedSnippetCollection {
		return &a
	})
}

func NewTemplate(c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(`<template>{{.Children }}</template>`, map[string]string{}, c...)
}

func NewScript(c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	// Add script tag type as input? Like async, module or things like that
	return NewEnclosedSnippetCollection(`<script>{{.Children }}</script>`, map[string]string{}, c...)
}

func NewTimeout(t time.Duration, c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(`setTimeout(() => { {{.Children}} }, {{.Data.ms}})`, map[string]int64{
		"ms": t.Milliseconds(),
	}, c...)
}

func NewGroup(c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	// Add script tag type as input? Like async, module or things like that
	return NewEnclosedSnippetCollection(`{{.Children }}`, map[string]string{}, c...)
}

func NewMap(mc Map) EnclosedSnippetCollectionRenderable {
	j, err := json.Marshal(mc)
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		`const map = new mapboxgl.Map({{ .Data.data }});`,
		map[string]string{"data": string(j)},
	)
}

// TODO: This is probably way more sane than to have a HTMX endpoint do the same thing? Or?
// TODO: Could this be more modular? Or is it good to have a super specific function like this?
// TODO: Make the input here a struct instead?
func NewMapOnEventLayerPairFeatureState(event1, event2, layer, source, feature, event1Value, event2Value string) EnclosedSnippetCollectionRenderable {
	variable := "x" + strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(rand.Int()))), "=", "")
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		return NewEnclosedSnippetCollection(
			`let {{.Data.variable}} = null; {{.Children}}`,
			map[string]string{"event1": event1, "event2": event2, "layer": layer, "variable": variable},
			NewMapOnEventLayer(event1, layer,
				NewMapSetFeatureState(source, "", "e.features[0].id", map[string]string{
					feature: event1Value,
				}),
				NewEnclosedSnippetCollection("{{.Data.variable}} = e.features[0].id;", map[string]string{
					"variable": variable,
				}),
			),
			NewMapOnEventLayer(event2, layer,
				NewMapSetFeatureState(source, "", variable, map[string]string{
					feature: event2Value,
				}),
				NewEnclosedSnippetCollection("{{.Data.variable}} = null;", map[string]string{
					"variable": variable,
				}),
			),
		)(rc)
	}
}

func NewMapFlyTo(co CameraOptions, fo FlyToOptions) EnclosedSnippetCollectionRenderable {
	cJson, err := json.Marshal(co)
	if err != nil {
		panic(err)
	}
	fJson, err := json.Marshal(fo)
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(`map.flyTo({...{{.Data.co}}, ...{{.Data.fo}}});`, map[string]string{
		"co": string(cJson),
		"fo": string(fJson),
	})
}

// map.on("event", layer, (e) => { ... });
func NewMapOnEventLayer(event, layer string, c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(`map.on("{{.Data.event}}", "{{.Data.layer}}", (e) => { {{- .Children -}} });`, map[string]string{
		"event": event,
		"layer": layer,
	}, c...)
}

func NewMapOnEventLayerCursor(event, layer, cursor string) EnclosedSnippetCollectionRenderable {
	return NewMapOnEventLayer(event, layer, EnclosedSnippetCollectionRenderable(func(rc RenderConfig) *EnclosedSnippetCollection {
		return &EnclosedSnippetCollection{
			Template: `map.getCanvas().style.cursor = "{{.Data.data}}"`,
			Data:     map[string]string{"data": cursor},
		}
	}))
}

// TODO add all options in Popup constructor
func NewPopup(lngLat geojson.Point, config PopupConfig, html string) EnclosedSnippetCollectionRenderable {
	coord, err := json.Marshal([2]float64{lngLat[0], lngLat[1]})
	if err != nil {
		panic(err)
	}
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	return NewEnclosedSnippetCollection(
		`(new mapboxgl.Popup( {{.Data.config}} )).setLngLat({{.Data.coord}}).setHTML("{{.Data.html}}").addTo(map);`,
		map[string]any{
			"coord":  string(coord),
			"config": string(jsonConfig),
			"html":   html,
		},
	)

}

func NewMapSetBasemapConfig(b BasemapConfig) EnclosedSnippetCollectionRenderable {
	basemapJson, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection("map.setConfig('basemap', {{.Data.data}});", map[string]string{"data": string(basemapJson)})
}

func NewConsoleWarn(log string) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection("console.warn({{.Data.data}});", map[string]string{"data": log})
}

func NewConsoleLog(log string) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection("console.log({{.Data.data}});", map[string]string{"data": log})
}

func NewMapAddImageCircle(name string, rad int, border float64, borderColor, circleColor color.RGBA) EnclosedSnippetCollectionRenderable {
	circleImg := []byte{}
	for i := 0; i < rad*2; i++ {
		for j := 0; j < rad*2; j++ {
			c := math.Pow(float64(i)-float64(rad), 2) + math.Pow(float64(j)-float64(rad), 2)
			r := math.Pow(float64(rad), 2)
			b := math.Pow(float64(rad)-border, 2)
			if c > b && c < r {
				circleImg = append(circleImg, borderColor.R, borderColor.G, borderColor.B, borderColor.A)
			} else if c < r {
				circleImg = append(circleImg, circleColor.R, circleColor.G, circleColor.B, circleColor.A)
			} else {
				circleImg = append(circleImg, 0, 0, 0, 0)
			}
		}
	}
	return NewMapAddImageBase64(name, base64.StdEncoding.EncodeToString(circleImg), rad*2, rad*2)
}

func NewMapAddImageRectangle(name string, height, width, border int) EnclosedSnippetCollectionRenderable {
	squareImg := []byte{}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if j < border || i < border || j > width-border-1 || i > height-border-1 {
				squareImg = append(squareImg, 30, 30, 30, 255)
			} else {
				squareImg = append(squareImg, 60, 200, 150, 255)
			}
		}
	}
	return NewMapAddImageBase64(name, base64.StdEncoding.EncodeToString(squareImg), height, width)
}

func missingImageRGBA() image.RGBA {
	const w, h = 64, 64
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// light gray background
	bg := color.RGBA{200, 200, 200, 255}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	// dark X
	xcol := color.RGBA{80, 80, 80, 255}
	for i := 0; i < w; i++ {
		img.SetRGBA(i, i, xcol)     // \ diagonal
		img.SetRGBA(w-1-i, i, xcol) // / diagonal
	}

	return *img
}

func rawPixelsToBase64(img image.Image) (string, image.Rectangle) {
	switch v := img.(type) {
	case *image.RGBA:
		return base64.StdEncoding.EncodeToString(v.Pix), img.Bounds()
	case *image.NRGBA:
		return base64.StdEncoding.EncodeToString(v.Pix), img.Bounds()
	default:
		fmt.Println(fmt.Errorf("failed to encode image"))
	}
	i := missingImageRGBA()
	return base64.StdEncoding.EncodeToString(i.Pix), i.Bounds()
}

func NewMapAddImage(name string, image image.Image) EnclosedSnippetCollectionRenderable {
	b64, bounds := rawPixelsToBase64(image)
	return NewMapAddImageBase64(name, b64, bounds.Dy(), bounds.Dx())
}

func NewMapAddImageBase64(name, imgBase64 string, height, width int) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(
		`map.addImage("{{.Data.name}}",{width: {{.Data.width}}, height: {{.Data.height}}, data: Uint8Array.fromBase64("{{.Data.img}}")});`,
		map[string]string{
			"name":   name,
			"img":    imgBase64,
			"width":  strconv.Itoa(width),
			"height": strconv.Itoa(height),
		})
}

func NewMapOnEvent(event string, c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(`map.on("{{.Data.data}}", (e) => { {{- .Children -}} });`, map[string]string{"data": event}, c...)
}

func NewMapOnLoad(c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	return NewMapOnEvent("load", c...)
}

func (esc EnclosedSnippetCollectionRenderable) MustRenderDefault() string {
	s, err := esc.Render(RenderConfig{})
	if err != nil {
		panic(err)
	}
	return string(s)
}
func (esc EnclosedSnippetCollectionRenderable) MustRender(config RenderConfig) htmltemplate.JS {
	s, err := esc.Render(config)
	if err != nil {
		panic(err)
	}
	return s
}

func (esc EnclosedSnippetCollectionRenderable) Render(config RenderConfig) (htmltemplate.JS, error) {
	c := esc(config)
	t, err := template.New("page").Funcs(template.FuncMap{
		// "safeJS": func(input any) template.JS {
		// 	if s, ok := input.(string); ok {
		// 		return template.JS(s)
		// 	} else if s, ok := input.(template.JS); ok {
		// 		return s
		// 	}
		// 	return template.JS("")
		// },
	}).Parse(c.Template)
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	children := strings.Builder{}

	for _, sc := range c.Children {
		scb, err := sc.Render(config) // .Render()
		if err != nil {
			return "", err
		}
		fmt.Fprint(&children, scb)
	}
	if err := t.Execute(&b, struct {
		Children string
		Data     any // map[string]string
	}{
		Children: children.String(),
		Data:     c.Data,
	}); err != nil {
		return "", err
	}
	return htmltemplate.JS(b.String()), nil
}

type MapLayer struct {
	Id     string        `json:"id,omitempty"`
	Type   string        `json:"type,omitempty"`
	Source any           `json:"source,omitempty"`
	Layout MapLayout     `json:"layout,omitempty"`
	Paint  MapLayerPaint `json:"paint,omitempty"`
}

type MapSource struct {
	Type       string `json:"type,omitempty"`
	Data       any    `json:"data,omitempty"`
	GenerateId bool   `json:"generateId,omitempty"`
}

func NewMapAddSourceFeatureCollection(id string, fc geojson.FeatureCollection) EnclosedSnippetCollectionRenderable {
	return NewMapAddSource(id, MapSource{
		Type:       "geojson",
		Data:       fc,
		GenerateId: true,
	})
}

// TODO: Make this stricter/break out into more than one? Different types/formats of sources
func NewMapAddSource(id string, ms MapSource) EnclosedSnippetCollectionRenderable {
	j, err := json.Marshal(ms) // Indent(ms, "", "\t")
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		`map.addSource("{{.Data.id}}", {{.Data.data}});`,
		map[string]string{
			"id":   id,
			"data": string(j),
		},
	)
}

func NewMapSourceSetData(id string, d any) EnclosedSnippetCollectionRenderable {
	data, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		`map.getSource("{{.Data.id}}").setData({{.Data.data}});`,
		map[string]string{"id": id, "data": string(data)},
	)
}

func NewMapRemoveLayer(layerId string) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		return NewEnclosedSnippetCollection(
			`map.removeLayer("{{.Data.data}}");`,
			map[string]string{"data": layerId},
		)(rc)
	}
}

type HtmxAjax struct {
	Path    string
	Verb    string
	Context HtmxAjaxContext
}

var HtmxAjaxContextEventValuesFull map[string]string = map[string]string{
	"eventType": "e.type",
	"lat":       "e.lngLat.lat",
	"lng":       "e.lngLat.lng",
	"featureId": "e.features[0].id",
	"layerId":   "e.features[0].layer.id",
	"layerType": "e.features[0].layer.type",
	"sourceId":  "e.features[0].source",
	"type":      "e.features[0].type",
	"x":         "e.point.x",
	"y":         "e.point.y",
}

var HtmxAjaxContextEventValuesNoFeature map[string]string = map[string]string{
	"eventType": "e.type",
	"lat":       "e.lngLat.lat",
	"lng":       "e.lngLat.lng",
	"x":         "e.point.x",
	"y":         "e.point.y",
}

type HtmxAjaxContext struct {
	Headers map[string]string
	Values  map[string]string
	Swap    string
	Target  string
	Event   string
	Handler string
	Select  string
	Source  string
}

// TODO: Remove this!?!? Just better off using the NewMapOnEventLayer instead
// NOTE: This one needs to be inside an event, so that it has access to the event variable e
func NewMapOnEventLayerHtmxAjaxEventData(event, layer, verb, path string, h HtmxAjaxContext) EnclosedSnippetCollectionRenderable {
	return NewMapOnEventLayer(event, layer, NewHtmxAjax(HtmxAjax{
		Path: path, Verb: verb, Context: h,
	}))
}

// TODO: Be able to pass data down into children? Messy?
func NewMapOnEventLayerHtmxAjax(event, layer string, data HtmxAjax) EnclosedSnippetCollectionRenderable {
	return NewMapOnEventLayer(event, layer, NewHtmxAjax(data))
}

func NewMapOnEventHtmxAjax(event string, data HtmxAjax) EnclosedSnippetCollectionRenderable {
	return NewMapOnEvent(event, NewHtmxAjax(data))
}

func NewHtmxAjax(htmxAjax HtmxAjax) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		// j, err := json.Marshal(data) // TODO: How to marshal JS code (not JSON)?
		// if err != nil {
		// 	panic(err)
		// }
		s := `htmx.ajax(
			"{{.Data.Verb}}",
			"{{.Data.Path}}",
			{
				"values": {
					{{range $k, $v := .Data.Context.Values}}
						"{{$k}}": {{$v}},
					{{end}}
				},
				{{if .Data.Context.Swap}}"swap": "{{.Data.Context.Swap}}",{{end}}
				{{if .Data.Context.Target}}"target": "{{.Data.Context.Target}}",{{end}}
			}
		);`
		if !rc.KeepWhitespace {
			s = strings.ReplaceAll(s, "\t", "")
			s = strings.ReplaceAll(s, "\n", " ")
		}
		return NewEnclosedSnippetCollection(s, htmxAjax)(rc)
	}
} // TODO: Template a HTMX hx-vals

func NewHtmxAjaxRaw(verb, path, data string) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		return NewEnclosedSnippetCollection(
			`htmx.ajax("{{.Data.verb}}", "{{.Data.path}}")`, map[string]string{},
		)(rc)
	}
}

// TODO: For hover things - generate a UUID variable name to use for keeping track of "hovered line" id?
// TODO: This can't set other values than string atm
func NewMapSetLayoutProperty(layerId, propertry, value string) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		s := `map.setLayoutProperty("{{.Data.layerId}}", "{{.Data.property}}", "{{.Data.value}}");`
		return NewEnclosedSnippetCollection(s, map[string]string{
			"layerId":  layerId,
			"property": propertry,
			"value":    value,
		})(rc)
	}
}

// TODO: We really need to make it clear which values are not wrapped in quotes? General structs with different fields?
func NewMapSetFeatureState(source, sourceLayer, id string, features map[string]string) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		s := `map.setFeatureState({
			source: "{{.Data.source}}",
			{{- if .Data.sourceLayer }} sourceLayer: "{{.Data.sourceLayer}}",{{end -}}
			id: {{.Data.id}} 
		}, {
			{{range $k, $v := .Data.features}}
				"{{$k}}": {{$v}},
			{{end}}
		});`
		return NewEnclosedSnippetCollection(s, map[string]any{
			"source":      source,
			"sourceLayer": sourceLayer,
			"id":          id,
			"features":    features,
		})(rc)
	}
}

// Sould this just use a setData from source id?
func NewMapSourceSetDataFromLayer(layerId string, data any) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		j, err := json.Marshal(data) // Indent(ml, "", "\t")
		if err != nil {
			panic(err)
		}
		s := `
			(
				() => {
				 	let l = map.getLayer("{{.Data.layer}}");
					if(l) {
						let s = map.getSource(l.source);
						if(s) {
							s.setData({{.Data.data}});
						} else {
							console.warn("could not find source \"", l.source, "\" in layer \"{{.Data.layer}}\""); 
						}
					} else {
						console.warn("could not find layer \"{{.Data.layer}}\""); 
					}
				}
			)();`
		if !rc.KeepWhitespace {
			s = strings.ReplaceAll(s, "\t", "")
			s = strings.ReplaceAll(s, "\n", " ")
		}
		// TODO: Add in RenderConfig option to ignore console log warnings? I.e. if/else in template that skips on a bool input
		return NewEnclosedSnippetCollection( // TODO: Use embed here to read from real JS files?
			s,
			map[string]string{
				"layer": layerId,
				"data":  string(j),
			},
		)(rc)
	}
}

func NewMapAddLayer(ml MapLayer) EnclosedSnippetCollectionRenderable {
	return func(rc RenderConfig) *EnclosedSnippetCollection {
		// if rc.UseCustomJsonMarshal {

		// }
		j, err := json.Marshal(ml) // Indent(ml, "", "\t")
		if err != nil {
			panic(err)
		}
		return NewEnclosedSnippetCollection(
			"map.addLayer({{.Data.data}});",
			map[string]string{"data": string(j)},
		)(rc)
	}
}
