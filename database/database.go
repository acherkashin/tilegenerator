package database

import (
	"database/sql"
	"fmt"
	"github.com/TerraFactory/tilegenerator/geo"
	_ "github.com/lib/pq" //we want to use blank import here
	"strconv"
	"strings"
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
		err := tmpRows.Scan(&geometry.ID, &geometry.TypeID, &geometry.Value)
		if err == nil {
			geometries = append(geometries, geometry)
		}
	}
	return geometries
}

// Transfer raw sql rows into a slice of BaseGeometry structs
func (gdb *GeometryDB) rowsToAttrs(rows *sql.Rows) []geo.BaseAttribute {
	attrs := []geo.BaseAttribute{}
	attr := geo.BaseAttribute{}
	tmpRows := *rows
	defer tmpRows.Close()

	for tmpRows.Next() {
		err := tmpRows.Scan(&attr.Value, &attr.Code, &attr.ObjectID)
		if err == nil {
			attrs = append(attrs, attr)
		}
	}
	return attrs
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
	q := fmt.Sprintf("SELECT id, type_id, ST_AsText( ST_Transform( %s, 4326 ) ) from %s;", gdb.geomcol, gdb.geomtable)
	rows, err := gdb.conn.Query(q)
	if err != nil {
		fmt.Printf("Query error: %v", err)
	} else {
		geometries := gdb.rowsToGeometries(rows)
		return geometries, err
	}
	return
}

// GetAllPatrollingAreas is a tmp method (don't use it in production parts of the app)
func (gdb *GeometryDB) GetAllPatrollingAreas() (geometries []geo.BaseGeometry, err error) {
	q := fmt.Sprintf("SELECT id, type_id, ST_AsText( ST_Transform( %s, 4326 ) ) from %s WHERE type_id = 47 OR type_id = 74;", gdb.geomcol, gdb.geomtable)
	rows, err := gdb.conn.Query(q)
	if err != nil {
		fmt.Printf("Query error: %s\n", err.Error())
	} else {
		geometries := gdb.rowsToGeometries(rows)
		return geometries, err
	}
	return
}

// GetAllAttributes returns all the attributes of an objects
func (gdb *GeometryDB) GetAllAttributes(ids []int) (attrs []geo.BaseAttribute, err error) {
	var tmpIds []string
	for _, id := range ids {
		tmpIds = append(tmpIds, strconv.Itoa(id))
	}

	q := fmt.Sprintf(`
SELECT av.value, a.code_attribute, state.object_id FROM 
maps.object_attribute_values av
inner join maps.object_attributes a on a.id=av.attribute_id 
inner join 
(
	select id, object_id from maps.object_state_histories
	where object_id in (%s) order by id desc limit(%v)
) as state
on state.id=av.object_state_history_id;
	`, strings.Join(tmpIds, ","), len(ids))
	rows, err := gdb.conn.Query(q)
	if err != nil {
		fmt.Printf("Query error: %s\n", err.Error())
	} else {
		attrs := gdb.rowsToAttrs(rows)
		return attrs, err
	}
	return
}
