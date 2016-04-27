package mapobjects

import "github.com/paulsmith/gogeos/geos"

type MapObject struct {
	Id       int
	Geometry geos.Geometry
	CSS      string
}

func NewObject(id int, wkt, css string) (*MapObject, error) {
	geometry, err := geos.FromWKT(wkt)
	if (err != nil) {
		return nil, err
	}
	return &MapObject{Id: id, Geometry: *geometry, CSS: css }, nil
}
