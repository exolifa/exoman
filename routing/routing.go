package routing

import (
	"fmt"

	"exolifa.com/exoman/params"
	"exolifa.com/exoman/renders"
	"github.com/gin-gonic/gin"
)

// SetRoutes defines de different routes for gin services
func SetRoutes() *gin.Engine {

	r := gin.Default()
	templatedir := params.Getconfig("Webtemplates")
	fmt.Printf("reading templates from:%v\n", templatedir)
	r.LoadHTMLGlob(templatedir)
	// this is the first page ...equivalent to the index.html reference
	r.GET("/", renders.Statuspage)
	// this is he route to hanle all requests from the form (carform.html)
	//r.POST("/formcars", processors.FormCars)
	return r
}
