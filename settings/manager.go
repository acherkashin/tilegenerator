package settings

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"sync"
)

// Settings is a singleton object, which contains configuration of a tile server
type Settings struct {
	DBConnectionString string
	DBGeometryTable    string
	DBGeometryColumn   string
	DBInstanceName     string
	HTTPPort           string
}

var instance *Settings
var once sync.Once

func readSettings() *Settings {
	var settings Settings
	config, err := toml.LoadFile("./config.toml")
	if err != nil {
		fmt.Println("Error ", err.Error())
	} else {
		settings = Settings{
			DBConnectionString: config.Get("database.connection_string").(string),
			DBGeometryTable:    config.Get("database.geometry_table").(string),
			DBGeometryColumn:   config.Get("database.geometry_column").(string),
			DBInstanceName:     config.Get("database.instance_name").(string),
			HTTPPort:           config.Get("http.port").(string),
		}
	}
	return &settings
}

// GetSettings returns single instance of the Settings structure.
// By default it reads configuration from "config.toml" file
func GetSettings() *Settings {
	once.Do(func() {
		instance = readSettings()
	})
	return instance
}
