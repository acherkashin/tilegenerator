package utils

import "math"

func RotatePoint(centerX, centerY, pointX, pointY, angle int) (x, y int) {
	radianAngle := (float64)(angle) * math.Pi / 180

	x = centerX + (int)((float64)(pointX-centerX)*math.Cos(radianAngle)-(float64)(pointY-centerY)*math.Sin(radianAngle))
	y = centerY + (int)((float64)(pointX-centerX)*math.Sin(radianAngle)+(float64)(pointY-centerY)*math.Cos(radianAngle))
	return x, y
}

//GetPointOnLine return point, which lies on the line is defined by two points
func GetPointOnLine(BeginX, BeginY, EndX, EndY int, percentPosition float64) (pointX, pointY int) {
	pointX = BeginX + (int)((float64)(EndX-BeginX)*percentPosition)
	pointY = BeginY + (int)((float64)(EndY-BeginY)*percentPosition)

	return pointX, pointY
}

func GetPointOnLineFloat(BeginX, BeginY, EndX, EndY float64, percentPosition float64) (pointX, pointY float64) {
	pointX = BeginX + (EndX-BeginX)*percentPosition
	pointY = BeginY + (EndY-BeginY)*percentPosition

	return pointX, pointY
}

//LengthPolyline calculate lenght of polyline
func LengthPolyline(xs, ys []float64) float64 {
	sum := 0.0

	for i := 0; i < len(xs)-1; i++ {
		sum += DistanceBetweenPoints(xs[i], ys[i], xs[i+1], ys[i+1])
	}

	return sum
}

//DistanceBetweenPoints returns distance between points (x1, y1) and (x2, y2)
func DistanceBetweenPoints(x1, y1, x2, y2 float64) float64 {
	a := x1 - x2
	b := y1 - y2

	return float64(math.Sqrt(float64(a*a + b*b)))
}

func LineCenter(x1, y1, x2, y2 float64) (x, y float64) {
	x = (x2 + x1) / 2
	y = (y2 + y1) / 2

	return x, y
}

// t - percent point's position on curve bezier
func PointBezier(bx1, by1, cx2, cy2, ex3, ey3, t float64) (float64, float64) {
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
func BezierToPolyline(bx1, by1, cx2, cy2, ex3, ey3 int, step float64) ([]int, []int) {
	var xs []int
	var ys []int

	for i := .0; i <= 1; i += step {
		x, y := PointBezier(float64(bx1), float64(by1), float64(cx2), float64(cy2), float64(ex3), float64(ey3), i)
		xs = append(xs, int(x))
		ys = append(ys, int(y))
	}

	return xs, ys
}

//PolylineToCurvePoints convert points of polyline to points of curve bezier
func PolylineToCurvePoints(xs, ys []int) ([]int, []int) {
	count := len(xs)
	percentLength := 0.5

	curveXs := []int{xs[0]}
	curveYs := []int{ys[0]}

	for i := 0; i < count-2; i++ {
		beginArcX, beginArcY := GetPointOnLine(xs[i], ys[i], xs[i+1], ys[i+1], 1-percentLength)
		endArcX, endArcY := GetPointOnLine(xs[i+1], ys[i+1], xs[i+2], ys[i+2], percentLength)

		bezierXs, bezierYs := BezierToPolyline(beginArcX, beginArcY, xs[i+1], ys[i+1], endArcX, endArcY, 0.1)

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

//IntArrayToFloat64 - convert int array to float64
func IntArrayToFloat64(xs []int) []float64 {
	var floatXs = []float64{}
	for _, item := range xs {
		floatXs = append(floatXs, float64(item))
	}

	return floatXs
}
