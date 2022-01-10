package main

import (
	"fmt"
	"icombo/pkg/icombo"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

const defaultConfigFilePath string = "icombo.toml"

var processImagesInput icombo.ProcessImagesInput

func main() {
	log.Println("starting icombo")
	initConfig()
	startTime := time.Now()
	if err := icombo.ProcessImages(processImagesInput); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("icombo complete in", time.Since(startTime))
}

func initConfig() {

	// first check to see if we have config file in current dir
	if _, err := os.Stat(fmt.Sprint("./", defaultConfigFilePath)); os.IsNotExist(err) {
		log.Print("no config file found at ./", defaultConfigFilePath, ", checking for example folder used during development")

		// if not try to move to example folder used during dev
		if _, err := os.Stat(fmt.Sprint("./example/", defaultConfigFilePath)); os.IsNotExist(err) {
			log.Fatalln("no config file found")
		}
		log.Print("config file found a ./example/", defaultConfigFilePath, " executing from ./example")
		os.Chdir("./example")
	}

	currentDirectory, _ := os.Getwd()
	viper.AddConfigPath(currentDirectory)
	viper.SetConfigType("toml")
	viper.SetConfigName("icombo.toml")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&processImagesInput); err != nil {
		log.Fatalln(err)
	}
}
