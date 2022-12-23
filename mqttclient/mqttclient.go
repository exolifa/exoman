package mqttclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"math/rand"
	"strconv"
	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/params"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/secsy/goftp"
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

var scanreceived bool

// Client is the  session used by the different processes
var Client mqtt.Client
var clientid string


func addtoiotlist(candi string, payload []byte) {
	tmp := new(Iotdevice)
	//	fmt.Printf("device: %s and payload %v\n", candi, string(payload[:]))
	_ = json.Unmarshal([]byte(payload), &tmp)
	tmp.Complist = map[string]*IoTComp{}
	//	fmt.Printf("tmp=%v\n", tmp)
	iotlist[candi] = tmp
	iotlist[candi].Complist = make(map[string]*IoTComp)
}

//Commande function collects all the probes from a device
func Commande(target string, oper string) {
	topic := "cmd/" + target
	switch oper {
	case "scan":
		msg := "SCAN"
		scanreceived = false
		Publication(Client, topic, msg)
		for scanreceived {
			time.Sleep(1 * time.Second)
		}
	case "reboot":
		msg := "REBOOT"
		Publication(Client, topic, msg)
	case "register":
		msg := "REGISTER"
		Publication(Client, topic, msg)
	case "info":
		msg := "INFO"
		Publication(Client, topic, msg)
	case "inventaire":
		topic := "broadcast/"
		msg := "INFO"
		Publication(Client, topic, msg)
	}
}

//Getlist send the list of currently discovered IoT
func Getlist() map[string]*Iotdevice {
	return iotlist
}

// Getdevices send a list of all discovered devices
func Getdevices() []string {
	var tmp []string
	for cle := range iotlist {
		tmp = append(tmp, cle)
	}
	return tmp
}

// GetIP return the ip address of a specific device
func GetIP(target string) string {
	go logger.Logme("global", "mqttclient", "GetIP", "info", fmt.Sprintf("liste des devices: %s\n",Getdevices()))
	go logger.Logme("global", "mqttclient", "GetIP", "info", fmt.Sprintf("IP adress for %s is %s\n", target,iotlist[target].Devip))
	return iotlist[target].Devip
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
	client = Connect(clientid)
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
				if rcvtopic[2] == "scan" { //this a response to a SCAN request
					Buildfromscan(rcviot, iotlist[rcviot].Devip, rcvpayload)
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
	rand.Seed(time.Now().UnixNano())
	clientid= "exoman-client"+strconv.Itoa(rand.Intn(2))
	Client = Connect(clientid)

	fmt.Println("subscribing to all messages")
	Client.Subscribe("#", 0, Processmsg)
	iotlist = make(map[string]*Iotdevice)
}

// Ftp is the function to communicate with the file system of the devices
func Ftp(target string, ip string, direction int) {
	// Create client object with default config
	cible := params.Getconfig("Configfiles") + "actual\\" + target + ".json"
	config := goftp.Config{
		User:               iotlist[target].Ftpuser,
		Password:           iotlist[target].Ftppw,
		ConnectionsPerHost: 1,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}
	client, err := goftp.DialConfig(config, ip)
	if err != nil {
		go logger.Logme("global", "config", "Ftp", "fatal", fmt.Sprintf("error creating ftp client:%v for ip:%v and config: %v", err,ip,config))
		panic(err)
	}
	switch direction {
	case 0: //this is a download
		fic, err := os.Create(cible)
		if err != nil {
			go logger.Logme("global", "config", "Ftp", "fatal", fmt.Sprintf("error downloading disk config:%v", err))
		}
		err = client.Retrieve("config.json", fic)
		if err != nil {
			go logger.Logme("global", "config", "Ftp", "fatal", fmt.Sprintf("error retrieving device config:%v", err))
		}
	case 1: //this is an upload
		if fileExists(cible) {
			fic, err := os.Open(cible)
			if err != nil {
				go logger.Logme("global", "config", "FTP", "fatal", fmt.Sprintf("error open disk config:%v", err))
			}
			err = client.Store("config.json", fic)
			if err != nil {
				go logger.Logme("global", "config", "FTP", "fatal", fmt.Sprintf("error writing device config:%v", err))
			}
		} else {
			go logger.Logme("global", "config", "FTP", "error", fmt.Sprintf("file does not exists on disk(%s):%v", cible, err))
		}
	}
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Buildfromscan write an actual file from scan result
func Buildfromscan(target string, ip string, content []byte) {
	// first get the json object from actual file
	// import the actual configuration
	Ftp(target, ip, 0)
	cible := params.Getconfig("Configfiles") + "actual/" + target + ".json"
	readcfg, err := ioutil.ReadFile(cible)
	if err != nil {
		go logger.Logme("global", "config", "Buildfromscan", "error", fmt.Sprintf("error reading actual  config:%v", err))
	}
	currcfg := map[string]interface{}{}
	scancfg := map[string]interface{}{}
	json.Unmarshal([]byte(readcfg), &currcfg)
	json.Unmarshal(content, &scancfg)
	currcfg["bus"] = scancfg["bus"]
	cfgfile, _ := os.OpenFile(cible, os.O_CREATE, os.ModePerm)
	defer cfgfile.Close()
	encoder := json.NewEncoder(cfgfile)
	encoder.Encode(currcfg)
}
