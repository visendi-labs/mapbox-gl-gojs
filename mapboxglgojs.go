package mapboxglgojs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	htmltemplate "html/template"
	"math"
	"strconv"
	"strings"
	"text/template"

	"github.com/paulmach/orb/geojson"
)

type RenderConfigOption int

const (
	UseCustomJsonMarshal RenderConfigOption = iota
)

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

func NewScript(c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	// Add script tag type as input? Like async, module or things like that
	return NewEnclosedSnippetCollection(`<script>{{.Children }}</script>`, map[string]string{}, c...)
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

func NewConsoleWarn(log string) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection("console.warn({{.Data.data}});", map[string]string{"data": log})
}

func NewConsoleLog(log string) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection("console.log({{.Data.data}});", map[string]string{"data": log})
}

func NewMapAddImageCircle(name string, size int) EnclosedSnippetCollectionRenderable {
	circleImg := []byte{}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			c := math.Pow(float64(i)-float64(size)/2, 2) + math.Pow(float64(j)-float64(size)/2, 2)
			r := math.Pow(float64(size)/2, 2)
			if c < r {
				circleImg = append(circleImg, 60, 150, 200, 255)
			} else if c > r-20 && c < r+20 {
				circleImg = append(circleImg, 30, 30, 30, 255)
			} else {
				circleImg = append(circleImg, 0, 0, 0, 0)
			}
		}
	}
	return NewMapAddImage(name, base64.StdEncoding.EncodeToString(circleImg), size, size)
}

func NewMapAddImageRectangle(name string, height, width int) EnclosedSnippetCollectionRenderable {
	squareImg := []byte{}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if j == 0 || i == 0 || j == width-1 || i == height-1 {
				squareImg = append(squareImg, 30, 30, 30, 255)
			} else {
				squareImg = append(squareImg, 60, 200, 150, 255)
			}
		}
	}
	return NewMapAddImage(name, base64.StdEncoding.EncodeToString(squareImg), height, width)
}

func NewMapAddImage(name, imgBase64 string, width, height int) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(
		`map.addImage("{{.Data.name}}",{width:{{.Data.width}},height:{{.Data.height}},data:Uint8Array.fromBase64("{{.Data.img}}") });`,
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
	Id     string    `json:"id,omitempty"`
	Type   string    `json:"type,omitempty"`
	Source any       `json:"source,omitempty"`
	Layout MapLayout `json:"layout,omitempty"`
	Paint  string    `json:"paint,omitempty"`
}

type MapSource struct {
	Type       string `json:"type,omitempty"`
	Data       any    `json:"data,omitempty"`
	GenerateId bool   `json:"generateId,omitempty"`
}

func NewMapAddSourceFeatureCollection(fc geojson.FeatureCollection) EnclosedSnippetCollectionRenderable {
	return NewMapAddSource(fc)
}

func NewMapAddSource(ms any) EnclosedSnippetCollectionRenderable {
	j, err := json.Marshal(ms) // Indent(ms, "", "\t")
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		"map.addSource({{.Data.data}});",
		map[string]string{"data": string(j)},
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

// NOTE: This one needs to be inside an event, so that it has access to the event variable e
func NewMapOnEventLayerHtmxAjaxEventData(event, layer, verb, path string) EnclosedSnippetCollectionRenderable {
	return NewMapOnEventLayer(event, layer, NewHtmxAjax(HtmxAjax{
		Path: path,
		Verb: verb,
		Context: HtmxAjaxContext{
			Values: map[string]string{
				"eventData": "ok",
			},
		},
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

func NewMapSourceSetData(sourceId string, data any) {}

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
