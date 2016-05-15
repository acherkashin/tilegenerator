package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TerraFactory/tilegenerator/database"
	"github.com/TerraFactory/tilegenerator/geo"
	"github.com/TerraFactory/tilegenerator/mapobjects"
	"github.com/TerraFactory/tilegenerator/svg"
	"github.com/gorilla/mux"
)

var db database.GeometryDB

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiles/{z}/{x}/{y}.svg", getTile)
	db = database.GeometryDB{}
	db.InitConnection("postgres", "host=localhost user=postgres dbname=okenit.new sslmode=disable", "maps.maps_objects", "the_geom")
	fmt.Println("Server has been started on 'localhost:8000'")
	log.Fatal(http.ListenAndServe(":8000", router))
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

	results, err := db.GetAllPatrollingAreas()
	if err != nil {
		fmt.Errorf("err: %s", err.Error())
	} else {
		for _, r := range results {
			obj, err := createMapObject(r)
			if err == nil {
				objects = append(objects, *obj)
			} else {
				fmt.Errorf("object creation err: %s", err.Error())
			}
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
	tile := mapobjects.NewTile(x, y, z)
	writer.Header().Set("Content-Type", "image/svg+xml")
	svg.RenderTile(tile, &objects, writer)
}
