package settings

import (
	"errors"
	"fmt"
	"sync"

	"github.com/TerraFactory/tilegenerator/utils"
	"github.com/pelletier/go-toml"
)

// Settings is a singleton object, which contains configuration of a tile server
type Settings struct {
	DBConnectionString string
	DBGeometryTable    string
	DBGeometryColumn   string
	DBInstanceName     string
	HTTPPort           string
	StylesDirectory    string
	UrlAPI             string
	LogDirectory       string
}

var instance *Settings
var once sync.Once

func readSettings(conf_path *string) (*Settings, error) {
	var settings Settings
	if !utils.FileExists(conf_path) {
		return nil, errors.New(fmt.Sprintf("Error. File %s does not exist.", *conf_path))
	}
	config, err := toml.LoadFile(*conf_path)
	if err != nil {
		return nil, err
	} else {
		settings = Settings{
			DBConnectionString: config.Get("database.connection_string").(string),
			DBGeometryTable:    config.Get("database.geometry_table").(string),
			DBGeometryColumn:   config.Get("database.geometry_column").(string),
			DBInstanceName:     config.Get("database.instance_name").(string),
			HTTPPort:           config.Get("http.port").(string),
			StylesDirectory:    config.Get("styles.directory").(string),
			UrlAPI:             config.Get("api.url").(string),
			LogDirectory:       config.Get("logging.directory").(string),
		}
	}
	return &settings, nil
}

// GetSettings returns single instance of the Settings structure.
// By default it reads configuration from "config.toml" file
func GetSettings(conf_path *string) (*Settings, error) {
	var err error
	once.Do(func() {
		instance, err = readSettings(conf_path)
	})
	return instance, err
}
