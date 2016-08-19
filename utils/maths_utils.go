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

//IntArrayToFloat64 - convert int array to float64
func IntArrayToFloat64(xs []int) []float64 {
	var floatXs = []float64{}

	for _, item := range xs {
		floatXs = append(floatXs, float64(item))
	}

	return floatXs
}

func LineCenter(x1, y1, x2, y2 float64) (x, y float64) {
	x = (x2 + x1) / 2
	y = (y2 + y1) / 2

	return x, y
}
