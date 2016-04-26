package mapobjects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeg2NumConversion(t *testing.T) {
	tile := Tile{Z: 12, X: 2475, Y: 1280, Lat: 55.77657301866769, Lon: 37.529296875}

	x, y := tile.Deg2num(&tile)

	assert.Equal(t, tile.X, x, "should be equal")
	assert.Equal(t, tile.Y, y, "should be equal")
}

func TestNum2degConversion(t *testing.T) {
	tile := Tile{Z: 12, X: 2475, Y: 1280, Lat: 55.77657301866768, Lon: 37.529296875}

	lat, lon := tile.Num2deg(&tile)

	assert.Equal(t, tile.Lat, lat, "should be equal")
	assert.Equal(t, tile.Lon, lon, "should be equal")
}

func TestNewTile(t *testing.T) {
	tile := NewTile(0, 0, 0)

	assert.Equal(t, tile.BoundingBox.North, 85.05112877980659)
	assert.Equal(t, tile.BoundingBox.South, -85.05112877980659)
	assert.Equal(t, tile.BoundingBox.East, 180.0)
	assert.Equal(t, tile.BoundingBox.West, -180.0)
}
