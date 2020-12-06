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
	// check if the specified directories exist if not create them
	if _, err := os.Stat(Conf.Configfiles); os.IsNotExist(err) {
		// config files directory does not exist
		err := os.Mkdir(Conf.Configfiles, 0755)
		if err != nil {
			fmt.Printf("unable to create %s ;error =%v\n", Conf.Configfiles, err)
		}
	}
	actualfiles := Conf.Configfiles + "actual\\"
	if _, err := os.Stat(actualfiles); os.IsNotExist(err) {
		// actuel config files directory does not exist
		err := os.Mkdir(actualfiles, 0755)
		if err != nil {
			fmt.Printf("unable to create %s ;error =%v\n", actualfiles, err)
		}
	}
	if _, err := os.Stat(Conf.Logfiles); os.IsNotExist(err) {
		// log files directory does not exist
		err := os.Mkdir(Conf.Logfiles, 0755)
		if err != nil {
			fmt.Printf("unable to create %s ;error =%v\n", Conf.Logfiles, err)
		}
	}
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
