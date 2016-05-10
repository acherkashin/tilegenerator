package mapobjects

import "github.com/paulsmith/gogeos/geos"

// MapObject represents a geometry on a map
type MapObject struct {
	ID       int
	Geometry geos.Geometry
	CSS      string
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, wkt, css string) (*MapObject, error) {
	geometry, err := geos.FromWKT(wkt)
	if err != nil {
		return nil, err
	}
	return &MapObject{ID: id, Geometry: *geometry, CSS: css}, nil
}
