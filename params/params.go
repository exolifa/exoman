package params

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Repertoires contains the list of different directories used
type Repertoires struct {
	Webtemplates string `json:"webtemplates"`
	Configfiles  string `json:"configfiles"`
	Logfiles     string `json:"logfiles"`
}

// Topics contains the list of the topics to subscribe to
type Topics struct {
	Topic string `json:"topic"`
}

// Mqttparams contains the list of the parameters to connect on MQTT brocker
type Mqttparams struct {
	User      string    `json:"mqttuser"`
	Password  string    `json:"mqttpw"`
	Brockerid string    `json:"mqttbrocker"`
	Subs      []*Topics `json:"topics"`
}

// Config structure contains the information read from the JSON parameter file
type Config struct {
	Dir     *Repertoires
	Mqtt    *Mqttparams
	Tcpport string `json:"tcpport"`
}

// Conf contains actual parameters
var Conf Config

func init() {
	conffyle := os.Args[1]
	fic, ficerr := ioutil.ReadFile(conffyle)
	if ficerr != nil {
		fmt.Printf("error ioutil : %v \n", ficerr)
	}
	_ = json.Unmarshal([]byte(fic), &Conf)
	fmt.Printf("Init summary\n=================\n")
	fmt.Printf("Received parameter:%v\n", conffyle)
	fmt.Printf("Param file title :%v\n", string(fic))
	fmt.Printf("lecture faite: %v\n", Conf)
}

// Getconfig allows all module to get value of any parameter
func Getconfig(cle string) string {

	return "error"
}
