package main

import (
	"github.com/TerraFactory/tilegenerator/listeners"
	"github.com/TerraFactory/tilegenerator/settings"
)

func main(){
	conf := settings.GetSettings()
	listeners.Listen(conf)
}