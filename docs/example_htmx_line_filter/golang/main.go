package main

import (
	"html/template"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_line_filter/common"
)

// / ### [demo]
func main() {
	common.CreateLines()
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))

	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "public, max-age=3600")
		tmpl.Execute(ctx.Writer, template.JS(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.GET("/filter", func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "public, max-age=3600")
		ctx.Writer.Write([]byte(template.JS(common.Filter(ctx.Query("distance")))))
	})
	router.Run(":8080")
}

/// [demo]
