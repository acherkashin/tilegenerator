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
	NeedMirrorReflection       bool
	ColorOuter                 string
	ColorInner                 string
	Code                       string
	Scale                      float64
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeId int, wkt string, isAntenna, needShowAzimuthalGrid, needShowDirectionalDiagram, needMirrorReflection bool, beamWidth, sidelobes, azimut, distance float64, colorOuter, colorInner, code string, scale float64) (*MapObject, error) {
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
		NeedMirrorReflection:       needMirrorReflection,
		ColorOuter:                 colorOuter,
		ColorInner:                 colorInner,
		Code:                       code,
		Scale:                      scale}, nil
}
