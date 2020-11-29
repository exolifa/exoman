package params

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config structure contains the information read from the JSON parameter file
type Config struct {
	Webtemplates string `json:"webtemplates"`
	Configfiles  string `json:"configfiles"`
	Logfiles     string `json:"logfiles"`
	Tcpport      string `json:"tcpport"`
	User         string `json:"mqttuser"`
	Password     string `json:"mqttpw"`
	Brockerid    string `json:"mqttbrocker"`
	Location     string `json:"timezone"`
}

// Conf contains actual parameters
var Conf Config

func init() {
	conffyle := os.Args[1]
	fic, ficerr := ioutil.ReadFile(conffyle)
	if ficerr != nil {
		fmt.Printf("error ioutil : %v \n", ficerr)
	}
	//	fmt.Printf(" content is :%v\n", fic)
	_ = json.Unmarshal([]byte(fic), &Conf)
	fmt.Printf("Init summary\n=================\n")
	fmt.Printf("Received parameter:%v\n", conffyle)
	fmt.Printf("Param file title :%v\n", string(fic))
	fmt.Printf("lecture faite: %v\n", Conf)
	fmt.Printf("%v, %v, %v, %v,%v", Conf.Configfiles, Conf.Logfiles, Conf.Webtemplates, Conf.Brockerid, Conf.Tcpport)
}

// Getconfig allows all module to get value of any parameter
func Getconfig(cle string) string {
	switch cle {
	case "Webtemplates":
		return Conf.Webtemplates
	case "Configfiles":
		return Conf.Configfiles
	case "Logfiles":
		return Conf.Logfiles
	case "Tcpport":
		return Conf.Tcpport
	case "User":
		return Conf.User
	case "Password":
		return Conf.Password
	case "Brockerid":
		return Conf.Brockerid
	case "Location":
		return Conf.Location
	}
	return "error"
}
