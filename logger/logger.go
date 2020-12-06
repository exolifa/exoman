package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"

	"exolifa.com/exoman/params"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

// Logdata is the structure of the log files
type Logdata struct {
	Time   string `json:"time"`
	Msg    string `json:"msg"`
	Module string `json:"module"`
	Origin string `json:"origin"`
	Level  string `json:"level"`
}

var mutex sync.Mutex

// init will verify the log directory exists and initiate log rotate on it
func init() {
	logdir := params.Getconfig("Logfiles")
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)
	path := logdir + "/old.UTC."
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", path, "%Y-%m-%d.%H:%M:%S"),
		rotatelogs.WithLinkName(logdir),
		rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationTime(time.Hour*6),
	)
	if err != nil {
		Logme("global", "logger", "not MQTT", "fatal", fmt.Sprintf("Failed to Initialize Log File %s", err))
	}
	Logme("global", "logger", "not MQTT", "info", "System Logger started")
	log.SetOutput(writer)
	return
}

//Getlogs provides a list of all existing logs
func Getlogs() []string {
	var tmp []string
	logdir := params.Getconfig("Logfiles")
	f, err := os.Open(logdir)
	if err != nil {
		Logme("global", "logger", "GetLog", "fatal", fmt.Sprintf("Failed to Open log directory %s", err))
		return nil
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		Logme("global", "logger", "GetLog", "fatal", fmt.Sprintf("Failed list directory content %s", err))
		return nil
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	for i := 0; i < len(list); i++ {
		tmp = append(tmp, list[i].Name())
	}
	return tmp

}

// Logme will allow logging to all programs based on logrus (it just provide a single point to change if decided to use other logger)
func Logme(target string, module string, topic string, level string, payload string) {
	mutex.Lock()
	logdir := params.Getconfig("Logfiles")
	glblog := logdir + target + ".log"
	file, err := os.OpenFile(glblog, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("raté l'ouverture:%v\n", err)
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	switch level {
	case "info":
		log.WithFields(log.Fields{"module": module, "origin": topic}).Info(payload)
	case "debug":
		log.WithFields(log.Fields{"module": module, "origin": topic}).Debug(payload)
	case "fatal":
		log.WithFields(log.Fields{"module": module, "origin": topic}).Fatal(payload)
	case "error":
		log.WithFields(log.Fields{"module": module, "origin": topic}).Error(payload)
	default:
		log.WithFields(log.Fields{"module": module, "origin": topic}).Warn(payload)
	}
	mutex.Unlock()
}

// Loglist returns the list of available log files
func Loglist() []string {
	var tmp []string
	f, err := os.Open(params.Getconfig("Logfiles"))
	if err != nil {
		Logme("global", "logger", "Loglist", "fatal", fmt.Sprintf("Failed to open %s error is %v", params.Getconfig("Logfiles"), err))
		return nil
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		Logme("global", "logger", "Loglist", "fatal", fmt.Sprintf("Failed to read %s error is %v", params.Getconfig("Logfiles"), err))
		return nil
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	for i := 0; i < len(list); i++ {
		tmp = append(tmp, list[i].Name())
	}
	return tmp
}

// Logview returns a string with the log records for the targetted log file
func Logview(target string) []Logdata {
	var tmp []Logdata
	logfyle := params.Getconfig("Logfiles") + target
	fic, ficerr := ioutil.ReadFile(logfyle)
	if ficerr != nil {
		Logme("global", "logger", "Loglist", "error", fmt.Sprintf("Failed to read %s error is %v", logfyle, ficerr))
	} else {
		for _, line := range bytes.Split(fic, []byte{'\n'}) {
			var v Logdata
			if err := json.Unmarshal(line, &v); err != nil {
				Logme("global", "logger", "Loglist", "erro", fmt.Sprintf("Failed to unmarshall %v error is %v", line, err))
			}
			tmp = append(tmp, v)
		}
	}
	return tmp
}
