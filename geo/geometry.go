package geo

import (
	"github.com/paulsmith/gogeos/geos"
)

func FromWKT(wkt string) (*geos.Geometry, error) {
	return geos.FromWKT(wkt)
}
