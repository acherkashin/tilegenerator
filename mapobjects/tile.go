package mapobjects

import (
	"math"
)
/* Size of each tile in pixels */
const TILE_SIZE = 256.0;

// region BoundingBox
/*
 Contains the most north/south/east/west coordinates of tile.
 */
type BoundingBox struct {
	North float64
	East  float64
	South float64
	West  float64
}

/*
 Takes point latitude and longitude and returns true if this point is inside of this BoundingBox.
 */
func (bbox BoundingBox) Contains(lat, lon float64) bool {
	return (bbox.North >= lat && bbox.South <= lat) && (bbox.West <= lon && bbox.East >= lon)
}
// endregion

// region Tile
/*
 Contains tile properties

 Z,X,Y - tile coordinates according to OSM specs(see http://wiki.openstreetmap.org/wiki/Slippy_map_tilenames)

 Bounding box - geographical coordinates of each side of tile
 */
type Tile struct {
	Z, X, Y     int
	Lat         float64
	Lon         float64
	BoundingBox BoundingBox
}

/*
 Returns longitude of the tile top side
 */
func Tile2lon(x int, z int) float64 {
	return float64(x) / math.Pow(2.0, float64(z)) * 360.0 - 180.0;
}

/*
 Returns latitude of the tile left side
 */
func Tile2lat(y int, z int) float64 {
	n := math.Pi - (2.0 * math.Pi * float64(y)) / math.Pow(2.0, float64(z));
	return math.Atan(math.Sinh(float64(n))) * 180 / math.Pi;
}

/*
 Tile factory function
 */
func NewTile(x int, y int, z int) Tile {
	return Tile{
		X: x, Y: y, Z:z,
		BoundingBox: BoundingBox{
			North : Tile2lat(y, z),
			South : Tile2lat(y + 1, z),
			West : Tile2lon(x, z),
			East : Tile2lon(x + 1, z),
		}}
}

func (tile *Tile) Lon2TileX(zoom int, lon_deg float64) int {
	x := (lon_deg + 180.0) / 360.0 * (math.Exp2(float64(zoom)))
	return int(math.Floor(TILE_SIZE * (x - float64(tile.X))))
}

func (tile *Tile) Lat2TileY(zoom int, lat_deg float64) int {
	y := (1.0 - math.Log(math.Tan(lat_deg * math.Pi / 180.0) + 1.0 / math.Cos(lat_deg * math.Pi / 180.0)) / math.Pi) / 2.0 * (math.Exp2(float64(zoom)))
	return int(math.Floor(TILE_SIZE * (y - float64(tile.Y))))
}


/*
 Takes point latitude and longitude and returns pixel coordinates of point on some tile.

 May return negative values as well as values outside of tile
 */
func (tile *Tile) Degrees2Pixels(lat, lon float64) (x int, y int) {
	return tile.Lon2TileX(tile.Z, lon), tile.Lat2TileY(tile.Z, lat)
}

/*
 Takes point latitude and longitude and returns true if this point is present on this tile.
 */
func (tile Tile) Contains(lat, lon float64) bool {
	return tile.BoundingBox.Contains(lat, lon)
}

// endregion