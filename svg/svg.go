package svg

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/mapobjects"
)

func prefixSelectors(css string, id int) string {
	reg, err := regexp.Compile("(}?[a-z0-9_ -]{1,256}{)")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(css, "#id"+strconv.Itoa(id)+" $0")
}

//func renderPointObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//coords, err := object.Geometry.Coords()
//if err != nil {
//return err
//}
//x, y := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
//canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
//canvas.CSS(prefixSelectors(object.CSS, object.ID))
//canvas.Circle(x, y, 5, "")
//canvas.Gend()
//return nil
//}

//func renderMultiPointObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//n, err := object.Geometry.NGeometry()
//if err != nil {
//return err
//}
//canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
//canvas.CSS(prefixSelectors(object.CSS, object.ID))
//for i := 0; i < n; i++ {
//g, err := object.Geometry.Geometry(i)
//if err != nil {
//return err
//}
//coords, err := g.Coords()
//if err != nil {
//return err
//}
//x, y := tile.Degrees2Pixels(coords[0].Y, coords[0].X)
//canvas.Circle(x, y, 5, "")

//}
//canvas.Gend()
//return nil
//}

//func renderPolylineObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//coords, err := object.Geometry.Coords()
//if err != nil {
//return err
//}
//xs := []int{}
//ys := []int{}
//for _, coord := range coords {
//x, y := tile.Degrees2Pixels(coord.Y, coord.X)
//xs = append(xs, x)
//ys = append(ys, y)
//}
//canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
//canvas.CSS(prefixSelectors(object.CSS, object.ID))
//canvas.Polyline(xs, ys, "")
//canvas.Gend()
//return nil
//}

//func renderMultiPolylineObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//n, err := object.Geometry.NGeometry()
//if err != nil {
//return err
//}
//canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
//canvas.CSS(prefixSelectors(object.CSS, object.ID))
//for i := 0; i < n; i++ {
//g, err := object.Geometry.Geometry(i)
//if err != nil {
//return err
//}
//coords, err := g.Coords()
//if err != nil {
//return err
//}
//xs := []int{}
//ys := []int{}
//for _, coord := range coords {
//x, y := tile.Degrees2Pixels(coord.Y, coord.X)
//xs = append(xs, x)
//ys = append(ys, y)
//}
//canvas.Polyline(xs, ys, "")

//}
//canvas.Gend()
//return nil
//}

//func renderPolygon(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//boundary, err := object.Geometry.Boundary()
//if err != nil {
//return err
//}
//coords, err := boundary.Coords()
//if err != nil {
//return err
//}
//xs := []int{}
//ys := []int{}
//for _, coord := range coords {
//x, y := tile.Degrees2Pixels(coord.Y, coord.X)
//xs = append(xs, x)
//ys = append(ys, y)
//}
//canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
//canvas.CSS(prefixSelectors(object.CSS, object.ID))
//canvas.Polygon(xs, ys, "")
//canvas.Gend()
//return nil
//}

//func renderMultiPolygonObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//n, err := object.Geometry.NGeometry()
//if err != nil {
//return err
//}
//canvas.Group("id=\"id" + strconv.Itoa(object.ID) + "\"")
//canvas.CSS(prefixSelectors(object.CSS, object.ID))
//for i := 0; i < n; i++ {
//g, err := object.Geometry.Geometry(i)
//if err != nil {
//return err
//}
//boundary, err := g.Boundary()
//if err != nil {
//return err
//}
//coords, err := boundary.Coords()
//if err != nil {
//return err
//}
//xs := []int{}
//ys := []int{}
//for _, coord := range coords {
//x, y := tile.Degrees2Pixels(coord.Y, coord.X)
//xs = append(xs, x)
//ys = append(ys, y)
//}
//canvas.Polygon(xs, ys, "")
//}
//canvas.Gend()
//return nil
//}

//func renderObject(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
//geometryType, err := object.Geometry.Type()
//if err != nil {
//return err
//}
//switch geometryType {
//case geos.POINT:
//renderPointObject(canvas, object, tile)
//case geos.MULTIPOINT:
//renderMultiPointObject(canvas, object, tile)
//case geos.LINESTRING:
//renderPolylineObject(canvas, object, tile)
//case geos.MULTILINESTRING:
//renderMultiPolylineObject(canvas, object, tile)
//case geos.POLYGON:
//renderPolygon(canvas, object, tile)
//case geos.MULTIPOLYGON:
//renderMultiPolygonObject(canvas, object, tile)
//default:
//return errors.New(fmt.Sprintf("Unexpected geometry type: %v", geometryType))
//}
//return nil
//}

// RenderTile takes a tile struct, map objects and then draws these objects on the tile
func RenderTile(tile *mapobjects.Tile, objects *[]mapobjects.MapObject, writer io.Writer) {
	canvas := svg.New(writer)
	canvas.Start(mapobjects.TileSize, mapobjects.TileSize)
	for _, obj := range *objects {
		renderSpecialObjects(canvas, &obj, tile)
	}
	canvas.End()
}

func renderSpecialObjects(canvas *svg.SVG, object *mapobjects.MapObject, tile *mapobjects.Tile) error {
	switch object.TypeID {
	case 47:
		RenderPatrollingArea(canvas, object, tile)
	case 74:
		RenderRouteAviationFlight(canvas, object, tile)
	default:
		return fmt.Errorf("Unexpected geometry type: %v", object.TypeID)
	}
	return nil

	//patrollingArea, _ := mapobjects.NewObject(
	//32,
	//"LINESTRING (70.6 50.6, 16.183333 51.716667)",
	//`polyline, path, line {
	//stroke:black;
	//stroke-width:1;
	//fill: none
	//}`)

	//RenderPatrollingArea(canvas, patrollingArea, tile)

	//routeAviationFlight, _ := mapobjects.NewObject(
	//30,
	//"LINESTRING (10.6 30.6, 46.183333 21.716667)",
	//`line {
	//fill: none;
	//stroke: red;
	//}`)

	//RenderRouteAviationFlight(canvas, routeAviationFlight, tile)

	//routeAviationFlight2, _ := mapobjects.NewObject(
	//30,
	//"LINESTRING ( -46.183333 -21.716667, 10.6 -30.6)",
	//`line {
	//fill: none;
	//stroke: red;
	//}`)

	//RenderRouteAviationFlight(canvas, routeAviationFlight2, tile)

	//beamDiagram := &beamDiagram{radius: 50, angelRotation: 0, sliderBeamWidth: 4, sidelobes: 0}
	//point, _ := mapobjects.NewObject(
	//31,
	//"POINT (-50 30)",
	//`circle {
	//fill: red;
	//}`)

	//RenderBeamDiagram(canvas, point, tile, beamDiagram)

	//point2, _ := mapobjects.NewObject(
	//32,
	//"POINT (-50 -50)",
	//`circle {
	//fill: red;
	//}`)

	//radioModeule1, _ := mapobjects.NewObject(
	//33,
	//"POINT (-49 -49)",
	//`circle {
	//fill: red;
	//}`)

	//radioModeule2, _ := mapobjects.NewObject(
	//34,
	//"POINT (-51 -50)",
	//`circle {
	//fill: red;
	//}`)

	//radioModeule3, _ := mapobjects.NewObject(
	//35,
	//"POINT (-53 -50)",
	//`circle {
	//fill: red;
	//}`)

	//radioModules := []*mapobjects.MapObject{radioModeule1, radioModeule2, radioModeule3}

	//RenderSatelliteVisibility(canvas, point2, radioModules, tile)
}
