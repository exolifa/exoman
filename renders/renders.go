package renders

import (
	"fmt"
	"net/http"

	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/mqttclient"
	"github.com/gin-gonic/gin"
)

// render to handle all types of request (html ,json,xml
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

// Statuspage return full list of IoTs
func Statuspage(c *gin.Context) {
	go logger.Logme("global", "renders", "statuspage", "info", fmt.Sprint("building status page"))
	mycomponentslist := mqttclient.Getlist()
	render(c, gin.H{
		"title":   "Status of all identified IoTs",
		"payload": mycomponentslist}, "status.html")
}
