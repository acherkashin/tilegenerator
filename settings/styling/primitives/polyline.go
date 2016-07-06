package primitives

import (
	"fmt"
	"strings"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
)

type PolylinePrimitive struct {
	Width     int64
	Stroke    string
	DashStyle string
	End       Primitive
	EndId     string
}

func (line PolylinePrimitive) Render(svg *svg.SVG, object *entities.MapObject) {
	setDefaulColorIfNeed(object)

	if strings.Contains(line.Stroke, "${stroke}") {
		line.Stroke = strings.Replace(line.Stroke, "${stroke}", object.ColorOuter, 1)
	}

	linePoints, _ := object.Geometry.AsLineString()
	coords := linePoints.Coordinates
	xs := []int{}
	ys := []int{}

	for i := 0; i <= len(coords)-1; i++ {
		xs = append(xs, int(coords[i].X))
		ys = append(ys, int(coords[i].Y))
	}

	line.End.Render(svg, object)

	style := fmt.Sprintf(
		"stroke: %v; stroke-width: %v; stroke-dasharray: %v; fill: none;marker-end: url(#%v)",
		line.Stroke,
		line.Width,
		line.DashStyle,
		line.EndId)
	svg.Polyline(xs, ys, style)
	// `marker-start="url(#dot)"`,
	// `marker-mid="url(#arrow)"`,)
}

func NewPolylinePrimitive(params *map[string]interface{}) (PolylinePrimitive, error) {
	line := PolylinePrimitive{}
	for key, value := range *params {
		switch strings.ToUpper(key) {
		case "WIDTH":
			line.Width = value.(int64)
		case "STROKE":
			line.Stroke = value.(string)
		case "DASH_STYLE":
			line.DashStyle = value.(string)
		case "END":
			{
				end := value.(map[string]interface{})
				line.End, _ = NewPrimitive(end["Type"].(string), end)
				line.EndId = end["Id"].(string)
				// line.End = line_end
				// if err != nil {
				// 	line.End = nil
				// }

				// fmt.Println(err)

			}
		}
	}
	return line, nil
}
