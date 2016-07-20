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

				if object.IsAntenna && object.NeedShowDirectionalDiagram {
					RenderBeamDiagram(canvas, &object, tile)
				}

				if object.NeedShowAzimuthalGrid {
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

type patrollingArea struct {
	x1, y1, x2, y2                  int
	radiusX, radiusY                int
	leftLinePointX, rightLinePointX int
	leftLinePointY, rightLinePointY int
	centerX, centerY                int
	rotateAngel                     float64
	distance                        float64
	leftArrowXs, leftArrowYs        []int
	rightArrowXs, rightArrowYs      []int
}

type bigArrow struct {
	pointXs, pointYs, arrowXs, arrowYs []int
	centerX, centerY                   int

	rotateAngel float64
}

func newPatrollingArea(tile *Tile, coords []geometry.Coord, firstIndex int) *patrollingArea {
	x1, y1 := int(coords[firstIndex].X), int(coords[firstIndex].Y)
	x2, y2 := int(coords[firstIndex+1].X), int(coords[firstIndex+1].Y)

	radiusX := int(distanceBeetweenPoints(x1, y1, x2, y2) / 4)
	radiusY := int(radiusX / 2)
	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	rightLinePointX, rightLinePointY := centerX+int(distance/2), centerY
	leftLinePointX, leftLinePointY := centerX-int(distance/2), centerY
	leftArrowXs, leftArrowYs := getLeftArrowPoints(leftLinePointX, leftLinePointY, int(distance))
	rightArrowXs, rightArrowYs := getRightArrowPoints(rightLinePointX, rightLinePointY, int(distance))

	return &patrollingArea{
		x1: x1, y1: y1, x2: x2, y2: y2,
		radiusX:         radiusX,
		radiusY:         radiusY,
		centerX:         centerX,
		centerY:         centerY,
		distance:        distance,
		rightLinePointX: rightLinePointX,
		rightLinePointY: rightLinePointY,
		leftLinePointX:  leftLinePointX,
		leftLinePointY:  leftLinePointY,
		leftArrowXs:     leftArrowXs,
		leftArrowYs:     leftArrowYs,
		rightArrowXs:    rightArrowXs,
		rightArrowYs:    rightArrowYs,
		rotateAngel:     getAngel(x1, y1, x2, y2),
	}
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

func RenderPit(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}

	coords := line.Coordinates
	weight := 1
	if object.ColorOuter == "" {
		object.ColorOuter = "black"
	}

	style := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.ColorOuter, weight)

	for i := 0; i < len(coords)-1; i++ {
		canvas.Line(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), style)
		drawHatching(canvas, int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), style)
	}

	return nil
}

func drawHatching(canvas *svg.SVG, beginX, beginY, endX, endY int, style string) {
	length := 8
	distance := distanceBeetweenPoints(beginX, beginY, endX, endY)
	percentSizeSegment := float64(length) / distance
	count := int(distance) / length

	currentPointPercent := percentSizeSegment
	for i := 0; i < count-1; i += 2 {
		x1, y1 := getPointOnLine(beginX, beginY, endX, endY, currentPointPercent)
		currentPointPercent += percentSizeSegment
		x2, y2 := getPointOnLine(beginX, beginY, endX, endY, currentPointPercent)
		currentPointPercent += percentSizeSegment

		resultX, resultY := rotatePoint(x1, y1, x2, y2, -90)

		canvas.Line(x1, y1, resultX, resultY, style)
	}
}

func RenderAttackMainDirection(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	if object.ColorInner == "" {
		object.ColorInner = "red"
	}

	if object.ColorOuter == "" {
		object.ColorOuter = "red"
	}

	weight := 1

	styleLine := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.ColorOuter, weight)
	styleArrow := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: %v;", object.ColorOuter, weight, object.ColorInner)

	return renderBigArrow(canvas, object, tile, styleLine, styleArrow)
}

func RenderPlannedAttackMainDirection(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	if object.ColorInner == "" {
		object.ColorInner = "red"
	}

	if object.ColorOuter == "" {
		object.ColorOuter = "red"
	}

	weight := 1

	styleLine := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;stroke-dasharray: 10;", object.ColorOuter, weight)
	styleArrow := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: %v;", object.ColorOuter, weight, object.ColorInner)

	return renderBigArrow(canvas, object, tile, styleLine, styleArrow)
}

func RenderCompletedProvideAction(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	weight := 1
	if object.ColorOuter != "" {
		object.ColorOuter = "red"
	}

	style := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.ColorOuter, weight)

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

	style := fmt.Sprintf("stroke: %v; stroke-width: %v; fill: none; stroke-dasharray: 10;", object.ColorOuter, 1)
	// renderPolyline(canvas, coords, style)
	renderCurve(canvas, coords, style)

	if object.Code != "1000000004" {
		xs, ys := coordToXsYs(coords)
		curveXs, curveYs := polylineToCurvePoints(xs, ys)
		x, y, angel := getCenterPolylineAndAngel(curveXs, curveYs)
		renderImageOnRouteAviation(canvas, object, angel, x, y, tile.Z)
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

	centerX, centerY := getPointOnLine(xs[i-1], ys[i-1], xs[i], ys[i], percentPosition)

	angel = getAngel(xs[i-1], ys[i-1], xs[i], ys[i])

	return centerX, centerY, angel
}

func renderArrowRouteAviationFlight(coords []geometry.Coord, canvas *svg.SVG, object *entities.MapObject, tile *Tile) {
	lengthArrow := 5.0
	weight := 1
	i := len(coords) - 2
	lastLineLength := distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))

	styleArrow := fmt.Sprintf("stroke:%v; stroke-width: %v; fill: %v;", object.ColorInner, weight, object.ColorInner)

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

	setDefaultColor(object)

	canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")

	xs, ys := coordToXsYs(coords)
	halfLength := getLengthPolyline(xs, ys) / 2
	alreadyDrawn := false

	for i := 0; i < len(coords)-1; i++ {

		area := newPatrollingArea(tile, coords, i)
		lineLength := distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))
		halfLength -= lineLength

		transformation := fmt.Sprintf("rotate(%v,%v,%v)", area.rotateAngel, area.centerX, area.centerY)
		canvas.Gtransform(transformation)
		if object.Code != "1000000002" {
			if halfLength <= 0 && !alreadyDrawn {
				percentPosition := (-1.0) * halfLength / lineLength

				x := float64(area.leftLinePointX) + (float64(area.rightLinePointX-area.leftLinePointX) * percentPosition)
				y := float64(area.leftLinePointY) + (float64(area.rightLinePointY-area.leftLinePointY) * percentPosition)
				renderImageOnPatrollingArea(canvas, object, area, x, y, tile.Z)
				alreadyDrawn = true
			}
		}

		canvas.Line(area.rightLinePointX, area.rightLinePointY, area.leftLinePointX, area.leftLinePointY, fmt.Sprintf("stroke: %v; fill: none;", object.ColorOuter))
		if i == 0 {
			renderRightPartPatrollingArea(canvas, area, object)
		}

		if i == len(coords)-2 {
			renderLeftPartPatrollingArea(canvas, area, object)
		}
		canvas.Gend()
	}
	canvas.Gend()
	return nil
}

func renderRightPartPatrollingArea(canvas *svg.SVG, area *patrollingArea, object *entities.MapObject) {
	canvas.Arc(area.rightLinePointX,
		area.rightLinePointY,
		area.radiusX,
		area.radiusY,
		0, false, true,
		area.rightLinePointX,
		area.rightLinePointY+int(2*area.radiusY),
		fmt.Sprintf("stroke: %v; fill: none;", object.ColorOuter))
	if area.radiusX > 10 {
		canvas.Polyline(area.rightArrowXs, area.rightArrowYs, fmt.Sprintf("stroke: %v; fill: %v;", object.ColorInner, object.ColorInner))
	}
}

func renderLeftPartPatrollingArea(canvas *svg.SVG, area *patrollingArea, object *entities.MapObject) {
	canvas.Arc(area.leftLinePointX,
		area.leftLinePointY,
		area.radiusX,
		area.radiusY,
		0, false, true,
		area.leftLinePointX,
		area.leftLinePointY-int(2*area.radiusY),
		fmt.Sprintf("stroke: %v; fill: none", object.ColorOuter))
	if area.radiusX > 10 {
		canvas.Polyline(area.leftArrowXs, area.leftArrowYs, fmt.Sprintf("stroke: %v; fill: %v;", object.ColorInner, object.ColorInner))
	}
}

func setDefaultColor(object *entities.MapObject) {
	if object.ColorInner == "" {
		object.ColorInner = "black"
	}
	if object.ColorOuter == "" {
		object.ColorOuter = "black"
	}
}

func renderImageOnPatrollingArea(canvas *svg.SVG, object *entities.MapObject, area *patrollingArea, x, y float64, zoom int) {
	pathConfig := "./config.toml"
	settings, err := settings.GetSettings(&pathConfig)

	if err == nil {
		href := fmt.Sprintf("%v/api/maps/object/%v/png", settings.UrlAPI, object.ID)
		if result, err := utils.GetImgByURL(href); err == nil {
			imgBase64Str := base64.StdEncoding.EncodeToString(result)

			img2html := "data:image/png;base64," + imgBase64Str

			imageWidth := 5 + 5*zoom
			imageHeight := 7 + 6*zoom

			canvas.Image(int(x)-(int)(imageWidth/2.0),
				int(y)-(int)(imageHeight/2.0),
				int(imageWidth),
				int(imageHeight),
				img2html,
				fmt.Sprintf("transform=\"rotate(%v,%v,%v)\"", -90, x, y))
		}
	}
}

func renderImageOnRouteAviation(canvas *svg.SVG, object *entities.MapObject, angel float64, x, y, zoom int) {
	pathConfig := "./config.toml"
	settings, err := settings.GetSettings(&pathConfig)

	if err == nil {
		href := fmt.Sprintf("%v/api/maps/object/%v/png", settings.UrlAPI, object.ID)
		if result, err := utils.GetImgByURL(href); err == nil {
			imgBase64Str := base64.StdEncoding.EncodeToString(result)

			img2html := "data:image/png;base64," + imgBase64Str

			imageWidth := 5 + 5*zoom
			imageHeight := 7 + 6*zoom

			canvas.Image(x-(int)(imageWidth/2.0),
				y-(int)(imageHeight/2.0),
				(int)(imageWidth),
				(int)(imageHeight),
				img2html,
				fmt.Sprintf("transform=\"rotate(%v,%v,%v)\"", angel-90, x, y))
		}
	}
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
	centerX, centerY := getPointOnLine(BeginX, BeginY, EndX, EndY, percentSize)
	rotatedPointX, rotatedPointY := EndX, EndY

	p1X, p1Y := rotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, angel)
	p2X, p2Y := rotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, -angel)

	xs := []int{p1X, rotatedPointX, p2X}
	ys := []int{p1Y, rotatedPointY, p2Y}
	return xs, ys
}

//equation of the line is defined by two points
func getPointOnLine(BeginX, BeginY, EndX, EndY int, percentSize float64) (pointX, pointY int) {
	pointX = BeginX + (int)((float64)(EndX-BeginX)*percentSize)
	pointY = BeginY + (int)((float64)(EndY-BeginY)*percentSize)

	return pointX, pointY
}

func rotatePoint(centerX, centerY, pointX, pointY, angel int) (x, y int) {
	x = centerX + (int)((float64)(pointX-centerX)*math.Cos((float64)(angel)*math.Pi/180)) - (int)((float64)(pointY-centerY)*math.Sin((float64)(angel)/180*math.Pi))
	y = centerY + (int)((float64)(pointX-centerX)*math.Sin((float64)(angel)*math.Pi/180)) + (int)((float64)(pointY-centerY)*math.Cos((float64)(angel)/180*math.Pi))
	return x, y
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

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", object.Azimut, centerX, centerY))
	xs, ys := getBeamDiagramPoints(centerX, centerY, int(object.BeamWidth), object.Sidelobes, radius)

	if object.ColorOuter == "" {
		object.ColorOuter = "red"
	}

	canvas.Polygon(xs, ys, fmt.Sprintf("stroke:%v; stroke-width:%v; fill: none;", object.ColorOuter, strokeWidth))

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

	canvas.Gtransform(fmt.Sprintf("rotate(%v,%v,%v)", object.Azimut, centerX, centerY))
	renderGradationAzimuthalGrid(canvas, object, centerX, centerY, tile.Z)
	canvas.Circle(centerX, centerY, int(radius*0.67), fmt.Sprintf(templateStyle, "yellow", strokeWidth))
	canvas.Circle(centerX, centerY, int(radius), fmt.Sprintf(templateStyle, "green", strokeWidth))
	canvas.Gend()

	return nil
}

func renderGradationAzimuthalGrid(canvas *svg.SVG, object *entities.MapObject, centerX, centerY, zoom int) {
	if object.ColorInner == "" {
		object.ColorInner = "gray"
	}
	radius := getRadiusAzimuthalGrid(zoom)
	strokeWidth := getStrokeWidthAzimuthalGrid(radius)
	styleAzimuthalGrid := fmt.Sprintf("stroke:%v; stroke-width:%v; fill: none;", object.ColorInner, strokeWidth)

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
	count := len(xs)

	if count <= 1 {
		panic("object must have more than one point")
	}

	percentLength := 0.5
	if count >= 3 {
		for i := 0; i <= count-3; i++ {
			endArcX, endArcY := getPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], percentLength)
			beginArcX, beginArcY := getPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], 1-percentLength)

			if i == 0 {
				canvas.Line(xs[i], ys[i], beginArcX, beginArcY, style)
			} else {
				x, y := getPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], percentLength)
				canvas.Line(x, y, beginArcX, beginArcY, style)
			}

			canvas.Qbez(beginArcX, beginArcY, xs[i+1], ys[i+1], endArcX, endArcY, style)

			if i == count-3 {
				canvas.Line(endArcX, endArcY, xs[i+2], ys[i+2], style)
			} else {
				x, y := getPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], 1-percentLength)
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

func polylineToCurvePoints(xs, ys []int) ([]int, []int) {
	count := len(xs)
	percentLength := 0.5

	curveXs := []int{xs[0]}
	curveYs := []int{ys[0]}

	for i := 0; i < count-2; i++ {
		beginArcX, beginArcY := getPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], 1-percentLength)
		endArcX, endArcY := getPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], percentLength)

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
