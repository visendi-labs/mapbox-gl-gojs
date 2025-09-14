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

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type RenderConfigOption int

const (
	UseCustomJsonMarshal RenderConfigOption = iota
)

type RenderConfig struct {
	UseCustomJsonMarshal bool
}

type EnclosedSnippetCollectionRenderable func(RenderConfig) *EnclosedSnippetCollection

type EnclosedSnippetCollection struct {
	Template string
	Data     map[string]string
	Children []EnclosedSnippetCollectionRenderable
}

func NewEnclosedSnippetCollection(template string, data map[string]string, c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	a := EnclosedSnippetCollection{
		Template: template,
		Data:     data,
		Children: c,
	}
	return EnclosedSnippetCollectionRenderable(func(rc RenderConfig) *EnclosedSnippetCollection {
		return &a
	})
}

func NewMapScript(c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	// Add script tag type as input? Like async, module or things like that
	return NewEnclosedSnippetCollection(`<script>{{.Children }}</script>`, map[string]string{}, c...)
}

func NewMap(mc MapConfig) EnclosedSnippetCollectionRenderable {
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
	return NewEnclosedSnippetCollection(`map.on("{{.Data.event}}", "{{.Data.layer}}", () => { {{- .Children -}} });`, map[string]string{
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
	return NewEnclosedSnippetCollection(`
	map.addImage("{{.Data.name}}",{width:{{.Data.width}},height:{{.Data.height}},data:Uint8Array.fromBase64("{{.Data.img}}") });
	`, map[string]string{
		"name":   name,
		"img":    imgBase64,
		"width":  strconv.Itoa(width),
		"height": strconv.Itoa(height),
	})
}

func NewMapOnEvent(event string, c ...EnclosedSnippetCollectionRenderable) EnclosedSnippetCollectionRenderable {
	return NewEnclosedSnippetCollection(`map.on("{{.Data.data}}", () => { {{- .Children -}} });`, map[string]string{"data": event}, c...)
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
		Data     map[string]string
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
	Type string `json:"type,omitempty"`
	Data any    `json:"data,omitempty"`
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

type Map struct {
	Container   string    `json:"container,omitempty"`
	Style       string    `json:"style,omitempty"`
	Center      orb.Point `json:"center,omitempty"`
	Zoom        float64   `json:"zoom,omitempty"`
	AccessToken string    `json:"accessToken,omitempty"`
}
