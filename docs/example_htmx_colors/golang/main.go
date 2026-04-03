package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_htmx_colors/common"
)

// / ### [demo]
func main() {
	router := gin.Default()

	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(c *gin.Context) {
		tmpl.Execute(c.Writer, template.JS(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.POST("/thickness", func(c *gin.Context) {
		thickStr := c.Query("thick")
		thick, _ := strconv.Atoi(thickStr)
		c.Writer.Write([]byte(template.JS(template.JS(common.PaintProperty("line-width", thick)))))
	})
	router.POST("/color", func(c *gin.Context) {
		c.Writer.Write([]byte(template.JS(template.JS(common.PaintProperty("line-color", fmt.Sprintf("\"%s\"", c.Query("color")))))))
	})
	router.Run(":8080")
}

/// [demo]
