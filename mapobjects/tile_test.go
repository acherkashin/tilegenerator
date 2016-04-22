package mapobjects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeg2NumConversion(t *testing.T) {
	tile := Tile{Z: 12, X: 2475, Y: 1280, Lat: 55.77657301866769, Lon: 37.529296875}

	x, y := tile.Deg2num(&tile)

	assert.Equal(t, x, tile.X, "should be equal")
	assert.Equal(t, y, tile.Y, "should be equal")
}

func TestNum2degConversion(t *testing.T) {
	tile := Tile{Z: 12, X: 2475, Y: 1280, Lat: 55.77657301866769, Lon: 37.529296875}

	lat, lon := tile.Num2deg(&tile)

	assert.Equal(t, lat, tile.Lat, "should be equal")
	assert.Equal(t, lon, tile.Lon, "should be equal")
}
