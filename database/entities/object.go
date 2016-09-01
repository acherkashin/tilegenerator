package entities

import (
	"fmt"

	"github.com/TerraFactory/wktparser"
	"github.com/TerraFactory/wktparser/geometry"
)

//define markerPosition constans
const (
	LeftUp = 1 + iota
	Up
	RightUp
	Left
	Center
	Right
	LeftDown
	Down
	RightDown
)

// MapObject represents a geometry on a map
type MapObject struct {
	ID             int
	StyleName      string
	TypeID         int
	Label          string
	Code           string
	MarkerPosition int
	Position       string
	AzimuthalGrid  AzimuthalGrid
	View           View
	Geometry       geometry.Geometry
}

//AzimuthalGrid - information about azimuthalGrid for  map object
type AzimuthalGrid struct {
	BeamWidth                  float64
	Sidelobes                  float64
	Azimut                     float64
	IsAntenna                  bool
	NeedShowAzimuthalGrid      bool
	NeedShowDirectionalDiagram bool
}

//View of map object
type View struct {
	ColorOuter           string
	ColorInner           string
	Scale                float64
	Size                 int
	UseCurveBezier       bool
	NeedMirrorReflection bool
}

// NewObject creates new MapObject with a parsed from WKT geometry
func NewObject(id int, typeId, markerPosition int, wkt string, code string, azimuthalGrid AzimuthalGrid, view View) (*MapObject, error) {
	geo, err := wktparser.Parse(wkt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if !isValidMarkerPosition(markerPosition) {
		markerPosition = Center
	}

	return &MapObject{
		ID:             id,
		TypeID:         typeId,
		Geometry:       geo,
		AzimuthalGrid:  azimuthalGrid,
		MarkerPosition: markerPosition,
		View:           view,
		Code:           code}, nil
}

//NewView creates new view for map object
func NewView(colorOuter, colorInner string, needMirrorReflection, useCurveBezier bool, scale float64) *View {
	return &View{
		ColorOuter:           colorOuter,
		ColorInner:           colorInner,
		NeedMirrorReflection: needMirrorReflection,
		UseCurveBezier:       useCurveBezier,
		Scale:                scale,
		// Size:                 size
	}
}

//NewAzimuthalGrid creates new object-information about azimuthalGrid for map object
func NewAzimuthalGrid(beamWidth, sidelobes, azimut float64,
	isShortwaveAntenna, needShowAzimuthalGrid, needShowDirectionalDiagram bool) *AzimuthalGrid {

	return &AzimuthalGrid{
		BeamWidth:                  beamWidth,
		Sidelobes:                  sidelobes,
		Azimut:                     azimut,
		IsAntenna:                  isShortwaveAntenna,
		NeedShowAzimuthalGrid:      needShowAzimuthalGrid,
		NeedShowDirectionalDiagram: needShowDirectionalDiagram}
}

func isValidMarkerPosition(position int) bool {
	return position > 0 && position < 10
}
