package main

import (
	"flag"

	"github.com/TerraFactory/tilegenerator/listeners"
	"github.com/TerraFactory/tilegenerator/settings"
	"os"
	"fmt"
)

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

	if(err != nil) {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	listeners.Listen(conf)
}
