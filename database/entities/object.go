package entities

import (
	"fmt"
	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID        int
	StyleName string
	Geometry  geometry.Geometry
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, wkt string) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &MapObject{ID: id, Geometry: geo}, nil
}
