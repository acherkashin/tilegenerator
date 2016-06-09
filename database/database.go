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
		var ID int
		var wkt string
		err := tmpRows.Scan(&ID, &wkt)
		if err == nil {
			mapObj, mapObjErr := entities.NewObject(ID, wkt)
			if mapObjErr == nil {
				mapObjects = append(mapObjects, *mapObj)
			} else {
				fmt.Println(errors.New("Can't create map object"))
			}
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
func (gdb *GeometryDB) GetGeometriesForTile(tile *tiles.Tile) (mapObjects []entities.MapObject, err error) {
	q := fmt.Sprintf(`

		SELECT id, ST_AsText( ST_Transform( %s, 4326 ) ) from %s
		where type_id not in (170, 11) and
		ST_Contains(ST_SetSRID(ST_MakeBox2D(ST_Point(%v, %v), ST_Point(%v, %v)), 4326), the_geom);
		
		`, gdb.geomcol, gdb.geomtable, tile.BoundingBox.West, tile.BoundingBox.North, tile.BoundingBox.East, tile.BoundingBox.South)

	rows, err := gdb.conn.Query(q)
	if err == nil {
		mapObjects, scanErr := gdb.rowsToMapObjects(rows)
		return mapObjects, scanErr
	} else {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}
}
