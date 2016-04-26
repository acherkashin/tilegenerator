package mapobjects

import "github.com/paulsmith/gogeos/geos"

type MapObject struct {
	Geometry geos.Geometry
	CSS      string
}

func NewObject(wkt, css string) (*MapObject, error) {
	geometry, err := geos.FromWKT(wkt)
	if (err != nil) {
		return nil, err
	}
	return &MapObject{Geometry: *geometry, CSS: css }, nil
}
