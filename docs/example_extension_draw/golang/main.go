package main

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_extension_draw/common"
)

// / ### [demo]
func main() {
	router := gin.Default()
	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(c *gin.Context) {
		tmpl.Execute(c.Writer, template.JS(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.Run(":8080")
}

/// [demo]
