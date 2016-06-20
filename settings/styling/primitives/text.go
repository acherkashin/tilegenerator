package primitives

import (
	"strings"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
)

type TextPrimitive struct {
	Weight   int64
	Size     int64
	Style    string
	Position string
	Content  string
}

func (text TextPrimitive) Render(svg *svg.SVG, object *entities.MapObject) {
	point, _ := object.Geometry.AsPoint()
	//Temporary solution. I used static shift value because we need to know image size, but we don't know it
	// we have to move such primitives as labels and so on into the image primitive
	var xShift float64
	var yShift float64

	if object.Label != "" {
		text.Content = strings.Replace(text.Content, "${label}", object.Label, 1)
	} else {
		text.Content = ""
	}

	switch text.Position {
	case "top":
		xShift = -20
		yShift = -40
	case "bottom":
		xShift = -20
		yShift = 40
	case "left":
		xShift = -60
		yShift = 10
	case "right":
		xShift = 20
		yShift = 10
	}

	svg.Text(
		int(point.Coordinates.X+xShift),
		int(point.Coordinates.Y+yShift),
		text.Content)
}

func NewTextPrimitive(params *map[string]interface{}) (TextPrimitive, error) {
	text := TextPrimitive{}
	for key, value := range *params {
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

	return text, nil
}
