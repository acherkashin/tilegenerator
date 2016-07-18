package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/TerraFactory/tilegenerator/listeners"
	"github.com/TerraFactory/tilegenerator/settings"
	"github.com/TerraFactory/tilegenerator/utils"

	"os"
)

func setOutputFileForLog(folder string) error {
	pathLogFile := folder + "/logfile.log"

	if !utils.FileExists(&pathLogFile) {
		_, err := os.Create(pathLogFile)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(pathLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file for: %v", err)
	}

	log.SetOutput(f)

	return nil
}

func main() {
	var help = flag.Bool("h", false, "Display this message.")
	var conf_path = flag.String("c", "./config.toml", "Absolute path to configuration file")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
		return
	}
	conf, err := settings.GetSettings(conf_path)

	if err = setOutputFileForLog(conf.LogDirectory); err != nil {
		fmt.Println(err.Error())
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	listeners.StartApplication(conf)
}
