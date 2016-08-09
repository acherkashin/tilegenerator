package entities

import (
	"fmt"

	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID            int
	StyleName     string
	TypeID        int
	Label         string
	Code          string
	Position      string
	AzimuthalGrid AzimuthalGrid
	View          View
	Geometry      geometry.Geometry
}

type AzimuthalGrid struct {
	BeamWidth                  float64
	Sidelobes                  float64
	Azimut                     float64
	IsAntenna                  bool
	NeedShowAzimuthalGrid      bool
	NeedShowDirectionalDiagram bool
}

type View struct {
	ColorOuter           string
	ColorInner           string
	NeedMirrorReflection bool
	Scale                float64
	Size                 int
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeId int, wkt string, azimuthalGrid AzimuthalGrid, view View, code string) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &MapObject{
		ID:            id,
		TypeID:        typeId,
		Geometry:      geo,
		AzimuthalGrid: azimuthalGrid,
		View:          view,
		Code:          code}, nil
}
