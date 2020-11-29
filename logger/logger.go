package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"exolifa.com/exoman/params"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

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
