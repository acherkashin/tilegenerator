package tiles

import (
	"math"
)

// TileSize is a size of each tile in pixels
const TileSize = 256

// Tile contains tile properties
// Z,X,Y - tile coordinates according to OSM specs(see http://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)
// Bounding box - geographical coordinates of each side of tile
type Tile struct {
	Z, X, Y     int
	Lat         float64
	Lon         float64
	BoundingBox BoundingBox
}

// Tile2lon returns longitude of the tile top side
func Tile2lon(x int, z int) float64 {
	return float64(x)/math.Pow(2.0, float64(z))*360.0 - 180.0
}

// Tile2lat returns latitude of the tile left side
func Tile2lat(y int, z int) float64 {
	n := math.Pi - (2.0*math.Pi*float64(y))/math.Pow(2.0, float64(z))
	return math.Atan(math.Sinh(float64(n))) * 180 / math.Pi
}

// NewTile is a tile factory function
func NewTile(x int, y int, z int) *Tile {
	return &Tile{
		X: x, Y: y, Z: z,
		BoundingBox: BoundingBox{
			North: Tile2lat(y, z),
			South: Tile2lat(y+1, z),
			West:  Tile2lon(x, z),
			East:  Tile2lon(x+1, z),
		}}
}

// Lon2TileX converts longitude into a tile X coordinate
func (tile *Tile) Lon2TileX(zoom int, lonDeg float64) int {
	x := (lonDeg + 180.0) / 360.0 * (math.Exp2(float64(zoom)))
	return int(math.Floor(TileSize * (x - float64(tile.X))))
}

// Lat2TileY converts latitude into a tile Y coordinate
func (tile *Tile) Lat2TileY(zoom int, latDeg float64) int {
	y := (1.0 - math.Log(math.Tan(latDeg*math.Pi/180.0)+1.0/math.Cos(latDeg*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))
	return int(math.Floor(TileSize * (y - float64(tile.Y))))
}

// Degrees2Pixels takes point latitude and longitude and returns pixel coordinates of point on some tile.
// May return negative values as well as values outside of tile
func (tile *Tile) Degrees2Pixels(lat, lon float64) (x int, y int) {
	return tile.Lon2TileX(tile.Z, lon), tile.Lat2TileY(tile.Z, lat)
}

// Contains takes point latitude and longitude and returns true if this point is present on this tile.
func (tile *Tile) Contains(lat, lon float64) bool {
	return tile.BoundingBox.Contains(lat, lon)
}
