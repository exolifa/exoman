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
	Oldtarget  string
}

// Logrec Screen of log
type Logrec struct {
	Oldlog   string
	Logfyles []string
	Logenreg []logger.Logdata
}

// version of the software
const version = "1.0.0"
const author = "Daniel Vielvoye"
const contact = "exolifa@gmail.com"
const pubdate = "6 décembre 2020"
const licence = "All Rigth Reserved"

// Getversion provides info for About pages
func Getversion() []string {
	var tmp []string
	tmp = append(tmp, version)
	tmp = append(tmp, author)
	tmp = append(tmp, contact)
	tmp = append(tmp, pubdate)
	tmp = append(tmp, licence)
	return tmp
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
		mqttclient.Commande(target, "scan")
	case "Reboot":
		mqttclient.Commande(target, "reboot")
	case "Register":
		mqttclient.Commande(target, "register")
	case "Refresh":
		mqttclient.Commande(target, "info")
	case "Inventaire":
		mqttclient.Commande(target, "inventaire")
	case "Update":
		cfg := c.PostForm("configfile")
		config.Save(target, []byte(cfg))
	case "Upload":
		mqttclient.Ftp(target, mqttclient.GetIP(target), 1)
	case "Retrieve":
		mqttclient.Ftp(target, mqttclient.GetIP(target), 0)
	default:
	}
	myconfiglist.Oldtarget = target
	myconfiglist.Devicelist = mqttclient.Getdevices()
	myconfiglist.Configlist = config.Get(target)
	fmt.Printf("données reçues:%v\n", myconfiglist)
	//c.HTML(http.StatusOK, "modules.html", data)
	render(c, gin.H{
		"title":   "Status of all identified IoTs",
		"payload": myconfiglist}, "modules.html")
}

// Logpage return full list options to configure, view logs ,etc of a module
func Logpage(c *gin.Context) {
	go logger.Logme("global", "renders", "logpage", "info", fmt.Sprint("building log viewer page"))
	var myloglist Logrec
	target := c.PostForm("target")
	if target == "" {
		target = "global.log"
	}
	myloglist.Oldlog = target
	myloglist.Logfyles = logger.Loglist()
	myloglist.Logenreg = logger.Logview(target)
	render(c, gin.H{
		"title":   "Status of all identified IoTs",
		"payload": myloglist}, "logviewer.html")
}

// Aboutpage return full list options to configure, view logs ,etc of a module
func Aboutpage(c *gin.Context) {
	go logger.Logme("global", "renders", "aboutpage", "info", fmt.Sprint("building about page"))
	render(c, gin.H{
		"title":   "Status of all identified IoTs",
		"payload": Getversion()}, "about.html")
}
