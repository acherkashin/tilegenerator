package svg

import (
	"github.com/paulsmith/gogeos/geos"
	"github.com/ajstarks/svgo"
	"github.com/terrafactory/tilegenerator/mapobjects"
	"io"
)

func renderPoint(canvas *svg.SVG, geometry *geos.Geometry, tile *mapobjects.Tile) bool {
	coords, err := geometry.Coords()
	if (err != nil) {
		return false;
	}
	if (len(coords) < 1) {
		return false
	}
	x, y := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
	canvas.Circle(x, y, 10, "fill: black;")
	return true
}

func renderGeometry(canvas *svg.SVG, geometry *geos.Geometry, tile *mapobjects.Tile) bool {
	geometryType, err := geometry.Type()
	if (err != nil) {
		return false
	}
	switch  geometryType{
	case geos.POINT:
		renderPoint(canvas, geometry, tile)
	default:
		return false
	}
	return true;
}

func RenderTile(tile *mapobjects.Tile, geometries *[]geos.Geometry, writer io.Writer) {
	canvas := svg.New(writer)
	canvas.Start(mapobjects.TILE_SIZE, mapobjects.TILE_SIZE)
	for _, geo := range *geometries {
		renderGeometry(canvas, &geo, tile)
	}
	canvas.End()
}