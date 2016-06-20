package entities

import (
	"fmt"

	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID                    int
	StyleName             string
	TypeID                int
	Label                 string
	Position              string
	Size                  int
	IsAntenna             bool
	NeedShowAzimuthalGrid bool
	Geometry              geometry.Geometry
	BeamWidth             float64
	Sidelobes             float64
	Azimut                float64
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeId int, wkt string, isAntenna, needShowAzimuthalGrid bool, beamWidth, sidelobes, azimut float64) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &MapObject{
		ID:                    id,
		TypeID:                typeId,
		Geometry:              geo,
		IsAntenna:             isAntenna,
		NeedShowAzimuthalGrid: needShowAzimuthalGrid,
		BeamWidth:             beamWidth,
		Sidelobes:             sidelobes,
		Azimut:                azimut}, nil
}
