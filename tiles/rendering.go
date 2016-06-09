package tiles

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings/styling"
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

			}
		}
	}

	canvas.End()
}

type chartPoint struct {
	x, y, z, value float64
}

type beamDiagram struct {
	radius          float64
	sidelobes       float64
	angelRotation   float64
	sliderBeamWidth int
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

func newPatrollingArea(tile *Tile, coords []geometry.Coord, firstIndex int) *patrollingArea {
	x1, y1 := tile.Degrees2Pixels(coords[firstIndex].Y, coords[firstIndex].X)
	x2, y2 := tile.Degrees2Pixels(coords[firstIndex+1].Y, coords[firstIndex+1].X)

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
func RenderPatrollingArea(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}
	coords := line.Coordinates

	for i := 0; i < len(coords)-1; i++ {
		area := newPatrollingArea(tile, coords, i)
		transformation := fmt.Sprintf("rotate(%v,%v,%v)", area.rotateAngel, area.centerX, area.centerY)
		canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
		// canvas.CSS(prefixSelectors(object.CSS, object.ID))
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
		// canvas.Gend()
	}
	return nil
}

/*
RenderBeamDiagram - drawing satellite directional diagram
*/
func RenderBeamDiagram(canvas *svg.SVG, object *entities.MapObject, tile *Tile, beamDiagram *beamDiagram) error {
	point, err := object.Geometry.AsPoint()
	if err != nil {
		return err
	}

	// needShow := getValueAttribute(object.Attrs, "NEED_SHOW_DIRECTIONAL_DIAGRAM")
	// isAntenna := getValueAttribute(object.Attrs, "IS_SHORTWAVE_ANTENNA")

	// if needShow == "TRUE" && isAntenna == "TRUE" {
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
	// }
	return nil
}

func GetArrowPoints(BeginX, BeginY, EndX, EndY, zoom int) ([]int, []int) {
	var angel int
	var centerX, centerY, rotatedPointX, rotatedPointY int

	percentSize := 0.95 + 0.0019999*(float64)(zoom)
	angel = 120
	centerX = BeginX + (int)((float64)(EndX-BeginX)*percentSize)
	centerY = BeginY + (int)((float64)(EndY-BeginY)*percentSize)
	rotatedPointX = EndX
	rotatedPointY = EndY

	p1X, p1Y := RotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, angel)
	p2X, p2Y := RotatePoint(centerX, centerY, rotatedPointX, rotatedPointY, -angel)

	xs := []int{rotatedPointX, p1X, centerX, p2X}
	ys := []int{rotatedPointY, p1Y, centerY, p2Y}
	return xs, ys
}

func RotatePoint(centerX, centerY, pointX, pointY, angel int) (x, y int) {
	x = centerX + (int)((float64)(pointX-centerX)*math.Cos((float64)(angel)*math.Pi/180)) - (int)((float64)(pointY-centerY)*math.Sin((float64)(angel)/180*math.Pi))
	y = centerY + (int)((float64)(pointX-centerX)*math.Sin((float64)(angel)*math.Pi/180)) + (int)((float64)(pointY-centerY)*math.Cos((float64)(angel)/180*math.Pi))
	return x, y
}

// RenderRouteAviationFlight renders an aviation route on a tile
func RenderRouteAviationFlight(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
	line, err := object.Geometry.AsLineString()
	if err != nil {
		return err
	}
	coords := line.Coordinates
	weight := 1
	style := fmt.Sprintf("stroke:black; stroke-width: %v; fill: none; stroke-dasharray: 10;", weight)
	styleArrow := fmt.Sprintf("stroke:black; stroke-width: %v; fill: black;", weight)
	canvas.Group()

	var x1, y1, x2, y2 int

	for i := 0; i < len(coords)-1; i++ {
		x1, y1 = tile.Degrees2Pixels(coords[i].Y, coords[i].X)
		x2, y2 = tile.Degrees2Pixels(coords[i+1].Y, coords[i+1].X)
		canvas.Line(x1, y1, x2, y2, style)
		xs, ys := GetArrowPoints(x1, y1, x2, y2, tile.Z)
		canvas.Polygon(xs, ys, styleArrow)
	}

	canvas.Gend()
	return nil
}

// RenderSatelliteVisibility renders a atellites visibility chart
func RenderSatelliteVisibility(canvas *svg.SVG, object *entities.MapObject, radiomodules []*entities.MapObject, tile *Tile) error {
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

// func RenderSatellite(canvas *svg.SVG, object *entities.MapObject, tile *Tile) error {
// 	point, err := object.Geometry.AsPoint()
// 	if err != nil {
// 		return err
// 	}
// 	coord := point.Coordinates
// 	centerX, centerY := tile.Degrees2Pixels(coord.Y, coord.X)
// 	radius := 10

// 	var typeSatellite = object.Geometry.GetType()
// 	fill := determineFill(typeSatellite)
// 	stroke := determineStroke(object)
// 	text := determineText(typeSatellite)

// 	canvas.Circle(centerX, centerY, radius, fmt.Sprintf("%v; %v", stroke, fill))
// 	canvas.Line(centerX, centerY+radius, centerX+radius+radius/2, centerY+radius+radius/3, stroke)
// 	canvas.Line(centerX, centerY-radius, centerX+radius+radius/2, centerY-radius-radius/3, stroke)
// 	canvas.Text(centerX-radius/2, centerY+radius/5, text, "font-size:10px;")

// 	return nil
// }

// func determineStroke(object *entities.MapObject) string {
// 	valueAttributeAlly := getValueAttribute(object.Attrs, "ALLY_ENEMY")
// 	isAlly := valueAttributeAlly == ""

// 	if isAlly {
// 		return "stroke: red"
// 	}

// 	return "stroke: blue"
// }
// func getValueAttribute(attrs []geo.BaseAttribute, code string) string {
// 	for _, attr := range attrs {
// 		if attr.Code == code {
// 			return attr.Value
// 		}
// 	}
// 	return ""
// }

func determineFill(typeSatellite int) string {
	var fill string
	switch typeSatellite {
	case SatellitePhotoReconnaissance:
		fill = "fill: rgb(174,198,219)"
	case SatelliteOptoelectronicReconnaissance:
		fill = "fill: rgb(255, 230, 153)"
	case SatelliteRadiolocatingReconnaissance:
		fill = "fill: rgb(255,80,80)"
	case SatelliteCommunicationL:
		fill = "fill: rgb(255,255,0)"
	case SatelliteCommunicationS:
		fill = "fill: rgb(102,255,255)"
	case SatelliteCommunicationC:
		fill = "fill: rgb(51,102,255)"
	case SatelliteCommunicationX:
		fill = "fill: rgb(102,255,102)"
	case SatelliteCommunicationKu:
		fill = "fill: rgb(255,0,255)"
	case SatelliteCommunicationKa:
		fill = "fill: rgb(217,217,217)"
	default:
		fill = "fill: none"
	}
	return fill
}

func determineText(typeSatellite int) string {
	var text string
	switch typeSatellite {
	case SatellitePhotoReconnaissance,
		SatelliteOptoelectronicReconnaissance,
		SatelliteRadiolocatingReconnaissance:
		text = "Р"
	case SatelliteNavigation:
		text = "Н"
	case SatelliteCommunicationC,
		SatelliteCommunicationK,
		SatelliteCommunicationKa,
		SatelliteCommunicationKu,
		SatelliteCommunicationL,
		SatelliteCommunicationS,
		SatelliteCommunicationX:
		text = "C"
	case SatelliteMeteorological:
		text = "М"
	case SatelliteExperimental:
		text = "Э"
	case SatelliteRepeater:
		text = "РТ"
	case SatelliteRemoteSoundingEarth:
		text = "ДЗ"
	case SatelliteAttackWarning:
		text = "ПР"
	default:
		text = ""
	}

	return text
}

const (
	/// Спутник
	Satellite = 149
	/// Спутник связи (L-диапазон)
	SatelliteCommunicationL = 150
	/// Спутник связи (S-диапазон)
	SatelliteCommunicationS = 151
	/// Спутник связи (C-диапазон)
	SatelliteCommunicationC = 152
	/// Спутник связи (X-диапазон)
	SatelliteCommunicationX = 153
	/// Спутник связи (Ku-диапазон)
	SatelliteCommunicationKu = 154
	/// Спутник связи (Ka-диапазон)
	SatelliteCommunicationKa = 155
	/// Спутник связи (K-диапазон)
	SatelliteCommunicationK = 156
	/// Спутник фото-разведки
	SatellitePhotoReconnaissance = 161
	/// Оптикоэлектронный спутник разведки
	SatelliteOptoelectronicReconnaissance = 162
	/// Радиолокационный спутник разведки
	SatelliteRadiolocatingReconnaissance = 163
	/// Спутник навигации
	SatelliteNavigation = 160
	/// Метеорологический спутник
	SatelliteMeteorological = 159
	/// Эксперементальный спутник
	SatelliteExperimental = 158
	/// Спутник - ретранслятор
	SatelliteRepeater = 165
	/// Спутник дистанционного зондирования земли
	SatelliteRemoteSoundingEarth = 164
	/// Спутник системы предупреждения о ракетном нападении
	SatelliteAttackWarning = 157
)
