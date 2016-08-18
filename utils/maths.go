package utils

import "math"

func RotatePoint(centerX, centerY, pointX, pointY, angle int) (x, y int) {
	radianAngle := (float64)(angle) * math.Pi / 180

	x = centerX + (int)((float64)(pointX-centerX)*math.Cos(radianAngle)-(float64)(pointY-centerY)*math.Sin(radianAngle))
	y = centerY + (int)((float64)(pointX-centerX)*math.Sin(radianAngle)+(float64)(pointY-centerY)*math.Cos(radianAngle))
	return x, y
}

//equation of the line is defined by two points
func GetPointOnLine(BeginX, BeginY, EndX, EndY int, percentSize float64) (pointX, pointY int) {
	pointX = BeginX + (int)((float64)(EndX-BeginX)*percentSize)
	pointY = BeginY + (int)((float64)(EndY-BeginY)*percentSize)

	return pointX, pointY
}

func GetPointOnLineFloat(BeginX, BeginY, EndX, EndY float64, percentSize float64) (pointX, pointY float64) {
	pointX = BeginX + (EndX-BeginX)*percentSize
	pointY = BeginY + (EndY-BeginY)*percentSize

	return pointX, pointY
}
