package mapobjects

import (
	"github.com/TerraFactory/tilegenerator/geo"
	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID       int
	TypeID   int
	Geometry geometry.Geometry
	CSS      string
	Attrs    []geo.BaseAttribute
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeID int, wkt, css string, Attrs []geo.BaseAttribute) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		return nil, err
	}
	return &MapObject{ID: id, TypeID: typeID, Geometry: geo, CSS: css, Attrs: nil}, nil
}
