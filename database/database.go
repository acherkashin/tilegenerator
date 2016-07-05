package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/tiles"
	_ "github.com/lib/pq" //we want to use blank import here
)

// GeometryDB is a structure which represents a DB connection
type GeometryDB struct {
	conn      *sql.DB
	geomtable string
	geomcol   string
}

// Transfer raw sql rows into a slice of BaseGeometry structs
func (gdb *GeometryDB) rowsToMapObjects(rows *sql.Rows) ([]entities.MapObject, error) {
	mapObjects := []entities.MapObject{}
	tmpRows := *rows
	defer tmpRows.Close()

	for tmpRows.Next() {
		var ID, typeID int
		var wkt, label, textPosition, colorOuter, colorInner, code string
		var isShortwaveAntenna, needShowAzimuthalGrid, needShowDirectionalDiagram bool
		var sidelobes, beamWidth, azimut, distance float64

		err := tmpRows.Scan(&ID, &typeID, &wkt, &label, &isShortwaveAntenna, &needShowAzimuthalGrid, &beamWidth, &sidelobes, &azimut, &distance, &needShowDirectionalDiagram, &textPosition, &colorOuter, &colorInner, &code)

		if err == nil {
			mapObj, mapObjErr := entities.NewObject(ID, typeID, wkt, isShortwaveAntenna, needShowAzimuthalGrid, needShowDirectionalDiagram, beamWidth, sidelobes, azimut, distance, colorOuter, colorInner, code)
			if mapObjErr == nil {
				mapObj.Label = label
				mapObj.Position = textPosition
				mapObjects = append(mapObjects, *mapObj)
			} else {
				fmt.Println(errors.New("Can't create map object"))
			}
		} else {
			fmt.Println(err)

		}
	}
	return mapObjects, nil
}

// InitConnection creates db connection. Use "geotable" parameter as a table with geometries and
// "geocol" as a geometry column
func (gdb *GeometryDB) InitConnection(username string, connstring string, geomtable string, geomcol string) {
	db, err := sql.Open(username, connstring)
	if err == nil {
		gdb.conn = db
		gdb.geomtable = geomtable
		gdb.geomcol = geomcol
	} else {
		fmt.Printf("Database connection error: %v\n", err)
		panic("DB Error")
	}
}

// Return slice of all geometries in a database
func (gdb *GeometryDB) GetGeometriesForTile(tile *tiles.Tile, situationsIds string) (mapObjects []entities.MapObject, err error) {
	var situationQuery string
	if situationsIds != "" {
		situationQuery = fmt.Sprintf("situation_id in (%v) and", situationsIds)
	}
	q := fmt.Sprintf(`
		SELECT id,type_id, ST_AsText( ST_Transform( %s, 4326 ) ), coalesce(text1, ''), coalesce(is_shortwave_antenna, false),
		coalesce(need_show_azimuthal_grid, false), coalesce(beam_width, '1'), coalesce(sidelobes, '1'), coalesce(azimut, '0'),
		coalesce(distance, '0'), coalesce(need_show_directional_diagram, 'false'), coalesce(text_position, 'bottom'),
		coalesce(color_outer, ''), coalesce(color_inner, ''), coalesce(code, '')  from %s
		WHERE type_id NOT in (170, 11) and 
		(min_zoom <= %v or min_zoom is null) and
		(max_zoom >= %v or max_zoom is null) and %v
		ST_Intersects(ST_SetSRID(ST_MakeBox2D(ST_Point(%v, %v), ST_Point(%v, %v)), 4326), the_geom);
		`, gdb.geomcol, gdb.geomtable, tile.Z, tile.Z, situationQuery, tile.BoundingBox.West, tile.BoundingBox.North, tile.BoundingBox.East, tile.BoundingBox.South)

	rows, err := gdb.conn.Query(q)
	if err == nil {
		mapObjects, scanErr := gdb.rowsToMapObjects(rows)
		return mapObjects, scanErr
	} else {
		fmt.Printf("Query error1: %v", err)
		return nil, err
	}
}

func (gdb *GeometryDB) GetAllSpecialObject(tile *tiles.Tile, situationsIds string) (mapObjects []entities.MapObject, err error) {
	var situationQuery string
	if situationsIds != "" {
		situationQuery = fmt.Sprintf("situation_id in (%v) and", situationsIds)
	}

	q := fmt.Sprintf(`SELECT id,type_id, ST_AsText( ST_Transform( %s, 4326 ) ), coalesce(text1, ''), coalesce(is_shortwave_antenna, false),
		coalesce(need_show_azimuthal_grid, false), coalesce(beam_width, '0'), coalesce(sidelobes, '1'),	coalesce(azimut, '1'),
		coalesce(distance, '0'), coalesce(need_show_directional_diagram, 'false'), coalesce(text_position, 'bottom'),
		coalesce(color_outer, ''), coalesce(color_inner, ''), coalesce(code, '')  from %s 
		WHERE (type_id BETWEEN 149 AND 165) OR (type_id IN (47,74,408,407,366,432)) and
		(min_zoom <= %v or min_zoom is null) and
		(max_zoom >= %v or max_zoom is null) and %v
		ST_Intersects(ST_SetSRID(ST_MakeBox2D(ST_Point(%v, %v), ST_Point(%v, %v)), 4326), the_geom);`, gdb.geomcol, gdb.geomtable, tile.Z, tile.Z, situationQuery, tile.BoundingBox.West, tile.BoundingBox.North, tile.BoundingBox.East, tile.BoundingBox.South)

	rows, err := gdb.conn.Query(q)
	if err == nil {
		mapObjects, scanErr := gdb.rowsToMapObjects(rows)
		return mapObjects, scanErr
	} else {
		fmt.Printf("Query error2: %v", err)
		return nil, err
	}
}
