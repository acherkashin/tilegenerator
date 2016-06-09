package styling

import (
	"errors"
	"fmt"
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings/styling/primitives"
	"github.com/TerraFactory/wktparser/geometry"
)

type Primitive interface {
	Render(svg *svg.SVG, geo geometry.Geometry, object *entities.MapObject)
}

type Style struct {
	GeometryType int
	Name         string
	Primitives   []Primitive
}

func (style *Style) ShouldRender(object *entities.MapObject) bool {
	return style.GeometryType == object.Geometry.GetType() && style.Name == object.StyleName
}

func (s *Style) Render(object *entities.MapObject, canvas *svg.SVG) {
	for _, p := range s.Primitives {
		p.Render(canvas, object.Geometry, object)
	}
}

func NewPrimitive(t string, params map[string]interface{}) (Primitive, error) {
	switch t {
	case "TEXT":
		return primitives.NewTextPrimitive(&params)
	case "IMAGE":
		return primitives.NewImagePrimitive(&params)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown primitive type %s.", t))
	}
}
