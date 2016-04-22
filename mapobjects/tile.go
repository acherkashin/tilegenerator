package mapobjects

import (
	"math"
)

type Tile struct {
	Z   int
	X   int
	Y   int
	Lat float64
	Lon float64
}

type Conversion interface {
	deg2num(t *Tile) (x int, y int)
	num2deg(t *Tile) (lat float64, lon float64)
}

func (*Tile) Deg2num(t *Tile) (x int, y int) {
	x = int(math.Floor((t.Lon + 180.0) / 360.0 * (math.Exp2(float64(t.Z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(t.Lat*math.Pi/180.0)+1.0/math.Cos(t.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(t.Z)))))
	return
}

func (*Tile) Num2deg(t *Tile) (lat float64, lon float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	lon = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return lat, lon
}
