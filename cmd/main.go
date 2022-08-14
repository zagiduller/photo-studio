package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	psMinio "photostudio/components/minio"
)

// @project photo-studio
// @created 27.07.2022

func init() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := psMinio.CreateNewUserClient("Sibay")
	if err != nil {
		log.Fatal(err)
		return
	}
	files, err := client.GetUserFiles()
	if err != nil {
		log.Error(err)
	}
	for i, file := range files {
		log.Infof("%d) Location: %s ", i+1, file.Key)
	}
}
