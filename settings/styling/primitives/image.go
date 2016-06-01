package primitives

import (
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
	"math"
	"strings"
)

type ImagePrimitive struct {
	Width  int64
	Height int64
	Href   string
}

func (img ImagePrimitive) Render(svg *svg.SVG, geo geometry.Geometry) {
	point, _ := geo.AsPoint()
	svg.Image(
		int(math.Floor(point.Coordinates.X + .5)),
		int(math.Floor(point.Coordinates.Y + 0.5)),
		int(img.Width), int(img.Height), img.Href)
}

func (text ImagePrimitive) SetParam(key string, value interface{}) {
	switch strings.ToUpper(key) {
	case "WIDTH":
		text.Width = value.(int64)
	case "HEIGHT":
		text.Height = value.(int64)
	case "HREF":
		text.Href = value.(string)
	}
}
