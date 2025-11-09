package main

import (
	"html/template"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example8/common"
)

// / ### [demo]
func main() {
	common.GeneratePoints()
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))

	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(ctx *gin.Context) {
		tmpl.Execute(ctx.Writer, template.JS(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.GET("/popup", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte(template.HTML(common.Popup(ctx.Query("featureId")))))
	})
	router.Run(":8080")
}

/// [demo]
