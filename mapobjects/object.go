package mapobjects

import (
	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID       int
	TypeID   int
	Geometry geometry.Geometry
	CSS      string
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeId int, wkt, css string) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		return nil, err
	}
	return &MapObject{ID: id, TypeID: typeId, Geometry: geo, CSS: css}, nil
}
