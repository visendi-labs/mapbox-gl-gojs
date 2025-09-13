package mapboxglgojs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"strings"

	"github.com/paulmach/orb"
)

type EnclosedSnippetCollection struct {
	Children []EnclosedSnippetCollection
	Template string
	Data     any
}

func NewEnclosedSnippetCollection(template string, data any, c ...EnclosedSnippetCollection) EnclosedSnippetCollection {
	return EnclosedSnippetCollection{
		Children: c,
		Template: template,
		Data:     data,
	}
}

func NewMapScript(c ...EnclosedSnippetCollection) EnclosedSnippetCollection {
	// Add script tag type as input? Like async, module or things like that
	return NewEnclosedSnippetCollection(`<script>
	{{.Children}}
	</script>`, "", c...)
}

func NewMap(mc Map) EnclosedSnippetCollection {
	j, err := json.Marshal(mc)
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		`
		const map = new mapboxgl.Map({{.Data}});
		`,
		string(j),
	)
}

func NewMapOnEventLayer(event, layer string, c ...EnclosedSnippetCollection) EnclosedSnippetCollection {
	return NewEnclosedSnippetCollection("map.on(\"{{.Data.Event}}\", \"{{.Data.Layer}}\", () => { {{.Children}} });", struct {
		Event, Layer string
	}{
		Event: event,
		Layer: layer,
	}, c...)
}

func NewMapOnEventLayerCursor(event, layer, cursor string) EnclosedSnippetCollection {
	return NewMapOnEventLayer(event, layer, EnclosedSnippetCollection{
		Template: `map.getCanvas().style.cursor = "{{.Data}}"`,
		Data:     cursor,
	})
}

func NewConsoleWarn(log string) EnclosedSnippetCollection {
	return NewEnclosedSnippetCollection("console.warn({{.Data}});", log)
}

func NewConsoleLog(log string) EnclosedSnippetCollection {
	return NewEnclosedSnippetCollection("console.log({{.Data}});", log)
}

func NewMapAddImageCircle(name string, size int) EnclosedSnippetCollection {
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

func NewMapAddImageRectangle(name string, height, width int) EnclosedSnippetCollection {
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

func NewMapAddImage(name, imgBase64 string, width, height int) EnclosedSnippetCollection {
	return NewEnclosedSnippetCollection(`
	map.addImage("{{.Data.Name}}", {width:{{.Data.Width}},height:{{.Data.Height}},data:Uint8Array.fromBase64("{{.Data.Img}}") });
	`, struct {
		Name, Img     string
		Width, Height int
	}{
		Name:   name,
		Img:    imgBase64,
		Width:  width,
		Height: height,
	})
}

func NewMapOnEvent(event string, c ...EnclosedSnippetCollection) EnclosedSnippetCollection {
	return NewEnclosedSnippetCollection(`
	map.on("{{.Data}}", () => { {{.Children}} });
	`, event, c...)
}

func NewMapOnLoad(c ...EnclosedSnippetCollection) EnclosedSnippetCollection {
	return NewMapOnEvent("load", c...)
}

func (esc *EnclosedSnippetCollection) Render() (string, error) {
	t, err := template.New("page").Parse(esc.Template)
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	children := strings.Builder{}

	for _, sc := range esc.Children {
		scb, err := sc.Render()
		if err != nil {
			return "", err
		}
		fmt.Fprint(&children, scb)
	}
	if err := t.Execute(&b, struct {
		Children string
		Data     any
	}{
		Children: children.String(),
		Data:     esc.Data,
	}); err != nil {
		return "", err
	}
	return b.String(), nil
}

type MapLayout struct {
	LineJoin         string `json:"line-join,omitempty"`
	LineCap          string `json:"line-cap,omitempty"`
	IconImage        string `json:"icon-image,omitempty"`
	TextField        string `json:"text-field,omitempty"`
	IconAllowOverlap bool   `json:"icon-allow-overlap,omitempty"`
	IconOffset       any    `json:"icon-offset,omitempty"`
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

func NewMapAddSource(ms MapSource) EnclosedSnippetCollection {
	j, err := json.MarshalIndent(ms, "", "\t")
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		"map.addSource({{.Data}});",
		string(j),
	)
}

func NewMapAddLayer(ml MapLayer) EnclosedSnippetCollection {
	j, err := json.MarshalIndent(ml, "", "\t")
	if err != nil {
		panic(err)
	}
	return NewEnclosedSnippetCollection(
		"map.addLayer({{.Data}});",
		string(j),
	)
}

type Map struct {
	Container   string    `json:"container,omitempty"`
	Style       string    `json:"style,omitempty"`
	Center      orb.Point `json:"center,omitempty"`
	Zoom        float64   `json:"zoom,omitempty"`
	AccessToken string    `json:"accessToken,omitempty"`
}
