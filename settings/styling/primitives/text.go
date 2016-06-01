package primitives

import (
	"math"
	"strings"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
)

type TextPrimitive struct {
	Weight   int64
	Size     int64
	Style    string
	Position string
	Content  string
}

func (text TextPrimitive) Render(svg *svg.SVG, geo geometry.Geometry) {
	point, _ := geo.AsPoint()
	svg.Text(
		int(math.Floor(point.Coordinates.X + .5)),
		int(math.Floor(point.Coordinates.Y + .5)),
		text.Content)
}

func (text TextPrimitive) SetParam(key string, value interface{}) {
	switch strings.ToUpper(key) { // Switch here is temporary workaround. I should use reflect instead.
	case "SIZE":
		text.Size = value.(int64)
	case "WEIGHT":
		text.Weight = value.(int64)
	case "STYLE":
		text.Style = value.(string)
	case "POSITION":
		text.Position = value.(string)
	case "CONTENT":
		text.Content = value.(string)
	}
}
