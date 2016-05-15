package svg

import (
	"fmt"
	"math"
	"strconv"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/mapobjects"
	"github.com/TerraFactory/wktparser/geometry"
)

type chartPoint struct {
	x, y, z, value float64
}

type beamDiagram struct {
	radius          float64
	sidelobes       float64
	angelRotation   float64
	sliderBeamWidth int
}

type routeAviationFlight struct {
	x1, y1, x2, y2                  int
	leftLinePointX, rightLinePointX int
	centerX, centerY                int
	arrowSize                       int
	rotateAngel                     float64
	arrowXs, arrowYs                []int
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

func newPatrollingArea(tile *mapobjects.Tile, coords []geometry.Coord) *patrollingArea {
	x1, y1 := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
	x2, y2 := tile.Degrees2Pixels(coords[1].Y, coords[1].X)

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

func newCharPoint(x, y, z, value float64) *chartPoint {
	var cp chartPoint
	cp.x = x
	cp.y = y
	cp.z = z
	cp.value = value

	return &cp
}

func newRouteAviationFlight(tile *mapobjects.Tile, coords []geometry.Coord) *routeAviationFlight {
	x1, y1 := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
	x2, y2 := tile.Degrees2Pixels(coords[1].Y, coords[1].X)

	centerX, centerY := getLineCenter(x1, y1, x2, y2)
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	rightLinePointX := centerX + int(distance)/2
	leftLinePoint := centerX - int(distance)/2
	arrowSize := int(distance / 15)

	arrowXs := []int{}
	arrowYs := []int{}

	if coords[0].X < coords[1].X {
		arrowXs = append(arrowXs, rightLinePointX-arrowSize,
			rightLinePointX,
			rightLinePointX-arrowSize)
	} else {
		arrowXs = append(arrowXs, leftLinePoint+arrowSize,
			leftLinePoint,
			leftLinePoint+arrowSize)
	}

	arrowYs = append(arrowYs, centerY+arrowSize/2, centerY, centerY-arrowSize/2)

	return &routeAviationFlight{
		x1: x1, y1: y1,
		x2: x2, y2: y2,
		rightLinePointX: rightLinePointX,
		leftLinePointX:  leftLinePoint,
		rotateAngel:     getAngel(x1, y1, x2, y2),
		arrowSize:       arrowSize,
		centerX:         centerX,
		centerY:         centerY,
		arrowXs:         arrowXs,
		arrowYs:         arrowYs,
	}
}

func (beamDiagram *beamDiagram) getPoints(centerX, centerY int) (xs, ys []int) {
	tempPointList := getTempPointsForBeamPoints(beamDiagram.sliderBeamWidth, beamDiagram.sidelobes)
	maxValue := getMax(tempPointList)

	for _, point := range tempPointList {
		xs = append(xs, centerX+int(point.y*beamDiagram.radius/maxValue))
		ys = append(ys, centerY+int(point.x*beamDiagram.radius/maxValue))
	}

	xs[0] = centerX + int(beamDiagram.radius)
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

/*
RenderPatrollingArea (район барражирования)
*/
func RenderPatrollingArea(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}
	coords := line.Coordinates

	area := newPatrollingArea(tile, coords)
	transformation := fmt.Sprintf("rotate(%v,%v,%v)", area.rotateAngel, area.centerX, area.centerY)
	canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
	canvas.CSS(prefixSelectors(object.CSS, object.ID))
	canvas.Gtransform(transformation)
	canvas.Line(area.rightLinePointX, area.rightLinePointY, area.leftLinePointX, area.leftLinePointY)
	canvas.Polyline(area.rightArrowXs, area.rightArrowYs)
	canvas.Polyline(area.leftArrowXs, area.leftArrowYs)
	canvas.Arc(area.rightLinePointX,
		area.rightLinePointY,
		area.radiusX,
		area.radiusY,
		0, false, true,
		area.rightLinePointX,
		area.rightLinePointY+int(2*area.radiusY))
	canvas.Arc(area.leftLinePointX,
		area.leftLinePointY,
		area.radiusX,
		area.radiusY,
		0, false, true,
		area.leftLinePointX,
		area.leftLinePointY-int(2*area.radiusY))
	canvas.Gend()
	canvas.Gend()

	return nil
}

/*
RenderBeamDiagram - drawing satellite directional diagram
*/
func RenderBeamDiagram(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile, beamDiagram *beamDiagram) error {
	point, err := object.Geometry.AsPoint()
	if err != nil {
		return err
	}
	coord := point.Coordinates

	beamDiagram.radius *= float64(tile.Z+1) / 3
	centerX, centerY := tile.Degrees2Pixels(coord.Y, coord.X)

	strokeWidth := float64(beamDiagram.radius) / float64(100)
	xs, ys := beamDiagram.getPoints(centerX, centerY)
	templateStyle := "stroke:%v; stroke-width:%v; fill: none;"
	rotation := fmt.Sprintf("rotate(%v,%v,%v)", beamDiagram.angelRotation, centerX, centerY)

	canvas.Gtransform(rotation)
	polarGridXs, polarGridYs := getPointsForPolarGrid(centerX, centerY, beamDiagram.radius)
	for i := 0; i < len(polarGridXs); i++ {
		canvas.Line(centerX, centerY, polarGridXs[i], polarGridYs[i], fmt.Sprintf(templateStyle, "gray", strokeWidth))
	}
	canvas.Circle(centerX, centerY, int(beamDiagram.radius*0.67), fmt.Sprintf(templateStyle, "yellow", strokeWidth))
	canvas.Circle(centerX, centerY, int(beamDiagram.radius), fmt.Sprintf(templateStyle, "green", strokeWidth))
	canvas.Polygon(xs, ys, fmt.Sprintf(templateStyle, "red", strokeWidth))
	canvas.Gend()

	return nil
}

// RenderRouteAviationFlight renders an aviation route on a tile
func RenderRouteAviationFlight(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}
	coords := line.Coordinates

	route := newRouteAviationFlight(tile, coords)
	weight := 1
	style := fmt.Sprintf("stroke:black; stroke-width: %v; fill: none; stroke-dasharray: 10;", weight)
	styleArrow := fmt.Sprintf("stroke:black; stroke-width: %v; fill: none;", weight)
	transformation := fmt.Sprintf("rotate(%v,%v,%v)", route.rotateAngel, route.centerX, route.centerY)

	canvas.Group(fmt.Sprintf("id=\"id%v\"  transform=\"%v\"", strconv.Itoa(object.ID), transformation))
	canvas.CSS(prefixSelectors(object.CSS, object.ID))
	canvas.Line(route.rightLinePointX, route.centerY, route.leftLinePointX, route.centerY, style)
	canvas.Polyline(route.arrowXs, route.arrowYs, styleArrow)
	canvas.Gend()

	return nil
}

// RenderSatelliteVisibility renders a atellites visibility chart
func RenderSatelliteVisibility(canvas *svg.SVG, object *mapobjects.MapObject, radiomodules []*mapobjects.MapObject, tile *mapobjects.Tile) error {
	point, err := object.Geometry.AsPoint()
	if err != nil {
		return err
	}

	coord := point.Coordinates

	x1, y1 := tile.Degrees2Pixels(coord.Y, coord.X)
	x2, y2 := tile.Degrees2Pixels(float64(coord.Y+2), float64(coord.X+2))
	distance := distanceBeetweenPoints(x1, y1, x2, y2)
	canvas.Group("fill-opacity=\".3\"")

	canvas.ClipPath("id=\"clip-ellipse\"")
	canvas.Ellipse(x1, y1, int(distance), int(distance*0.7))
	canvas.ClipEnd()

	canvas.Ellipse(x1, y1, int(distance), int(distance*0.7), "stroke:black; fill: green; ")

	for _, radioModule := range radiomodules {
		rmPoint, err := radioModule.Geometry.AsPoint()
		if err == nil {
			rmCoords := rmPoint.Coordinates
			x, y := tile.Degrees2Pixels(rmCoords.Y, rmCoords.X)
			canvas.Circle(x, y, int(distance*0.2), "stroke:black; fill: blue; clip-path: url(#clip-ellipse)")
		}
	}
	canvas.Gend()

	return nil
}
