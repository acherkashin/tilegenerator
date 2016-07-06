package primitives

import (
	"fmt"
	"strings"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
)

type ArrowPrimitive struct {
	Id       string
	Width    int64
	Height   int64
	Stroke   string
	Fill     string
	Rotate   float64
	Position string
}

func (arrow ArrowPrimitive) Render(svg *svg.SVG, object *entities.MapObject) {
	// point, _ := object.Geometry.AsPoint()
	setDefaulColorIfNeed(object)

	svg.Def()
	svg.Marker(arrow.Id, 0, int(arrow.Height/2), int(arrow.Width), int(arrow.Height), "orient=\"auto\"")

	var xs = []int{0, int(arrow.Width), 0}
	var ys = []int{0, int(arrow.Height / 2), int(arrow.Height)}

	svg.Polyline(xs, ys, createStyle(arrow, object))
	svg.MarkerEnd()
	svg.DefEnd()

}

func createStyle(arrow ArrowPrimitive, object *entities.MapObject) string {
	if strings.Contains(arrow.Stroke, "${stroke}") {
		arrow.Stroke = strings.Replace(arrow.Stroke, "${stroke}", object.ColorOuter, 1)
	}

	if strings.Contains(arrow.Fill, "${fill}") {
		arrow.Stroke = strings.Replace(arrow.Stroke, "${fill}", object.ColorInner, 1)
	}

	style := fmt.Sprintf("stroke: %v; fill: %v", arrow.Stroke, arrow.Fill)

	return style
}

func NewArrowPrimitive(params *map[string]interface{}) (ArrowPrimitive, error) {
	arrow := ArrowPrimitive{}
	for key, value := range *params {
		switch strings.ToUpper(key) {
		case "ID":
			arrow.Id = value.(string)
		case "HEIGHT":
			arrow.Height = value.(int64)
		case "WIDTH":
			arrow.Width = value.(int64)
		case "POSITION":
			arrow.Position = value.(string)
		case "FILL":
			arrow.Fill = value.(string)
		case "STROKE":
			arrow.Stroke = value.(string)
		case "ROTATE":
			arrow.Rotate = value.(float64)
		}
	}

	return arrow, nil
}
