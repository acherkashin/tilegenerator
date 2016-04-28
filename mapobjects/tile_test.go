package mapobjects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTile(t *testing.T) {
	tile := NewTile(0, 0, 0)
	assert.Equal(t, tile.BoundingBox.North, 85.05112877980659)
	assert.Equal(t, tile.BoundingBox.South, -85.05112877980659)
	assert.Equal(t, tile.BoundingBox.East, 180.0)
	assert.Equal(t, tile.BoundingBox.West, -180.0)

	tile = NewTile(2475, 1280, 12)
	assert.Equal(t, tile.BoundingBox.North, 55.77657301866768)
	assert.Equal(t, tile.BoundingBox.South, 55.727110085045986)
	assert.Equal(t, tile.BoundingBox.East, 37.6171875)
	assert.Equal(t, tile.BoundingBox.West, 37.529296875)
}

func TestTile_Contains(t *testing.T) {
	tile := NewTile(0, 0, 0)
	assert.True(t, tile.Contains(0, 0), "this point is present on this tile")

	tile = NewTile(0, 1, 1)
	assert.False(t, tile.Contains(1, 1), "this point is not present on this tile")
}

func TestTile_Degrees2Pixels(t *testing.T) {
	tile := NewTile(0, 0, 0)
	x, y := tile.Degrees2Pixels(0, 0)
	assert.Equal(t, TILE_SIZE/2, x, "point with (0,0) coords should be exactly in the center of whole world tile")
	assert.Equal(t, TILE_SIZE/2, y, "point with (0,0) coords should be exactly in the center of whole world tile")

	tile = NewTile(1, 0, 1)
	x, y = tile.Degrees2Pixels(0, 0)
	assert.Equal(t, 0, x, "point with (0,0) coords should be exactly in the center of whole world tile")
	assert.Equal(t, TILE_SIZE, y, "point with (0,0) coords should be exactly in the center of whole world tile")
}
