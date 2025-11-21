package main

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/visendi-labs/mapbox-gl-gojs/docs/example_ssr_extension_draw/common"
)

// / ### [demo]
func main() {
	router := gin.Default()
	tmpl := template.Must(template.ParseFiles("index.html"))
	router.GET("/", func(c *gin.Context) {
		tmpl.Execute(c.Writer, template.HTML(common.Example(os.Getenv("MAPBOX_ACCESS_TOKEN"))))
	})
	router.POST("/create", func(c *gin.Context) {
		c.Writer.Write([]byte(template.HTML(common.AddFeatures(c.PostForm("features")))))

	})
	router.Run(":8080")
}

/// [demo]
