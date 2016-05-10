package geo

import (
	"github.com/paulsmith/gogeos/geos"
)

type BaseGeometry struct {
	Id    int
	Value string
}

func (bg *BaseGeometry) FromWKT(wkt string) (*geos.Geometry, error) {
	return geos.FromWKT(wkt)
}
