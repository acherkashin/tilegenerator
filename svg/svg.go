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
	canvas.Circle(x, y, 5, "fill: black;")
	return true
}

func renderPolyline(canvas *svg.SVG, geometry *geos.Geometry, tile *mapobjects.Tile) bool {
	coords, err := geometry.Coords()
	if (err != nil) {
		return false;
	}
	xs := []int{}
	ys := []int{}
	for _, coord := range coords {
		x, y := tile.Degrees2Pixels(coord.Y, coord.X)
		xs = append(xs, x)
		ys = append(ys, y)
	}
	canvas.Polyline(xs, ys, "stroke: black; fill: none;")
	return true
}

func renderPolygon(canvas *svg.SVG, geometry *geos.Geometry, tile *mapobjects.Tile) bool {
	boundary, err := geometry.Boundary()
	if (err != nil) {
		return false;
	}
	coords, err := boundary.Coords()
	if (err != nil) {
		return false;
	}
	xs := []int{}
	ys := []int{}
	for _, coord := range coords {
		x, y := tile.Degrees2Pixels(coord.Y, coord.X)
		xs = append(xs, x)
		ys = append(ys, y)
	}
	canvas.Polygon(xs, ys, "stroke: black; fill: rgba(100, 100, 100, .5);")
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
	case geos.LINESTRING:
		renderPolyline(canvas, geometry, tile)
	case geos.POLYGON:
		renderPolygon(canvas, geometry, tile)
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