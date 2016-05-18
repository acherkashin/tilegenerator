package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TerraFactory/tilegenerator/database"
	"github.com/TerraFactory/tilegenerator/geo"
	"github.com/TerraFactory/tilegenerator/mapobjects"
	"github.com/TerraFactory/tilegenerator/settings"
	"github.com/TerraFactory/tilegenerator/svg"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var db database.GeometryDB

func printStartingMsg(config *settings.Settings) {
	fmt.Printf("Starting with the following settings:\n")
	fmt.Printf("\tGeometry table: %s\n", color.CyanString(config.DBGeometryTable))
	fmt.Printf("\tGeometry column: %s\n", color.CyanString(config.DBGeometryColumn))
	fmt.Printf("\tHTTP port: %s\n", color.CyanString(config.HTTPPort))
	color.Green("\n Started!\n")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiles/{z}/{x}/{y}.svg", getTile)
	db = database.GeometryDB{}
	conf := settings.GetSettings()
	db.InitConnection(conf.DBInstanceName, conf.DBConnectionString, conf.DBGeometryTable, conf.DBGeometryColumn)
	printStartingMsg(conf)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.HTTPPort), router))
}

func createMapObject(dbObj geo.BaseGeometry) (*mapobjects.MapObject, error) {
	switch dbObj.TypeID {
	case 47:
		return mapobjects.NewObject(
			dbObj.ID,
			dbObj.TypeID,
			dbObj.Value,
			`polyline, path, line {
					stroke: black;
					stroke-width: 1;
					fill: none
	       }`)
	case 74:
		return mapobjects.NewObject(
			dbObj.ID,
			dbObj.TypeID,
			dbObj.Value,
			`line {
			fill: none;
			stroke: red;
			}`)
	default:
		return nil, fmt.Errorf("Unexpected object type: %v", dbObj)
	}
}

func getTile(writer http.ResponseWriter, req *http.Request) {
	var objects []mapobjects.MapObject
	var ids []int

	results, objErr := db.GetAllPatrollingAreas()
	if objErr != nil {
		fmt.Println(objErr.Error())
		return
	}
	for _, obj := range results {
		ids = append(ids, obj.ID)
	}

	attrs, attrErr := db.GetAllAttributes(ids)
	if attrErr != nil {
		fmt.Println(attrErr.Error())
		return
	}

	for _, r := range results {
		var objAttrs []geo.BaseAttribute
		for _, attr := range attrs {
			if attr.ObjectID == r.ID {
				objAttrs = append(objAttrs, attr)
			}
		}

		r.Attrs = objAttrs

		obj, err := createMapObject(r)
		if err == nil {
			objects = append(objects, *obj)
		} else {
			fmt.Println(err.Error())
			writer.WriteHeader(400)
			return
		}
	}

	vars := mux.Vars(req)
	x, errX := strconv.Atoi(vars["x"])
	y, errY := strconv.Atoi(vars["y"])
	z, errZ := strconv.Atoi(vars["z"])
	if errX != nil || errY != nil || errZ != nil {
		writer.WriteHeader(400)
		return
	}
	fmt.Println("It's still OK")
	tile := mapobjects.NewTile(x, y, z)
	writer.Header().Set("Content-Type", "image/svg+xml")
	svg.RenderTile(tile, &objects, writer)
}
