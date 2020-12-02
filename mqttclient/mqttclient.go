package mqttclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/params"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// IoTComp contains information of the device's sensors and actuators it is maintained with all info circulating on the MQTT
type IoTComp struct {
	Lastvalue string
	Lasttime  string
}

// Iotdevice houses the list of discovered components by listening to all MQTT topics
type Iotdevice struct {
	Devip      string `json:"IP"`
	Devversion string `json:"Version"`
	Devname    string `json:"deviceName"`
	Ftpuser    string `json:"FTP user"`
	Ftppw      string `json:"FTP pw"`
	Complist   map[string]*IoTComp
}

func (iot Iotdevice) String() string {
	reponse := iot.Devname + "(" + iot.Devversion + ")" + " address=" + iot.Devip + "\n" + "FTP: user=" + iot.Ftpuser + " - pw=" + iot.Ftppw
	for cle, val := range iot.Complist {
		//			reponse = reponse + "\n composant: "+cle +" time=" + iot.Complist[cle].Lasttime +" data="+ iot.Complist[cle].Lastvalue
		reponse = reponse + "\n composant: " + cle + " time=" + val.Lasttime + " data=" + val.Lastvalue

	}
	return reponse
}

// iotlist
var iotlist map[string]*Iotdevice

func addtoiotlist(candi string, payload []byte) {
	tmp := new(Iotdevice)
	//	fmt.Printf("device: %s and payload %v\n", candi, string(payload[:]))
	_ = json.Unmarshal([]byte(payload), &tmp)
	tmp.Complist = map[string]*IoTComp{}
	//	fmt.Printf("tmp=%v\n", tmp)
	iotlist[candi] = tmp
	iotlist[candi].Complist = make(map[string]*IoTComp)
}

//Getlist send the list of currently discovered IoT
func Getlist() map[string]*Iotdevice {
	return iotlist
}

func Getdevices() []string {
	var tmp []string
	for cle := range iotlist {
		tmp = append(tmp, cle)
	}
	return tmp
}

// Connect initiates a connections to MQTT brocker
func Connect(clientID string) mqtt.Client {
	opts := createClientOptions(clientID)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		go logger.Logme("global", "mqttclient", "Connecting", "fatal", fmt.Sprint(token.Error()))
		panic(token.Error())
	} else {
		go logger.Logme("global", "mqttclient", "Connecting", "info", fmt.Sprintf("Connected to tcp://%s\n", params.Getconfig("Brockerid")))
	}
	return client
}

func reconnect(client mqtt.Client, err error) {
	go logger.Logme("global", "mqttclient", "Connection lost", "error", fmt.Sprintf("lost connection with error :%s\n", err))
	client = Connect("exoman-sub")
	client.Subscribe("#", 0, Processmsg)
}
func connLostHandler(c mqtt.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)

	//Perform additional action...
}

func createClientOptions(clientID string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", params.Getconfig("Brockerid")))
	opts.SetConnectionLostHandler(reconnect)
	opts.AutoReconnect = true
	opts.SetUsername(params.Getconfig("User"))
	password := params.Getconfig("Password")
	opts.SetPassword(password)
	opts.SetClientID(clientID)
	//	fmt.Printf("connection options: %v\n", opts)
	return opts
}
func gettime() string {
	loc, err := time.LoadLocation(params.Getconfig("Location"))
	if err != nil {
		go logger.Logme("global", "mqttclient", "gettime", "error", fmt.Sprint(err))
	}
	now := time.Now().In(loc)
	return fmt.Sprint(now)
}

// Processmsg wait to receive MQTT message
func Processmsg(client mqtt.Client, message mqtt.Message) {
	rcvtopic := strings.Split(message.Topic(), "/")
	rcvpayload := message.Payload()
	rcvtype := rcvtopic[0]
	rcviot := rcvtopic[1]
	go logger.Logme(rcviot, "mqttclient", rcvtype, "info", string(rcvpayload[:]))
	//	fmt.Printf("received message type: %s\n", rcvtype)
	if rcvtype == "info" || rcvtype == "debug" {
		if _, ok := iotlist[rcviot]; !ok {
			//			fmt.Printf("%v is new\n", rcviot)
			if rcvtype == "debug" {
				if rcvtopic[2] == "info" { //this is a response to a INFO request
					addtoiotlist(rcviot, rcvpayload)
				}
			} else {
				//publish cmd info
				cmdtopic := "cmd/" + rcviot
				go Publication(client, cmdtopic, "INFO")
			}
		} else {
			//device already registered : collecting sensors/actuator state
			if rcvtype == "info" {
				rcvcomp := rcvtopic[2]
				temp := new(IoTComp)
				temp.Lastvalue = fmt.Sprint(string(rcvpayload[:]))
				temp.Lasttime = fmt.Sprint(gettime())
				iotlist[rcviot].Complist[rcvcomp] = temp
			}
		}
	}
}

// Publication send a command on MQTT
func Publication(client mqtt.Client, topic string, msg string) {
	client.Publish(topic, 0, false, msg)
}

// Listen captures all messages in MQTT
func init() {
	fmt.Println("initiating connect ")
	client := Connect("exoman-sub")
	fmt.Println("subscribing to all messages")
	client.Subscribe("#", 0, Processmsg)
	iotlist = make(map[string]*Iotdevice)
}
