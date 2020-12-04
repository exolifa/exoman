package renders

import (
	"fmt"
	"net/http"

	"exolifa.com/exoman/config"
	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/mqttclient"
	"github.com/gin-gonic/gin"
)

// Configdata donnée pour envoyer vers module.html
type Configdata struct {
	Devicelist []string
	Configlist []string
}

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

// Modulepage return full list options to configure, view logs ,etc of a module
func Modulepage(c *gin.Context) {
	go logger.Logme("global", "renders", "modulespage", "info", fmt.Sprint("building modules page"))
	var myconfiglist Configdata
	action := c.PostForm("action")
	target := c.PostForm("target")
	switch action {
	case "Scan":
	case "Update":
	case "Save":
	case "Retrieve":
	default:
	}
	myconfiglist.Devicelist = mqttclient.Getdevices()
	myconfiglist.Configlist = config.Get(target)
	fmt.Printf("données reçues:%v\n", myconfiglist)
	data := make(map[string][]string)
	data["demo"] = myconfiglist.Devicelist
	data["autre"] = myconfiglist.Configlist
	//c.HTML(http.StatusOK, "modules.html", data)
	render(c, gin.H{
		"title":   "Status of all identified IoTs",
		"payload": myconfiglist}, "modules.html")
}
