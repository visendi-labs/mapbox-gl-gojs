package main

import (
	"fmt"
	"html/template"
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"

	"github.com/paulmach/orb/geojson"
	mbgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

func main() {
	lines := geojson.NewFeatureCollection()
	points1 := geojson.NewFeatureCollection()
	points2 := geojson.NewFeatureCollection()

	for i := 0; i < 1; i++ {
		line := orb.LineString{}
		for j := 0; j < 2; j++ {
			line = append(line, orb.Point{-30 + rand.Float64()*60, -30 + rand.Float64()*60})
		}
		lines.Append(geojson.NewFeature(line))
	}
	for i := 0; i < 10; i++ {
		points1.Append(geojson.NewFeature(orb.Point{rand.Float64() * 50, rand.Float64() * 50}))
		points2.Append(geojson.NewFeature(orb.Point{20 + rand.Float64()*50, 20 - rand.Float64()*50}))
	}

	mapbox := mbgojs.NewScript(
		mbgojs.NewMap(mbgojs.Map{
			Container:   "map",
			AccessToken: "<token>",
			Hash:        true,
			Config:      mbgojs.MapConfig{Basemap: mbgojs.BasemapConfig{Theme: "faded"}},
		}),
		mbgojs.NewMapOnLoad(
			mbgojs.NewMapAddImageRectangle("square", 20, 20, 4),
			mbgojs.NewMapAddImageCircle("circle", 10, 2),
			mbgojs.NewMapAddSource(
				"sourceId", mbgojs.MapSource{Type: "geojson", Data: lines, GenerateId: true},
			),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id: "layer", Type: "line", Source: "sourceId",
				Paint: mbgojs.MapLayerPaint{
					LineColor: "#116",
					LineWidth: []any{
						"case",
						[]any{"boolean", []any{"feature-state", "hover"}, false},
						6, 2,
					},
				},
			}),
			(func(sources ...string) mbgojs.EnclosedSnippetCollectionRenderable {
				r := []mbgojs.EnclosedSnippetCollectionRenderable{}
				for _, s := range sources {
					r = append(r,
						mbgojs.NewMapOnEventLayer("mouseover", "layer", mbgojs.NewHtmxAjax(
							mbgojs.HtmxAjax{
								Path: "/hover?hover=true",
								Verb: "GET",
								Context: mbgojs.HtmxAjaxContext{
									Values: mbgojs.HtmxAjaxContextEventValuesFull, Target: "#testtarget", Swap: "innerHTML",
								},
							}),
						),
						mbgojs.NewMapOnEventLayer("mouseout", "layer", mbgojs.NewHtmxAjax(
							mbgojs.HtmxAjax{
								Path:    fmt.Sprintf("/hover?hover=false&sourceId=%s", s),
								Verb:    "GET",
								Context: mbgojs.HtmxAjaxContext{Target: "#testtarget", Swap: "innerHTML"},
							},
						)),
					)
				}
				return mbgojs.NewGroup(r...)
			})("sourceId"),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id: "points1", Type: "symbol",
				Source: mbgojs.MapSource{Type: "geojson", Data: *points1, GenerateId: true},
				Layout: mbgojs.MapLayout{IconImage: "square", IconAllowOverlap: true},
			}),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id: `points2`, Type: "symbol",
				Source: mbgojs.MapSource{Type: "geojson", Data: *points2, GenerateId: true},
				Layout: mbgojs.MapLayout{IconImage: "circle", IconAllowOverlap: true},
			}),
			mbgojs.NewMapOnEventLayer("click", "points1", mbgojs.NewHtmxAjax(
				mbgojs.HtmxAjax{
					Path: "/click",
					Verb: "GET",
					Context: mbgojs.HtmxAjaxContext{
						Values: mbgojs.HtmxAjaxContextEventValuesFull,
						Target: "#testtarget",
						Swap:   "innerHTML",
					},
				}),
				mbgojs.NewMapSetLayoutProperty("points2", "visibility", "visible"),
			),
			mbgojs.NewMapOnEventLayer("click", "points2", mbgojs.NewMapSetLayoutProperty("points2", "visibility", "none")),
		),
		mbgojs.NewMapOnEventLayerCursor("mouseover", "points1", "pointer"),
		mbgojs.NewMapOnEventLayerCursor("mouseout", "points1", ""),
		mbgojs.NewMapOnEventLayerCursor("mouseover", "points2", "pointer"),
		mbgojs.NewMapOnEventLayerCursor("mouseout", "points2", ""),
	)

	s, err := mapbox.Render(mbgojs.RenderConfig{})
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	hoveredId := ""
	r.GET("/hover", func(c *gin.Context) {
		if f := c.Query("featureId"); f != "" {
			hoveredId = f
		}
		sourceId := c.Query("sourceId")
		hover := c.Query("hover")
		t, err := template.New("").Parse(string(mbgojs.NewScript(
			mbgojs.NewMapSetFeatureState(sourceId, "", hoveredId, map[string]string{
				"hover": hover,
			}),
		).MustRender(mbgojs.RenderConfig{})))
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}
		t.Execute(c.Writer, nil)
	})
	r.GET("/click", func(ctx *gin.Context) {
		feature := ctx.Query("feature")
		lng := ctx.Query("lng")
		lat := ctx.Query("lat")
		t, _ := template.New("page").Parse(`
		<div>{{.Feature}} -- {{.Lat}} , {{.Lng}}</div>
		`)
		if err := t.Execute(ctx.Writer, struct {
			Feature string
			Lat     string
			Lng     string
		}{
			Feature: feature,
			Lat:     lat,
			Lng:     lng,
		}); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	})
	r.GET("/", func(ctx *gin.Context) {
		t, _ := template.New("page").Parse(`<html><head>
			<script src='https://api.mapbox.com/mapbox-gl-js/v3.15.0/mapbox-gl.js'></script>
			<link href='https://api.mapbox.com/mapbox-gl-js/v3.15.0/mapbox-gl.css' rel='stylesheet' />
			<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.7/dist/htmx.min.js"></script>
		</head><body style="margin:0">
			<div id="testtarget"></div>
			<div id="map" style="width:100vw; height:100vh;"></div>
			{{.}}</body></html>`)

		if err := t.Execute(ctx.Writer, template.HTML(s)); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	})
	r.Run()
}

// curl --output data -H  "Accept-Encoding: gzip" localhost:8080
