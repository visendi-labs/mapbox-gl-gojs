package main

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_raster/common"
)

// / ### [demo]
func main() {
	router := gin.Default()
	common.ReadFiles()
	router.Static("/example_htmx_raster/common/weather/", "../common/weather") // Has to be this path for docs live variant to work..
	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(ctx *gin.Context) {
		tmpl.Execute(ctx.Writer, template.JS(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.GET("/timeline", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte(template.HTML(common.UpdateUrl(ctx.Query("t")))))
	})
	router.Run(":8080")
}

/// ### [demo]
