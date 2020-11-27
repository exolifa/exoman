package routing

import (
	"exolifa.com/exoman/params"
	"github.com/gin-gonic/gin"
)

// SetRoutes defines de different routes for gin services
func SetRoutes() *gin.Engine {

	r := gin.Default()
	templatedir := params.Getconfig("webtemplates")
	r.LoadHTMLGlob(templatedir)
	// this is the first page ...equivalent to the index.html reference
	//r.GET("/", processors.Carsliste)
	// this is he route to hanle all requests from the form (carform.html)
	//r.POST("/formcars", processors.FormCars)
	return r
}
