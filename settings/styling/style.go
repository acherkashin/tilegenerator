package styling

import (
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
	"errors"
	"fmt"
	"github.com/TerraFactory/tilegenerator/settings/styling/primitives"
)

type Primitive interface {
	Render(svg *svg.SVG, geo geometry.Geometry)
	SetParam(name string, value interface{})
}

type Style struct {
	GeometryType int
	Name         string
	Primitives   [] Primitive
}

func (s *Style) Render(canvas *svg.SVG) {

}

func NewPrimitive(t string) (Primitive, error) {
	switch t {
	case "TEXT":
		return primitives.TextPrimitive{}, nil
	case "IMAGE":
		return primitives.ImagePrimitive{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown primitive type %s.", t))
	}
}