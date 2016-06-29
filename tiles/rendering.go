package tiles

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings/styling"
	"github.com/TerraFactory/tilegenerator/utils"
	"github.com/TerraFactory/wktparser/geometry"
)

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
		object.TypeID = 181

		if object.TypeID == 47 || (object.TypeID >= 184 && object.TypeID <= 193) {
			RenderPatrollingArea(canvas, &object, tile)
		} else if object.TypeID == 74 || (object.TypeID >= 174 && object.TypeID <= 183) {
			RenderRouteAviationFlight(canvas, &object, tile)
		} else if object.TypeID == 408 {
			RenderPlannedAttackMainDirection(canvas, &object, tile)
		} else if object.TypeID == 407 {
			RenderAttackMainDirection(canvas, &object, tile)
		} else if object.TypeID == 366 {
			RenderCompletedProvideAction(canvas, &object, tile)
		} else if object.TypeID == 432 {
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

func getPointsForPolarGrid(centerX, centerY int, radius float64) (xs, ys []int) {
	for i := 0; i <= 24; i++ {
		grad := float64(i) * 15 / 180 * math.Pi

		xs = append(xs, centerX+int(radius*math.Cos(grad)))
		ys = append(ys, centerY+int(radius*math.Sin(grad)))
	}

	return xs, ys
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
	arrowSize := int(distance / 15)

	pointForArrowX := x1
	pointForArrowY := y1 + int(2*radiusY)
	xs, ys = getArrowPoints(pointForArrowX, pointForArrowY, arrowSize)

	return xs, ys
}

func getLeftArrowPoints(x1, y1, distance int) (xs []int, ys []int) {
	radiusX := distance / 4
	radiusY := radiusX / 2
	arrowSize := int(distance / 15)

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
		drawSegmentation(canvas, int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), style)
	}

	return nil
}

func drawSegmentation(canvas *svg.SVG, beginX, beginY, endX, endY int, style string) {
	length := 8
	distance := distanceBeetweenPoints(beginX, beginY, endX, endY)
	percentSizeSegment := float64(length) / distance
	count := int(distance) / length

	currentPointPercent := percentSizeSegment
	for i := 0; i < count-1; i += 2 {
		x1, y1 := equationLineByTwoPoints(beginX, beginY, endX, endY, currentPointPercent)
		currentPointPercent += percentSizeSegment
		x2, y2 := equationLineByTwoPoints(beginX, beginY, endX, endY, currentPointPercent)
		currentPointPercent += percentSizeSegment

		resultX, resultY := RotatePoint(x1, y1, x2, y2, -90)

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

	canvas.Group()

	renderPathRouteAviationFlight(coords, canvas, object, tile)
	renderArrowRouteAviationFlight(coords, canvas, object, tile)

	canvas.Gend()
	return nil
}

func renderPathRouteAviationFlight(coords []geometry.Coord, canvas *svg.SVG, object *entities.MapObject, tile *Tile) {
	fullLength := getLengthPolyline(coords, tile) / 2
	alreadyDrawn := false
	weight := 1

	var style string

	if object.ColorOuter != "" {
		style = fmt.Sprintf("stroke: %v; stroke-width: %v; fill: none; stroke-dasharray: 10;", object.ColorOuter, weight)
	} else {
		style = fmt.Sprintf("stroke:black; stroke-width: %v; fill: none; stroke-dasharray: 10;", weight)
	}

	var i int
	for i = 0; i < len(coords)-1; i++ {
		canvas.Line(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), style)

		lineLength := distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))
		fullLength -= lineLength

		if object.TypeID != 74 {
			if fullLength <= 0 && !alreadyDrawn {
				percentPosition := (-1.0) * fullLength / lineLength

				renderImageOnRouteAviation(canvas, object, int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y), tile.Z, percentPosition)
				alreadyDrawn = true
			}
		}
	}

}

func renderArrowRouteAviationFlight(coords []geometry.Coord, canvas *svg.SVG, object *entities.MapObject, tile *Tile) {
	lengthArrow := 5.0
	weight := 1
	i := len(coords) - 2
	lineLength := distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))

	var styleArrow string
	if object.ColorOuter != "" {
		styleArrow = fmt.Sprintf("stroke:%v; stroke-width: %v; fill: none;", object.ColorOuter, weight)
	} else {
		styleArrow = fmt.Sprintf("stroke:black; stroke-width: %v; fill: none;", weight)
	}

	if lineLength > lengthArrow {
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
	canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
	if object.ColorOuter == "" {
		canvas.CSS(`line, path, polyline {
			fill: none;
			stroke: black;
			}`)
	} else {
		canvas.CSS(fmt.Sprintf(`line, path, polyline {
			fill: none;
			stroke: %v;
			}`, object.ColorOuter))
	}

	fullLength := getLengthPolyline(coords, tile) / 2
	alreadyDrawn := false

	for i := 0; i < len(coords)-1; i++ {

		area := newPatrollingArea(tile, coords, i)
		lineLength := distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))
		fullLength -= lineLength

		transformation := fmt.Sprintf("rotate(%v,%v,%v)", area.rotateAngel, area.centerX, area.centerY)
		canvas.Gtransform(transformation)
		if object.TypeID != 47 {
			if fullLength <= 0 && !alreadyDrawn {
				percentPosition := (-1.0) * fullLength / lineLength

				x := float64(area.leftLinePointX) + (float64(area.rightLinePointX-area.leftLinePointX) * percentPosition)
				y := float64(area.leftLinePointY) + (float64(area.rightLinePointY-area.leftLinePointY) * percentPosition)

				renderImageOnPatrollingArea(canvas, object, area, x, y, tile.Z)
				alreadyDrawn = true
			}
		}

		canvas.Line(area.rightLinePointX, area.rightLinePointY, area.leftLinePointX, area.leftLinePointY)
		if i == 0 {
			canvas.Polyline(area.rightArrowXs, area.rightArrowYs)
			canvas.Arc(area.rightLinePointX,
				area.rightLinePointY,
				area.radiusX,
				area.radiusY,
				0, false, true,
				area.rightLinePointX,
				area.rightLinePointY+int(2*area.radiusY))
		}

		if i == len(coords)-2 {
			canvas.Polyline(area.leftArrowXs, area.leftArrowYs)
			canvas.Arc(area.leftLinePointX,
				area.leftLinePointY,
				area.radiusX,
				area.radiusY,
				0, false, true,
				area.leftLinePointX,
				area.leftLinePointY-int(2*area.radiusY))
		}
		canvas.Gend()
	}
	canvas.Gend()
	return nil
}

func renderImageOnPatrollingArea(canvas *svg.SVG, object *entities.MapObject, area *patrollingArea, x, y float64, zoom int) {
	bytesImg, err := utils.GetImgFromFile(filepath.FromSlash(fmt.Sprintf("images/%v.png", object.TypeID)))

	if err == nil {
		imgBase64Str := base64.StdEncoding.EncodeToString(bytesImg)

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

func renderImageOnRouteAviation(canvas *svg.SVG, object *entities.MapObject, x1, y1, x2, y2, zoom int, percentPosition float64) {
	bytesImg, err := utils.GetImgFromFile(filepath.FromSlash(fmt.Sprintf("images/%v.png", object.TypeID)))

	if err == nil {
		imgBase64Str := base64.StdEncoding.EncodeToString(bytesImg)

		img2html := "data:image/png;base64," + imgBase64Str

		imageWidth := 5 + 5*zoom
		imageHeight := 7 + 6*zoom
		angel := getAngel(x1, y1, x2, y2)

		centerX := x1 + (int)((float64)(x2-x1)*percentPosition)
		centerY := y1 + (int)((float64)(y2-y1)*percentPosition)

		canvas.Image(centerX-(int)(imageWidth/2.0),
			centerY-(int)(imageHeight/2.0),
			(int)(imageWidth),
			(int)(imageHeight),
			img2html,
			fmt.Sprintf("transform=\"rotate(%v,%v,%v)\"", angel-90, centerX, centerY))
	}
}

func getLengthPolyline(coords []geometry.Coord, tile *Tile) float64 {
	sum := 0.0

	for i := 0; i < len(coords)-1; i++ {
		sum += distanceBeetweenPoints(int(coords[i].X), int(coords[i].Y), int(coords[i+1].X), int(coords[i+1].Y))
	}

	return sum
}

func getArrowRouteAviationFlight(BeginX, BeginY, EndX, EndY, zoom int) ([]int, []int) {
	var angel int
	distance := distanceBeetweenPoints(BeginX, BeginY, EndX, EndY)
	percentSize := 1 - 5.0/distance
	angel = 120
	centerX, centerY := equationLineByTwoPoints(BeginX, BeginY, EndX, EndY, percentSize)
	rotatedPointX, rotatedPointY := EndX, EndY

	p1X, p1Y := RotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, angel)
	p2X, p2Y := RotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, -angel)

	xs := []int{p1X, rotatedPointX, p2X}
	ys := []int{p1Y, rotatedPointY, p2Y}
	return xs, ys
}

func equationLineByTwoPoints(BeginX, BeginY, EndX, EndY int, percentSize float64) (pointX, pointY int) {
	pointX = BeginX + (int)((float64)(EndX-BeginX)*percentSize)
	pointY = BeginY + (int)((float64)(EndY-BeginY)*percentSize)

	return pointX, pointY
}

func RotatePoint(centerX, centerY, pointX, pointY, angel int) (x, y int) {
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

	radius := 20 * float64(tile.Z+1) / 3
	centerX, centerY := int(point.Coordinates.X), int(point.Coordinates.Y)

	rotation := fmt.Sprintf("rotate(%v,%v,%v)", object.Azimut, centerX, centerY)
	canvas.Gtransform(rotation)
	strokeWidth := float64(radius) / float64(100)
	xs, ys := getBeamDiagramPoints(centerX, centerY, int(object.BeamWidth), object.Sidelobes, radius)
	templateStyle := "stroke:%v; stroke-width:%v; fill: none;"
	if object.ColorOuter != "" {
		canvas.Polygon(xs, ys, fmt.Sprintf(templateStyle, object.ColorOuter, strokeWidth))
	} else {
		canvas.Polygon(xs, ys, fmt.Sprintf(templateStyle, "red", strokeWidth))
	}
	canvas.Gend()
	return nil
}

//RenderAzimuthalGrid ...
func RenderAzimuthalGrid(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	point, err := object.Geometry.AsPoint()
	if err != nil {
		return err
	}

	radius := 20 * float64(tile.Z+1) / 3

	centerX, centerY := int(point.Coordinates.X), int(point.Coordinates.Y)
	strokeWidth := float64(radius) / float64(100)
	templateStyle := "stroke:%v; stroke-width:%v; fill: none;"
	rotation := fmt.Sprintf("rotate(%v,%v,%v)", object.Azimut, centerX, centerY)
	canvas.Gtransform(rotation)

	polarGridXs, polarGridYs := getPointsForPolarGrid(centerX, centerY, radius)

	var styleAzimuthalGrid string
	if object.ColorInner == "" {
		styleAzimuthalGrid = fmt.Sprintf(templateStyle, "gray", strokeWidth)
	} else {
		styleAzimuthalGrid = fmt.Sprintf(templateStyle, object.ColorInner, strokeWidth)
	}

	for i := 0; i < len(polarGridXs); i++ {
		canvas.Line(centerX, centerY, polarGridXs[i], polarGridYs[i], styleAzimuthalGrid)
	}

	indexZeroAzimut := 18
	canvas.Line(centerX, centerY, polarGridXs[indexZeroAzimut], polarGridYs[indexZeroAzimut], fmt.Sprintf(templateStyle, "red", strokeWidth))

	canvas.Circle(centerX, centerY, int(radius*0.67), fmt.Sprintf(templateStyle, "yellow", strokeWidth))
	canvas.Circle(centerX, centerY, int(radius), fmt.Sprintf(templateStyle, "green", strokeWidth))
	canvas.Gend()

	return nil
}
