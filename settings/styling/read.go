package styling

import (
	"sync"
	"github.com/TerraFactory/tilegenerator/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/pelletier/go-toml"
	"github.com/TerraFactory/tilegenerator/settings"
)

var styles *map[string]Style
var errs []error
var once sync.Once

func readStylesFile(filename string) (*Style, error) {
	toml, err := toml.LoadFile(filename)
	if err != nil {
		return nil, err
	}
	style := Style{}
	style.Name = toml.Get("name").(string)
	return &style, nil
}

func readStylesDirectory(directory string) (*map[string]Style, []error) {
	result := map[string]Style{}
	allErrors := []error{}

	if utils.IsDirectory(directory) {
		return nil, []error{errors.New(fmt.Sprint("Path %s is not a directory", directory))}
	}
	files, err := ioutil.ReadDir(directory)
	if (err != nil) {
		return nil, []error{err}
	}
	for _, file := range files {
		path := directory + "/" + file.Name()
		if file.IsDir() {
			styles, errs := readStylesDirectory(path)
			for _, style := range *styles {
				result[style.Name] = style
			}
			allErrors = append(allErrors, errs...)
		} else {
			style, err := readStylesFile(path)
			if err == nil {
				result[style.Name] = *style
			} else {
				allErrors = append(allErrors, err)
			}
		}
	}
	return &result, allErrors
}

func GetStyles(conf *settings.Settings) (*map[string]Style, []error) {
	once.Do(func() {
		styles, errs = readStylesDirectory(conf.StylesDirectory)
	})
	fmt.Println(*styles)
	return styles, errs
}