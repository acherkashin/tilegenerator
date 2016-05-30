package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //we want to use blank import here
)

// GeometryDB is a structure which represents a DB connection
type GeometryDB struct {
	conn      *sql.DB
	geomtable string
	geomcol   string
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
