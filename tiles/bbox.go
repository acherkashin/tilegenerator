package tiles

import (
	"math"
)

// BoundingBox contains the most north/south/east/west coordinates of tile.
type BoundingBox struct {
	North float64
	East  float64
	South float64
	West  float64
}

// Contains takes point latitude and longitude and returns true if this point is inside of this BoundingBox.
func (bbox *BoundingBox) Contains(lat, lon float64) bool {
	return (bbox.North >= lat && bbox.South <= lat) && (bbox.West <= lon && bbox.East >= lon)
}

// AddMargin calculates difference between North and South, devides it by 2 to get margin value, and then add this margin to BoundingBox.
// Example:
//  bbox:= BoundingBox{North: 50, West: 10, South: 10, East: 50}
//
//	bbox.AddMargin()
//
// now "bbox" equals to BoundingBox{North: 70, West: -10, South: -10, East: 70}
func (bbox *BoundingBox) AddMargin() {
	margin := math.Abs(bbox.North-bbox.South) / 2

	bbox.North = bbox.North + margin
	bbox.West = bbox.West - margin
	bbox.South = bbox.South - margin
	bbox.East = bbox.East + margin
}
