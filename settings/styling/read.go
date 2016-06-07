package styling

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/TerraFactory/tilegenerator/settings"
	"github.com/TerraFactory/tilegenerator/utils"
	"github.com/pelletier/go-toml"
	"strings"
	"github.com/TerraFactory/wktparser/geometry"
)

var styles *map[string]Style
var errs []error
var once sync.Once

func parseType(t string) (int, error) {
	switch strings.ToUpper(t) {
	case "POINT":
		return geometry.TPoint, nil
	case "MULTIPOINT":
		return geometry.TMultiPoint, nil
	case "LINE":
	case "POLYLINE":
	case "LINESTRING":
		return geometry.TLineString, nil
	case "MULTILINE":
	case "MULTIPOLYLINE":
	case "MULTILINESTRING":
		return geometry.TMultiLineString, nil
	case "POLYGON":
		return geometry.TPolygon, nil
	}
	return -1, errors.New(fmt.Sprintf("Failed to parse geometry type %s", t))
}

func readStylesFile(filename string) (*Style, error) {
	styles, err := toml.LoadFile(filename)
	if err != nil {
		return nil, err
	}
	style := Style{}
	style.Name = styles.Get("Name").(string)
	geometryType, err := parseType(styles.Get("GeometryType").(string))
	if err != nil {
		return nil, err
	}
	style.GeometryType = geometryType
	prims := styles.Get("primitives").([]*toml.TomlTree)
	for _, p := range prims {
		t := p.Get("Type").(string)
		primitive, _ := NewPrimitive(t, p.ToMap())
		style.Primitives = append(style.Primitives, primitive)
	}
	return &style, nil
}

func readStylesDirectory(directory string) (*map[string]Style, []error) {
	result := map[string]Style{}
	allErrors := []error{}

	if utils.IsDirectory(directory) {
		return nil, []error{errors.New(fmt.Sprint("Path %s is not a directory", directory))}
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
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
	return styles, errs
}
