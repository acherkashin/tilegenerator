package primitives

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
)

type Primitive interface {
	Render(svg *svg.SVG, object *entities.MapObject)
}

func NewPrimitive(t string, params map[string]interface{}) (Primitive, error) {
	switch strings.ToUpper(t) {
	case "TEXT":
		return NewTextPrimitive(&params)
	case "IMAGE":
		return NewImagePrimitive(&params)
	case "POLYLINE":
		return NewPolylinePrimitive(&params)
	case "ARROW":
		return NewArrowPrimitive(&params)

	default:
		return nil, errors.New(fmt.Sprintf("Unknown primitive type %s.", t))
	}
}

func setDefaulColorIfNeed(object *entities.MapObject) {
	if object.ColorOuter == "" {
		object.ColorOuter = "black"
	}

	if object.ColorInner == "" {
		object.ColorInner = "black"
	}

}
