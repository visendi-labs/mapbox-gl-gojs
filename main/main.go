package main

import (
	"fmt"
	"html/template"
	"math/rand/v2"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/paulmach/orb"

	"github.com/paulmach/orb/geojson"
	mapboxglgojs "github.com/visendi-labs/mapbox-gl-gojs"
	mbgojs "github.com/visendi-labs/mapbox-gl-gojs"
)

func main() {
	fmt.Println("oks")
	c := jsoniter.Config{
		EscapeHTML:              true,
		SortMapKeys:             false,
		MarshalFloatWith6Digits: true,
	}.Froze()

	geojson.CustomJSONMarshaler = c
	geojson.CustomJSONUnmarshaler = c

	lines := geojson.NewFeatureCollection()
	points1 := geojson.NewFeatureCollection()
	points2 := geojson.NewFeatureCollection()

	for i := 0; i < 1000; i++ {
		line := orb.LineString{}
		for i := 0; i < 20; i++ {
			line = append(line, orb.Point{-30 + rand.Float64()*60, -30 + rand.Float64()*60})
		}
		lines.Append(geojson.NewFeature(line))
	}
	for i := 0; i < 10000; i++ {
		points1.Append(geojson.NewFeature(orb.Point{rand.Float64() * 50, rand.Float64() * 50}))
		points2.Append(geojson.NewFeature(orb.Point{rand.Float64() * 50, rand.Float64() * 50}))
	}

	// mc := Map{Container: "map"}
	// s, err := mc.Render()
	// if err != nil {
	// 	panic(1)
	// }
	// println(s.String())

	esc := mbgojs.NewScript(
		mbgojs.NewMap(mbgojs.Map{
			Container:   "map",
			AccessToken: "<token>",
			Hash:        true,
			Config: mbgojs.MapConfig{
				Basemap: mbgojs.BasemapConfig{
					Theme: "faded",
				},
			},
		}),
		mbgojs.NewConsoleLog("map"),
		mbgojs.NewMapOnLoad(
			mbgojs.NewMapAddImageRectangle("square", 10, 10),
			mbgojs.NewMapAddImageCircle("circle", 15),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id:   "layer",
				Type: "line",
				Source: mbgojs.MapSource{
					Type: "geojson",
					Data: *lines,
				},
			}),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id:   "points1",
				Type: "symbol",
				Source: mbgojs.MapSource{
					Type:       "geojson",
					Data:       *points1,
					GenerateId: true,
				},
				Layout: mbgojs.MapLayout{
					IconImage:        "square",
					IconAllowOverlap: true,
				},
			}),
			mbgojs.NewMapAddLayer(mbgojs.MapLayer{
				Id:   `points2`,
				Type: "symbol",
				Source: mbgojs.MapSource{
					Type:       "geojson",
					Data:       *points2,
					GenerateId: true,
				},
				Layout: mbgojs.MapLayout{
					IconImage:        "circle",
					IconAllowOverlap: true,
				},
			}),
			mapboxglgojs.NewMapOnEventLayer("click", "points1", mapboxglgojs.NewHtmxAjax(
				mapboxglgojs.HtmxAjax{
					Path: "/click",
					Verb: "GET",
					Context: mapboxglgojs.HtmxAjaxContext{
						Values: map[string]string{
							"event":   `"ok"`,
							"lat":     "e.lngLat.lat",
							"lng":     "e.lngLat.lng",
							"feature": "e.features[0].id",
						},
						Target: "#testtarget",
						Swap:   "innerHTML",
					},
				}),
			),
		),
		mbgojs.NewMapOnEventLayerCursor("mouseover", "points1", "pointer"),
		mbgojs.NewMapOnEventLayerCursor("mouseout", "points1", ""),
		mbgojs.NewMapOnEventLayerCursor("mouseover", "points2", "pointer"),
		mbgojs.NewMapOnEventLayerCursor("mouseout", "points2", ""),
		mbgojs.NewMapOnEventLayer("mouseout", "layer"),
	)

	s, err := esc.Render(mbgojs.RenderConfig{})
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
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
		t, _ := template.New("page").Parse(`
		<html>
		<head>
			<script src='https://api.mapbox.com/mapbox-gl-js/v3.15.0/mapbox-gl.js'></script>
			<link href='https://api.mapbox.com/mapbox-gl-js/v3.15.0/mapbox-gl.css' rel='stylesheet' />
			<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.7/dist/htmx.min.js"></script>

		</head>
		<body>
			<div id="testtarget"></div>
			<div id="map" style="width:100vw; height:100vh"></div>
			{{.}}
		</body>
		</html>
		`)

		if err := t.Execute(ctx.Writer, template.HTML(s)); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// curl --output data -H  "Accept-Encoding: gzip" localhost:8080
