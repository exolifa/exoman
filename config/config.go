package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"exolifa.com/exoman/logger"
	"exolifa.com/exoman/params"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Get reads the config file from disk and
func Get(target string) []string {
	var temp []string
	cible := params.Getconfig("Configfiles") + target + ".json"
	if fileExists(cible) {
		content, err := ioutil.ReadFile(cible)
		if err != nil {
			go logger.Logme("global", "config", "ConfigGet", "fatal", fmt.Sprintf("error reading disk config:%v", err))
		}

		// Convert []byte to string and print to screen
		temp[0] = string(content)
	} else {
		temp[0] = ""
	}
	cible1 := params.Getconfig("Configfiles") + "actual/" + target + ".json"
	if fileExists(cible1) {
		content1, err := ioutil.ReadFile(cible1)
		if err != nil {
			go logger.Logme("global", "config", "ConfigGet", "fatal", fmt.Sprintf("error reading actual config:%v", err))
		}

		// Convert []byte to string and print to screen
		temp[1] = string(content1)
	} else {
		temp[1] = ""
	}
	return temp
}
