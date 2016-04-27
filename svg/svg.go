package svg

import (
	"github.com/paulsmith/gogeos/geos"
	"github.com/ajstarks/svgo"
	"github.com/terrafactory/tilegenerator/mapobjects"
	"io"
)

func renderPoint(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	coords, err := object.Geometry.Coords()
	if (err != nil) {
		return false;
	}
	x, y := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
	canvas.Circle(x, y, 5, object.CSS)
	return true
}

func renderMultiPoint(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	n, err := object.Geometry.NGeometry()
	if (err != nil) {
		return false
	}
	for i := 0; i < n; i++ {
		g, _ := object.Geometry.Geometry(i)

		coords, err := g.Coords()
		if (err != nil) {
			return false;
		}
		x, y := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
		canvas.Circle(x, y, 5, object.CSS)

	}
	return false;
}

func renderPolyline(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	coords, err := object.Geometry.Coords()
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
	canvas.Polyline(xs, ys, object.CSS)
	return true
}

func renderMultiPolyline(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	n, err := object.Geometry.NGeometry()
	if (err != nil) {
		return false
	}
	for i := 0; i < n; i++ {
		g, _ := object.Geometry.Geometry(i)
		coords, err := g.Coords()
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
		canvas.Polyline(xs, ys, object.CSS)

	}
	return false;
}

func renderPolygon(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	boundary, err := object.Geometry.Boundary()
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
	canvas.Polygon(xs, ys, object.CSS)
	return true
}

func renderMultiPolygon(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	n, err := object.Geometry.NGeometry()
	if (err != nil) {
		return false
	}
	for i := 0; i < n; i++ {
		g, _ := object.Geometry.Geometry(i)
		boundary, err := g.Boundary()
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
		canvas.Polygon(xs, ys, object.CSS)
	}
	return false;
}

func renderObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) bool {
	geometryType, err := object.Geometry.Type()
	if (err != nil) {
		return false
	}
	switch  geometryType{
	case geos.POINT:
		renderPoint(canvas, object, tile)
	case geos.MULTIPOINT:
		renderMultiPoint(canvas, object, tile)
	case geos.LINESTRING:
		renderPolyline(canvas, object, tile)
	case geos.MULTILINESTRING:
		renderMultiPolyline(canvas, object, tile)
	case geos.POLYGON:
		renderPolygon(canvas, object, tile)
	case geos.MULTIPOLYGON:
		renderMultiPolygon(canvas, object, tile)
	default:
		return false
	}
	return true;
}

func RenderTile(tile *mapobjects.Tile, objects *[]mapobjects.MapObject, writer io.Writer) {
	canvas := svg.New(writer)
	canvas.Start(mapobjects.TILE_SIZE, mapobjects.TILE_SIZE)
	for _, geo := range *objects {
		renderObject(canvas, &geo, tile)
	}
	canvas.End()
}