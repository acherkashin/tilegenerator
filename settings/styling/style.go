package styling

import (
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
	"errors"
	"fmt"
	"github.com/TerraFactory/tilegenerator/settings/styling/primitives"
	"github.com/TerraFactory/tilegenerator/database/entities"
)

type Primitive interface {
	Render(svg *svg.SVG, geo geometry.Geometry)
}

type Style struct {
	GeometryType int
	Name         string
	Primitives   [] Primitive
}

func (s *Style) Render(object *entities.MapObject, canvas *svg.SVG) {
	for _, p := range s.Primitives{
		p.Render(canvas, object.Geometry)
	}
}

func NewPrimitive(t string, params map[string]interface{}) (Primitive, error) {
	switch t {
	case "TEXT":
		return primitives.NewTextPrimitive(&params)
	case "IMAGE":
		return primitives.ImagePrimitive{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown primitive type %s.", t))
	}
}