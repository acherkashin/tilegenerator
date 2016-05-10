package geo

import (
	"github.com/paulsmith/gogeos/geos"
)

// BaseGeometry is a geometry structure
type BaseGeometry struct {
	ID    int
	Value string
}

// FromWKT parses WKT into a structure
func (bg *BaseGeometry) FromWKT(wkt string) (*geos.Geometry, error) {
	return geos.FromWKT(wkt)
}
