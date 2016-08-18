package tiles

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings"
	"github.com/TerraFactory/tilegenerator/settings/styling"
	"github.com/TerraFactory/tilegenerator/utils"
	"github.com/TerraFactory/wktparser/geometry"
)

var routeAviationsFlightCodes = []string{
	"121062001002010601010100000100",
	"1210620050020106010002000100",
	"1210620050030106010002000100",
	"1210620050030106010003000100",
	"12106200800201060102000100",
	"1210620050020106010003000100",
	"1000000004",
	"121062003002010601000300010100",
	"121062606030020801030001000000",
	"121062010002010601000300010100",
	"121062009002010601000300010100",
}

var patrollingAreaCodes = []string{
	"912106200800201060102000100",
	"9121062001002010601010100000100",
	"9121062606030020801030001000000",
	"91210620050020106010003000100",
	"91210620050020106010002000100",
	"9121062003002010601000300010100",
	"91210620050030106010002000100",
	"91210620050030106010003000100",
	"9121062010002010601000300010100",
	"9121062009002010601000300010100",
	"1000000002",
}

var attackMainDirectionCode = "009995004501010106"
var plannedAttackMainDirectionCode = "13147260003502080100010003"
var completedProvideActionCode = "13147260003502080100010002"
var pitCode = "21056441215588"

var hashTypes map[int]string

// RenderTile takes a tile struct, map objects and then draws these objects on the tile
func RenderTile(tile *Tile, objects *[]entities.MapObject, styles *map[string]styling.Style, writer io.Writer) {
	f := func(x, y float64) (float64, float64) {
		nx, ny := tile.Degrees2Pixels(y, x)
		return float64(nx), float64(ny)
	}

	canvas := svg.New(writer)
	canvas.Start(TileSize, TileSize)
	for _, object := range *objects {
		object.Geometry.ConvertCoords(f)

		for _, style := range *styles {

			if style.ShouldRender(&object) {
				style.Render(&object, canvas)

				if object.AzimuthalGrid.IsAntenna && object.AzimuthalGrid.NeedShowDirectionalDiagram {
					RenderBeamDiagram(canvas, &object, tile)
				}

				if object.AzimuthalGrid.NeedShowAzimuthalGrid {
					RenderAzimuthalGrid(canvas, &object, tile)
				}
			}
		}

		if contains(object.Code, patrollingAreaCodes) {
			RenderPatrollingArea(canvas, &object, tile)
		} else if contains(object.Code, routeAviationsFlightCodes) {
			RenderRouteAviationFlight(canvas, &object, tile)
		} else if object.Code == plannedAttackMainDirectionCode {
			RenderPlannedAttackMainDirection(canvas, &object, tile)
		} else if object.Code == attackMainDirectionCode {
			RenderAttackMainDirection(canvas, &object, tile)
		} else if object.Code == completedProvideActionCode {
			RenderCompletedProvideAction(canvas, &object, tile)
		} else if object.Code == pitCode {
			RenderPit(canvas, &object, tile)
		}
	}
	canvas.End()
}

type chartPoint struct {
	x, y, z, value float64
}

type bigArrow struct {
	pointXs, pointYs, arrowXs, arrowYs []int
	centerX, centerY                   int

	rotateAngel float64
}

func newBigArrow(tile *Tile, coords []geometry.Coord) *bigArrow {
	x1, y1 := int(coords[0].X), int(coords[0].Y)
	x2, y2 := int(coords[1].X), int(coords[1].Y)

	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	sizeEdgeArrow := int(distance / 10)

	line1X1, line1Y1 := centerX+int(distance/2), centerY+sizeEdgeArrow/2/3
	line1X2, line1Y2 := centerX-int(distance/2), centerY+sizeEdgeArrow/2

	line2X1, line2Y1 := centerX+int(distance/2), centerY-sizeEdgeArrow/2/3
	line2X2, line2Y2 := centerX-int(distance/2), centerY-sizeEdgeArrow/2

	rightLinePointX, rightLinePointY := centerX+int(distance/2), centerY

	arrowXs := []int{rightLinePointX, rightLinePointX + int(math.Sqrt(3)/2.0*float64(sizeEdgeArrow)), rightLinePointX}
	arrowYs := []int{rightLinePointY + sizeEdgeArrow/2, rightLinePointY, rightLinePointY - sizeEdgeArrow/2}

	pointXs := []int{line1X2, line1X1, arrowXs[0], arrowXs[1], arrowXs[2], line2X1, line2X2}
	pointYs := []int{line1Y2, line1Y1, arrowYs[0], arrowYs[1], arrowYs[2], line2Y1, line2Y2}

	return &bigArrow{
		centerX:     centerX,
		centerY:     centerY,
		pointXs:     pointXs,
		pointYs:     pointYs,
		arrowXs:     arrowXs,
		arrowYs:     arrowYs,
		rotateAngel: getAngel(x1, y1, x2, y2) + 180,
	}
}

func newCharPoint(x, y, z, value float64) *chartPoint {
	var cp chartPoint
	cp.x = x
	cp.y = y
	cp.z = z
	cp.value = value

	return &cp
}

func getBeamDiagramPoints(centerX, centerY, beamWidth int, sidelobes, radius float64) (xs, ys []int) {
	tempPointList := getTempPointsForBeamPoints(beamWidth, sidelobes)
	maxValue := getMax(tempPointList)

	for _, point := range tempPointList {
		xs = append(xs, centerX+int(point.y*radius/maxValue))
		ys = append(ys, centerY+int(point.x*radius/maxValue))
	}

	xs[0] = centerX + int(radius)
	ys[0] = centerY

	return xs, ys
}

func distanceBeetweenPoints(x1, y1, x2, y2 int) float64 {
	a := x1 - x2
	b := y1 - y2

	return float64(math.Sqrt(float64(a*a + b*b)))
}

func distanceBeetweenPointsFloat(x1, y1, x2, y2 float64) float64 {
	a := x1 - x2
	b := y1 - y2

	return float64(math.Sqrt(float64(a*a + b*b)))
}

func getLineCenter(x1, y1, x2, y2 int) (x, y int) {
	x = (x2 + x1) / 2
	y = (y2 + y1) / 2

	return x, y
}

func getAngel(x1, y1, x2, y2 int) float64 {
	r := distanceBeetweenPoints(x1, y1, x2, y2)
	angel := math.Acos((float64(x1)-float64(x2))/r) * 180.0 / math.Pi

	if ((angel < 0) && (y1 > y2)) || ((angel > 0) && (y1 < y2)) {
		angel *= -1
	}

	return angel
}

func getArrowPoints(x1, y1, dx int) (sx []int, sy []int) {
	sx = append(sx,
		x1+dx,
		x1,
		x1+dx)

	sy = append(sy,
		y1+dx/2,
		y1,
		y1-dx/2)

	return sx, sy
}

func getTempPointsForBeamPoints(k int, n float64) []*chartPoint {
	tempPointList := make([]*chartPoint, 0, 721) //нормализация

	for i := 0; i <= 720; i++ {
		angelRadian := (float64(i) * (math.Pi / 360)) //перевод в радианы

		value := сalculateValueFormula(k, angelRadian, n)
		x := value * math.Sin(angelRadian)
		y := value * math.Cos(angelRadian)

		tempPointList = append(tempPointList, newCharPoint(x, y, angelRadian, value))
	}

	return tempPointList
}

/*
 "F(q) = Sin(K * Pi * Sin(q)) / Sin(Pi * Sin(q)/ (5.1 - 4 * N)))"
 "k" - beam
 "q" - угол
 "n" - sidelobes(боковые лепестки)
*/
func сalculateValueFormula(k int, q, n float64) float64 {
	value1 := math.Sin(float64(k) * math.Pi * math.Sin(q))
	value2 := math.Sin(math.Pi * math.Sin(q) / (5.1 - 4*n))

	var val float64

	if value2 != 0 {
		val = value1 / value2
	} else {
		val = 0
	}

	return val
}

func getRightArrowPoints(x1, y1, distance int) (xs []int, ys []int) {
	radiusX := distance / 4
	radiusY := radiusX / 2
	arrowSize := 10

	pointForArrowX := x1
	pointForArrowY := y1 + int(2*radiusY)
	xs, ys = getArrowPoints(pointForArrowX, pointForArrowY, arrowSize)

	return xs, ys
}

func getLeftArrowPoints(x1, y1, distance int) (xs []int, ys []int) {
	radiusX := distance / 4
	radiusY := radiusX / 2
	arrowSize := 10

	pointForArrowX := x1
	pointForArrowY := y1 - int(2*radiusY)
	xs, ys = getArrowPoints(pointForArrowX, pointForArrowY, -arrowSize)

	return xs, ys
}

func getMax(chartPoints []*chartPoint) float64 {
	max := math.SmallestNonzeroFloat64

	for _, point := range chartPoints {
		if point.value > max {
			max = point.value
		}
	}

	return max
}

func RenderNewObject(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}
	setDefaultColor(object)

	styleLine := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;stroke-dasharray: 10 2 2 2;", object.View.ColorOuter, 2)

	renderCurveOrPolyline(canvas, object.View.UseCurveBezier, styleLine, line.Coordinates)

	xs, ys := coordToXsYs(line.Coordinates)
	count := len(xs)

	renderRightPartNewObject(canvas, xs[0], ys[0], xs[1], ys[1], object.View.ColorInner, object.View.ColorOuter)
	renderLeftPartNewObject(canvas, xs[count-2], ys[count-2], xs[count-1], ys[count-1], object.View.ColorInner, object.View.ColorOuter)

	return nil
}

func renderRightPartNewObject(canvas *svg.SVG, x1, y1, x2, y2 int, colorInner, colorOuter string) {
	radiusX := 10
	radiusY := radiusX / 2
	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	rightLinePointX, rightLinePointY := centerX+int(distance/2), centerY

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", getAngel(x1, y1, x2, y2), centerX, centerY))

	canvas.Arc(rightLinePointX,
		rightLinePointY,
		radiusX,
		radiusY,
		0, false, true,
		rightLinePointX+2*radiusX,
		rightLinePointY,
		fmt.Sprintf("stroke: %v; fill: none;stroke-width: %v;", colorOuter, 2))

	canvas.Gend()
}

func renderLeftPartNewObject(canvas *svg.SVG, x1, y1, x2, y2 int, colorInner, colorOuter string) {
	radiusX := 10
	radiusY := radiusX / 2
	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	leftLinePointX, leftLinePointY := centerX-int(distance/2), centerY

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", getAngel(x1, y1, x2, y2), centerX, centerY))

	canvas.Arc(leftLinePointX-2*radiusX,
		leftLinePointY,
		radiusX,
		radiusY,
		0, false, true,
		leftLinePointX,
		leftLinePointY,
		fmt.Sprintf("stroke: %v; fill: none;stroke-width: %v;", colorOuter, 2))

	canvas.Gend()
}

func RenderPit(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	canvas.Gid(fmt.Sprintf("id%v", strconv.Itoa(object.ID)))
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}

	setDefaultColor(object)
	style := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.View.ColorOuter, 2)

	if object.View.UseCurveBezier {
		renderCurveBezierWithHatching(canvas, line.Coordinates, style)
	} else {
		renderPolygonWithHatching(canvas, line.Coordinates, style)
	}
	canvas.Gend()
	return nil
}

func renderCurveBezierWithHatching(canvas *svg.SVG, coords []geometry.Coord, style string) {
	xs, ys := coordToXsYs(coords)
	percentLength := 0.5
	renderCurve(canvas, coords, style)

	if len(xs) >= 3 {
		for i := 0; i <= len(xs)-3; i++ {
			beginArcX, beginArcY := utils.GetPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], 1-percentLength)
			endArcX, endArcY := utils.GetPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], percentLength)

			if i == 0 {
				//render first line
				drawHatchingOnLine(canvas, xs[i], ys[i], beginArcX, beginArcY, style)
			} else {
				x, y := utils.GetPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], percentLength)
				drawHatchingOnLine(canvas, x, y, beginArcX, beginArcY, style)
			}

			renderHatchingOnBezier(canvas, beginArcX, beginArcY, xs[i+1], ys[i+1], endArcX, endArcY, style)

			if i == len(xs)-3 {
				//render last line
				drawHatchingOnLine(canvas, endArcX, endArcY, xs[i+2], ys[i+2], style)
			} else {
				x, y := utils.GetPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], 1-percentLength)
				drawHatchingOnLine(canvas, endArcX, endArcY, x, y, style)
			}
		}
	} else {
		drawHatchingOnLine(canvas, xs[0], ys[0], xs[1], ys[1], style)
	}
}

func renderPolygonWithHatching(canvas *svg.SVG, coords []geometry.Coord, style string) {
	for i := 0; i < len(coords)-1; i++ {
		canvas.Line(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), style)
		drawHatchingOnLine(canvas, int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), style)
	}

}

func renderHatchingOnBezier(canvas *svg.SVG, beginX, beginY, controlX, controlY, endX, endY int, style string) {
	bezierXs, bezierYs := bezierToPolyline(beginX, beginY, controlX, controlY, endX, endY, 0.1)
	lengthHatch := 8.0
	percentSizeSegment := lengthHatch / getLengthPolyline(bezierXs, bezierYs)
	if percentSizeSegment > 1 {
		percentSizeSegment = 1
	}

	for j := 0.; j < 1; j += percentSizeSegment {
		beginHatchX, beginHatchY := getPointBezier(float64(beginX), float64(beginY), float64(controlX), float64(controlY), float64(endX), float64(endY), j)

		nextX, nextY := getPointBezier(float64(beginX), float64(beginY), float64(controlX), float64(controlY), float64(endX), float64(endY), j+percentSizeSegment/5000)
		distance := distanceBeetweenPointsFloat(beginHatchX, beginHatchY, nextX, nextY)
		nextX, nextY = utils.GetPointOnLineFloat(beginHatchX, beginHatchY, nextX, nextY, lengthHatch/distance)
		endHatchX, endHatchY := utils.RotatePoint(int(beginHatchX), int(beginHatchY), int(nextX), int(nextY), -90)

		//sometimes endHatchX and endHatchY are infinity
		if math.Abs(float64(endHatchX)) > 9223372036 || math.Abs(float64(endHatchY)) > 9223372036 {
			continue
		}
		canvas.Line(int(beginHatchX), int(beginHatchY), int(endHatchX), int(endHatchY), style)

		nextX, nextY = getPointBezier(float64(beginX), float64(beginY), float64(controlX), float64(controlY), float64(endX), float64(endY), j+percentSizeSegment)

		//if distance between current and next point less than lengthHatch/2 then skip next point
		if distanceBeetweenPoints(int(beginHatchX), int(beginHatchY), int(nextX), int(nextY)) < lengthHatch/2 {
			j += percentSizeSegment
		}
	}

}
func drawHatchingOnLine(canvas *svg.SVG, beginX, beginY, endX, endY int, style string) {
	length := 8
	distance := distanceBeetweenPoints(beginX, beginY, endX, endY)
	percentSizeSegment := float64(length) / distance

	if percentSizeSegment >= 1 {
		percentSizeSegment = 1
	}

	count := int(distance) / length

	currentPointPercent := .0
	for i := 0; i <= count; i++ {
		x1, y1 := utils.GetPointOnLine(beginX, beginY, endX, endY, currentPointPercent)
		currentPointPercent += percentSizeSegment
		x2, y2 := utils.GetPointOnLine(beginX, beginY, endX, endY, currentPointPercent)

		resultX, resultY := utils.RotatePoint(x1, y1, x2, y2, -90)
		canvas.Line(x1, y1, resultX, resultY, style)
	}

	//draw last hatch
	x1, y1 := utils.GetPointOnLine(beginX, beginY, endX, endY, 1-percentSizeSegment)
	resultX, resultY := utils.RotatePoint(endX, endY, x1, y1, 90)
	canvas.Line(endX, endY, resultX, resultY, style)
}

func RenderAttackMainDirection(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	if object.View.ColorInner == "" {
		object.View.ColorInner = "red"
	}

	if object.View.ColorOuter == "" {
		object.View.ColorOuter = "red"
	}

	weight := 1

	styleLine := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.View.ColorOuter, weight)
	styleArrow := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: %v;", object.View.ColorOuter, weight, object.View.ColorInner)

	return renderBigArrow(canvas, object, tile, styleLine, styleArrow)
}

func RenderPlannedAttackMainDirection(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	if object.View.ColorInner == "" {
		object.View.ColorInner = "red"
	}

	if object.View.ColorOuter == "" {
		object.View.ColorOuter = "red"
	}

	weight := 1

	styleLine := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;stroke-dasharray: 10;", object.View.ColorOuter, weight)
	styleArrow := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: %v;", object.View.ColorOuter, weight, object.View.ColorInner)

	return renderBigArrow(canvas, object, tile, styleLine, styleArrow)
}

func RenderCompletedProvideAction(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	weight := 1
	if object.View.ColorOuter != "" {
		object.View.ColorOuter = "red"
	}

	style := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.View.ColorOuter, weight)

	return renderBigArrow(canvas, object, tile, style, style)
}

func renderBigArrow(canvas *svg.SVG, object *entities.MapObject, tile *Tile, styleLine, styleArrow string) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}

	direction := newBigArrow(tile, line.Coordinates)

	transformationStyle := fmt.Sprintf("rotate(%v,%v,%v)", direction.rotateAngel, direction.centerX, direction.centerY)

	canvas.Gtransform(transformationStyle)

	canvas.Polyline(direction.arrowXs, direction.arrowYs, styleArrow)
	canvas.Polyline(direction.pointXs, direction.pointYs, styleLine)

	canvas.Gend()

	return nil
}

// RenderRouteAviationFlight renders an aviation route on a tile
func RenderRouteAviationFlight(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}

	coords := line.Coordinates
	setDefaultColor(object)

	canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")

	style := fmt.Sprintf("stroke: %v; stroke-width: %v; fill: none; stroke-dasharray: 10;", object.View.ColorOuter, 1)

	renderCurveOrPolyline(canvas, object.View.UseCurveBezier, style, coords)

	if object.Code != "1000000004" {
		xs, ys := coordToXsYs(coords)
		if object.View.UseCurveBezier {
			xs, ys = polylineToCurvePoints(xs, ys)
		}

		x, y, angel := getCenterPolylineAndAngel(xs, ys)
		renderImageOnLine(canvas, object.View.Scale, angel, x, y, object.ID, tile.Z)
	}

	if object.Label != "" {
		xs, ys := coordToXsYs(coords)

		if object.View.UseCurveBezier {
			xs, ys = polylineToCurvePoints(xs, ys)
		}
		x, y, _ := getCenterPolylineAndAngel(xs, ys)
		renderTextOnLine(canvas, x, y, object.Label, object.Position)
	}

	renderArrowRouteAviationFlight(coords, canvas, object, tile)

	canvas.Gend()
	return nil
}

func renderPolyline(canvas *svg.SVG, coords []geometry.Coord, style string) {
	xs, ys := coordToXsYs(coords)

	canvas.Polyline(xs, ys, style)
}

//angel is angel between line and axis Ox
func getCenterPolylineAndAngel(xs, ys []int) (x, y int, angel float64) {
	var lineLength float64
	var i int
	halfLength := getLengthPolyline(xs, ys) / 2

	for i = 0; i < len(xs) && halfLength > 0; i++ {
		lineLength = distanceBeetweenPoints(xs[i], ys[i], xs[i+1], ys[i+1])
		halfLength -= lineLength
	}

	percentPosition := 1 - (-1.0)*halfLength/lineLength

	centerX, centerY := utils.GetPointOnLine(xs[i-1], ys[i-1], xs[i], ys[i], percentPosition)

	angel = getAngel(xs[i-1], ys[i-1], xs[i], ys[i])

	return centerX, centerY, angel
}

func renderArrowRouteAviationFlight(coords []geometry.Coord, canvas *svg.SVG, object *entities.MapObject, tile *Tile) {
	lengthArrow := 5.0
	weight := 1
	i := len(coords) - 2
	lastLineLength := distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))

	styleArrow := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: %v;", object.View.ColorInner, weight, object.View.ColorInner)

	if lastLineLength > lengthArrow {
		xs, ys := getArrowRouteAviationFlight(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), tile.Z)
		canvas.Polyline(xs, ys, styleArrow)
	}
}

/*
RenderPatrollingArea (район барражирования)
*/
func RenderPatrollingArea(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}
	coords := line.Coordinates
	count := len(coords)
	setDefaultColor(object)

	canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")

	xs, ys := coordToXsYs(coords)

	style := fmt.Sprintf("stroke: %v; fill: none;", object.View.ColorOuter)

	renderCurveOrPolyline(canvas, object.View.UseCurveBezier, style, coords)

	if object.Code != "1000000002" {
		if object.View.UseCurveBezier {
			curveXs, curveYs := polylineToCurvePoints(xs, ys)
			x, y, angle := getCenterPolylineAndAngel(curveXs, curveYs)
			renderImageOnLine(canvas, object.View.Scale, angle, x, y, object.ID, tile.Z)
		} else {
			x, y, angle := getCenterPolylineAndAngel(xs, ys)
			renderImageOnLine(canvas, object.View.Scale, angle, x, y, object.ID, tile.Z)
		}
	}

	if object.Label != "" {
		var x, y int
		if object.View.UseCurveBezier {
			x, y, _ = getCenterPolylineAndAngel(xs, ys)
		} else {
			curveXs, curveYs := polylineToCurvePoints(xs, ys)
			x, y, _ = getCenterPolylineAndAngel(curveXs, curveYs)
		}
		renderTextOnLine(canvas, x, y, object.Label, object.Position)
	}

	renderRightPartPatrollingArea(canvas, xs[0], ys[0], xs[1], ys[1], object.View.ColorInner, object.View.ColorOuter)
	renderLeftPartPatrollingArea(canvas, xs[count-2], ys[count-2], xs[count-1], ys[count-1], object.View.ColorInner, object.View.ColorOuter)

	canvas.Gend()

	return nil
}

//At first, we display arc and arrow horizontally and then rotate to needed angle
func renderRightPartPatrollingArea(canvas *svg.SVG, x1, y1, x2, y2 int, colorInner, colorOuter string) {
	radiusX := int(distanceBeetweenPoints(x1, y1, x2, y2) / 4)
	radiusY := int(radiusX / 2)
	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	rightLinePointX, rightLinePointY := centerX+int(distance/2), centerY
	rightArrowXs, rightArrowYs := getRightArrowPoints(rightLinePointX, rightLinePointY, int(distance))

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", getAngel(x1, y1, x2, y2), centerX, centerY))
	canvas.Arc(rightLinePointX,
		rightLinePointY,
		radiusX,
		radiusY,
		0, false, true,
		rightLinePointX,
		rightLinePointY+int(2*radiusY),
		fmt.Sprintf("stroke: %v; fill: none;", colorOuter))

	arrowSize := 10
	if radiusX > arrowSize {
		canvas.Polyline(rightArrowXs, rightArrowYs, fmt.Sprintf("stroke: %v; fill: %v;", colorInner, colorInner))
	}
	canvas.Gend()
}

//At first, we display arc and arrow horizontally and then rotate to needed angle
func renderLeftPartPatrollingArea(canvas *svg.SVG, x1, y1, x2, y2 int, colorInner, colorOuter string) {
	radiusX := int(distanceBeetweenPoints(x1, y1, x2, y2) / 4)
	radiusY := int(radiusX / 2)
	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	leftLinePointX, leftLinePointY := centerX-int(distance/2), centerY
	leftArrowXs, leftArrowYs := getLeftArrowPoints(leftLinePointX, leftLinePointY, int(distance))

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", getAngel(x1, y1, x2, y2), centerX, centerY))
	canvas.Arc(leftLinePointX,
		leftLinePointY,
		radiusX,
		radiusY,
		0, false, true,
		leftLinePointX,
		leftLinePointY-int(2*radiusY),
		fmt.Sprintf("stroke: %v; fill: none", colorOuter))

	arrowSize := 10
	if radiusX > arrowSize {
		canvas.Polyline(leftArrowXs, leftArrowYs, fmt.Sprintf("stroke: %v; fill: %v;", colorInner, colorInner))
	}
	canvas.Gend()
}

func setDefaultColor(object *entities.MapObject) {
	if object.View.ColorInner == "" {
		object.View.ColorInner = "black"
	}
	if object.View.ColorOuter == "" {
		object.View.ColorOuter = "black"
	}
}

//this function is used temporary, till we don't use styles for rendering of primitives
func renderImageOnLine(canvas *svg.SVG, scale, angel float64, x, y, id, zoom int) {
	pathConfig := "./config.toml"
	settings, err := settings.GetSettings(&pathConfig)

	if err == nil {
		href := fmt.Sprintf("%v/api/maps/object/%v/png", settings.UrlAPI, id)
		if result, err := utils.GetImgByURL(href); err == nil {
			imgBase64Str := base64.StdEncoding.EncodeToString(result)

			img2html := "data:image/png;base64," + imgBase64Str

			imageWidth := scale * float64(5+5*zoom)
			imageHeight := scale * float64(7+6*zoom)

			canvas.Image(x-(int)(imageWidth/2.0),
				y-(int)(imageHeight/2.0),
				(int)(imageWidth),
				(int)(imageHeight),
				img2html,
				fmt.Sprintf("transform=\"rotate(%v,%v,%v)\"", angel-90, x, y))
		}
	}
}

//this function is used temporary, till we don't use styles for rendering of primitives
func renderTextOnLine(svg *svg.SVG, x, y int, label, position string) {
	var xShift, yShift int

	if position == "" {
		position = "bottom"
	}

	switch position {
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

	svg.Text(int(x+xShift), int(y+yShift), label)
}

func getLengthPolyline(xs, ys []int) float64 {
	sum := 0.0

	for i := 0; i < len(xs)-1; i++ {
		sum += distanceBeetweenPoints(xs[i], ys[i], xs[i+1], ys[i+1])
	}

	return sum
}

func getArrowRouteAviationFlight(BeginX, BeginY, EndX, EndY, zoom int) ([]int, []int) {
	var angel int
	distance := distanceBeetweenPoints(BeginX, BeginY, EndX, EndY)
	percentSize := 1 - 5.0/distance
	angel = 120
	centerX, centerY := utils.GetPointOnLine(BeginX, BeginY, EndX, EndY, percentSize)
	rotatedPointX, rotatedPointY := EndX, EndY

	p1X, p1Y := utils.RotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, angel)
	p2X, p2Y := utils.RotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, -angel)

	xs := []int{p1X, rotatedPointX, p2X}
	ys := []int{p1Y, rotatedPointY, p2Y}
	return xs, ys
}

/*
RenderBeamDiagram - drawing satellite directional diagram
 "F(q) = Sin(K * Pi * Sin(q)) / Sin(Pi * Sin(q)/ (5.1 - 4 * N)))"
 "k" - beam
 "q" - угол
 "n" - sidelobes(боковые лепестки)
*/
func RenderBeamDiagram(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	point, err := object.Geometry.AsPoint()
	if err != nil {
		return err
	}

	radius := getRadiusAzimuthalGrid(tile.Z)
	strokeWidth := getStrokeWidthAzimuthalGrid(radius)
	centerX, centerY := int(point.Coordinates.X), int(point.Coordinates.Y)

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", object.AzimuthalGrid.Azimut, centerX, centerY))
	xs, ys := getBeamDiagramPoints(centerX, centerY, int(object.AzimuthalGrid.BeamWidth), object.AzimuthalGrid.Sidelobes, radius)

	if object.View.ColorOuter == "" {
		object.View.ColorOuter = "red"
	}

	canvas.Polygon(xs, ys, fmt.Sprintf("stroke:%v; stroke-width:%v; fill: none;", object.View.ColorOuter, strokeWidth))

	canvas.Gend()
	return nil
}

//RenderAzimuthalGrid ...
func RenderAzimuthalGrid(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	point, err := object.Geometry.AsPoint()
	if err != nil {
		return err
	}

	radius := getRadiusAzimuthalGrid(tile.Z)
	strokeWidth := getStrokeWidthAzimuthalGrid(radius)
	centerX, centerY := int(point.Coordinates.X), int(point.Coordinates.Y)
	templateStyle := "stroke:%v; stroke-width:%v; fill: none;"

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", object.AzimuthalGrid.Azimut, centerX, centerY))
	renderGradationAzimuthalGrid(canvas, object, centerX, centerY, tile.Z)
	canvas.Circle(centerX, centerY, int(radius*0.67), fmt.Sprintf(templateStyle, "yellow", strokeWidth))
	canvas.Circle(centerX, centerY, int(radius), fmt.Sprintf(templateStyle, "green", strokeWidth))
	canvas.Gend()

	return nil
}

func renderGradationAzimuthalGrid(canvas *svg.SVG, object *entities.MapObject, centerX, centerY, zoom int) {
	if object.View.ColorInner == "" {
		object.View.ColorInner = "gray"
	}
	radius := getRadiusAzimuthalGrid(zoom)
	strokeWidth := getStrokeWidthAzimuthalGrid(radius)
	styleAzimuthalGrid := fmt.Sprintf("stroke:%v; stroke-width:%v; fill: none;", object.View.ColorInner, strokeWidth)

	polarGridXs, polarGridYs := getPointsForPolarGrid(centerX, centerY, radius)

	for i := 0; i < len(polarGridXs); i++ {
		canvas.Line(centerX, centerY, polarGridXs[i], polarGridYs[i], styleAzimuthalGrid)
	}

	indexZeroAzimut := 18
	canvas.Line(centerX, centerY, polarGridXs[indexZeroAzimut], polarGridYs[indexZeroAzimut],
		fmt.Sprintf("stroke:%v; stroke-width:%v; fill: none;", "red", strokeWidth))
}

func getPointsForPolarGrid(centerX, centerY int, radius float64) (xs, ys []int) {
	for i := 0; i <= 24; i++ {
		grad := float64(i) * 15 / 180 * math.Pi

		xs = append(xs, centerX+int(radius*math.Cos(grad)))
		ys = append(ys, centerY+int(radius*math.Sin(grad)))
	}

	return xs, ys
}

func getRadiusAzimuthalGrid(zoom int) float64 {
	return 20 * float64(zoom+1) / 3
}

func getStrokeWidthAzimuthalGrid(radius float64) float64 {
	return float64(radius) / float64(100)
}

func contains(sample string, list []string) bool {
	for _, b := range list {
		if b == sample {
			return true
		}
	}
	return false
}

func renderCurve(canvas *svg.SVG, coords []geometry.Coord, style string) {
	xs, ys := coordToXsYs(coords)
	percentLength := 0.5

	if len(xs) <= 1 || len(xs) != len(ys) {
		panic("object must have more than one point")
	}

	if len(xs) >= 3 {
		for i := 0; i <= len(xs)-3; i++ {
			endArcX, endArcY := utils.GetPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], percentLength)
			beginArcX, beginArcY := utils.GetPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], 1-percentLength)

			if i == 0 {
				canvas.Line(xs[i], ys[i], beginArcX, beginArcY, style)
			} else {
				x, y := utils.GetPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], percentLength)
				canvas.Line(x, y, beginArcX, beginArcY, style)
			}

			canvas.Qbez(beginArcX, beginArcY, xs[i+1], ys[i+1], endArcX, endArcY, style)

			if i == len(xs)-3 {
				canvas.Line(endArcX, endArcY, xs[i+2], ys[i+2], style)
			} else {
				x, y := utils.GetPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], 1-percentLength)
				canvas.Line(endArcX, endArcY, x, y, style)
			}
		}
	} else {
		canvas.Line(xs[0], ys[0], xs[1], ys[1], style)
	}
}

func coordToXsYs(coords []geometry.Coord) ([]int, []int) {
	var xs, ys []int

	for _, coord := range coords {
		xs = append(xs, int(coord.X))
		ys = append(ys, int(coord.Y))
	}

	return xs, ys
}

func getPointBezier(bx1, by1, cx2, cy2, ex3, ey3, t float64) (float64, float64) {
	x := equationBezier(bx1, cx2, ex3, t)
	y := equationBezier(by1, cy2, ey3, t)

	return x, y
}

func equationBezier(px1, px2, px3, t float64) float64 {
	value := math.Pow(1-t, 2)*px1 + 2*t*(1-t)*px2 + math.Pow(t, 2)*px3

	return value
}

//step must be more 0 and less 1
//get points of bezier curve by three points
func bezierToPolyline(bx1, by1, cx2, cy2, ex3, ey3 int, step float64) ([]int, []int) {
	var xs []int
	var ys []int

	for i := .0; i <= 1; i += step {
		x, y := getPointBezier(float64(bx1), float64(by1), float64(cx2), float64(cy2), float64(ex3), float64(ey3), i)
		xs = append(xs, int(x))
		ys = append(ys, int(y))
	}

	return xs, ys
}

//convert points of polyline to points of curve bezier
func polylineToCurvePoints(xs, ys []int) ([]int, []int) {
	count := len(xs)
	percentLength := 0.5

	curveXs := []int{xs[0]}
	curveYs := []int{ys[0]}

	for i := 0; i < count-2; i++ {
		beginArcX, beginArcY := utils.GetPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], 1-percentLength)
		endArcX, endArcY := utils.GetPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], percentLength)

		bezierXs, bezierYs := bezierToPolyline(beginArcX, beginArcY, xs[i+1], ys[i+1], endArcX, endArcY, 0.1)

		curveXs = concat(curveXs, bezierXs)
		curveYs = concat(curveYs, bezierYs)
	}

	curveXs = append(curveXs, xs[count-1])
	curveYs = append(curveYs, ys[count-1])

	return curveXs, curveYs
}

func concat(array1, array2 []int) []int {
	newslice := make([]int, len(array1)+len(array2))
	copy(newslice, array1)
	copy(newslice[len(array1):], array2)
	return newslice
}

func renderCurveOrPolyline(canvas *svg.SVG, isRenderCurve bool, style string, coord []geometry.Coord) {
	if isRenderCurve {
		renderCurve(canvas, coord, style)
	} else {
		renderPolyline(canvas, coord, style)
	}

}
