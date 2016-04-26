package mapobjects

import (
	"math"
)

type Tile struct {
	Z           int
	X           int
	Y           int
	Lat         float64
	Lon         float64
	BoundingBox BoundingBox
}

type BoundingBox struct {
	North float64
	East  float64
	South float64
	West  float64
}

func Tile2lon(x int, z int) float64 {
	return float64(x) / math.Pow(2.0, float64(z)) * 360.0 - 180.0;
}

func Tile2lat(y int, z int) float64 {
	n := math.Pi - (2.0 * math.Pi * float64(y)) / math.Pow(2.0, float64(z));
	return math.Atan(math.Sinh(float64(n))) * 180 / math.Pi;
}

type Conversion interface {
	deg2num(t *Tile) (x int, y int)
	num2deg(t *Tile) (lat float64, lon float64)
}

func (*Tile) Deg2num(t *Tile) (x int, y int) {
	x = int(math.Floor((t.Lon + 180.0) / 360.0 * (math.Exp2(float64(t.Z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(t.Lat * math.Pi / 180.0) + 1.0 / math.Cos(t.Lat * math.Pi / 180.0)) / math.Pi) / 2.0 * (math.Exp2(float64(t.Z)))))
	return
}

func (*Tile) Num2deg(t *Tile) (lat float64, lon float64) {
	lat = Tile2lat(t.Y, t.Z)
	lon = Tile2lon(t.X, t.Z)
	return lat, lon
}

func NewTile(x int, y int, z int) *Tile {
	tile := new(Tile)
	tile.X = x
	tile.Y = y
	tile.Z = z
	tile.BoundingBox = BoundingBox{}
	tile.BoundingBox.North = Tile2lat(y, z)
	tile.BoundingBox.South = Tile2lat(y + 1, z)
	tile.BoundingBox.West = Tile2lon(x, z)
	tile.BoundingBox.East = Tile2lon(x + 1, z)
	return tile
}
