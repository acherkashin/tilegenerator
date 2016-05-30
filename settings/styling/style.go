package styling

import (
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
)

type Primitive interface {
	Render(*svg.SVG, *geometry.Geometry)
}

type Style struct {
	GeometryType int
	Name string
	Primitives[] Primitive
}

func (s *Style) Render(canvas *svg.SVG) {

}