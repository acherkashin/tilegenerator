package entities

import (
	"fmt"

	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID                         int
	StyleName                  string
	TypeID                     int
	Label                      string
	Position                   string
	Size                       int
	IsAntenna                  bool
	NeedShowAzimuthalGrid      bool
	Geometry                   geometry.Geometry
	BeamWidth                  float64
	Sidelobes                  float64
	Azimut                     float64
	Distance                   float64
	NeedShowDirectionalDiagram bool
	ColorOuter                 string
	ColorInner                 string
	Hash                       string
	Code                       string
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeId int, wkt string, isAntenna, needShowAzimuthalGrid, needShowDirectionalDiagram bool, beamWidth, sidelobes, azimut, distance float64, colorOuter, colorInner, code string) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &MapObject{
		ID:                         id,
		TypeID:                     typeId,
		Geometry:                   geo,
		IsAntenna:                  isAntenna,
		NeedShowAzimuthalGrid:      needShowAzimuthalGrid,
		BeamWidth:                  beamWidth,
		Sidelobes:                  sidelobes,
		Azimut:                     azimut,
		Distance:                   distance,
		NeedShowDirectionalDiagram: needShowDirectionalDiagram,
		ColorOuter:                 colorOuter,
		ColorInner:                 colorInner,
		Code:                       code}, nil
}
