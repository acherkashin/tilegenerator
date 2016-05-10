package database

import (
	"database/sql"
	"fmt"
	"github.com/TerraFactory/tilegenerator/geo"
	"github.com/lib/pq"
)

// GeometryDB is a structure which represents a DB connection
type GeometryDB struct {
	conn      *sql.DB
	geomtable string
	geomcol   string
}

// Transfer raw sql rows into a slice of BaseGeometry structs
func (gdb *GeometryDB) rowsToGeometries(rows *sql.Rows) []geo.BaseGeometry {
	geometries := []geo.BaseGeometry{}
	geometry := geo.BaseGeometry{}
	tmpRows := *rows
	defer tmpRows.Close()

	for tmpRows.Next() {
		err := tmpRows.Scan(&geometry.Id, &geometry.Value)
		if err == nil {
			geometries = append(geometries, geometry)
		}
	}
	return geometries
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

// GetAllGeometries returns a slice of all geometries in a database
func (gdb *GeometryDB) GetAllGeometries() (geometries []geo.BaseGeometry, err error) {
	q := fmt.Sprintf("SELECT id, ST_AsText( ST_Transform( %s, 4326 ) ) from %s;", gdb.geomcol, gdb.geomtable)
	rows, err := gdb.conn.Query(q)
	if err != nil {
		fmt.Printf("Query error: %v", err)
	} else {
		geometries := gdb.rowsToGeometries(rows)
		return geometries, err
	}
	return
}
