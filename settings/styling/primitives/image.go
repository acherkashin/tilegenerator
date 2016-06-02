package primitives

import (
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
	"math"
	"strings"
	"fmt"
)

type ImagePrimitive struct {
	Width  int64
	Height int64
	Href   string
}

func (img ImagePrimitive) Render(svg *svg.SVG, geo geometry.Geometry) {
	point, _ := geo.AsPoint()
	fmt.Println(img.Href)
	svg.Image(
		int(math.Floor(point.Coordinates.X + .5)),
		int(math.Floor(point.Coordinates.Y + 0.5)),
		int(img.Width), int(img.Height), img.Href)
}

func NewImagePrimitive(params *map[string]interface{}) (ImagePrimitive, error) {
	text := ImagePrimitive{}
	for key, value := range *params {
		switch strings.ToUpper(key) { // Switch here is temporary workaround. I should use reflect instead.
		case "WIDTH":
			text.Width = value.(int64)
		case "HEIGHT":
			text.Height = value.(int64)
		case "HREF":
			text.Href = value.(string)
		}
	}

	return text, nil
}