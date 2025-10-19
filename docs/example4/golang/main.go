package main

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example4/common"
)

// / [demo]
func main() {
	router := gin.Default()

	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(c *gin.Context) {
		tmpl.Execute(c.Writer, template.JS(common.Example()))
	})
	router.Run(":8080")
}

/// [demo]
