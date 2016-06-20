package tiles

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var bboxTestCases = map[BoundingBox]BoundingBox{
	BoundingBox{North: 50, West: 10, South: 10, East: 50}:   BoundingBox{North: 70, West: -10, South: -10, East: 70},
	BoundingBox{North: -50, West: 20, South: 20, East: -50}: BoundingBox{North: -15, West: -15, South: -15, East: -15},
	BoundingBox{North: 0, West: 0, South: 0, East: 0}:       BoundingBox{North: 0, West: 0, South: 0, East: 0},
}

func TestAddMarginToBbox(t *testing.T) {
	for key, value := range bboxTestCases {
		key.AddMargin()
		assert.Equal(t, value, key)
	}
}
