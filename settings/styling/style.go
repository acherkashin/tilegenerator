package styling

import (
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings/styling/primitives"
)

type Style struct {
	GeometryType int
	Name         string
	Primitives   []primitives.Primitive
}

func (style *Style) ShouldRender(object *entities.MapObject) bool {
	return style.GeometryType == object.Geometry.GetType() && style.Name == object.StyleName
}

func (s *Style) Render(object *entities.MapObject, canvas *svg.SVG) {
	for _, p := range s.Primitives {
		p.Render(canvas, object)
	}
}
