package main

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_kmeans/common"
)

// / ### [demo]
func main() {
	router := gin.Default()
	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(c *gin.Context) {
		tmpl.Execute(c.Writer, template.JS(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.GET("/kmean-lines", func(c *gin.Context) {
		c.Writer.Write([]byte(template.JS(common.KmeanClusterLines(c.Query("k")))))
	})
	router.GET("/kmean-points", func(c *gin.Context) {
		c.Writer.Write([]byte(template.JS(common.KmeanClusterPoints(c.Query("k")))))
	})
	router.Run(":8080")
}

/// [demo]
