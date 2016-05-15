package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/TerraFactory/tilegenerator/database"
	"github.com/TerraFactory/tilegenerator/mapobjects"
	"github.com/TerraFactory/tilegenerator/svg"
	"github.com/gorilla/mux"
)

var db database.GeometryDB

func main() {
	runtime.GOMAXPROCS(1) // Temporary workaround
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiles/{z}/{x}/{y}.svg", getTile)
	db = database.GeometryDB{}
	db.InitConnection("postgres", "host=localhost user=postgres dbname=okenit.new sslmode=disable", "maps.maps_objects", "the_geom")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getTile(writer http.ResponseWriter, req *http.Request) {
	var objects []mapobjects.MapObject
	//line, _ := mapobjects.NewObject(
	//3,
	//"LINESTRING (36.6 50.6, 36.183333 51.716667, 36.083333 52.966667)",
	//`polyline {
	//fill: none;
	//stroke: red;
	//}`)

	results, err := db.GetAllFlyRoutes()
	if err != nil {
		fmt.Errorf("err: %s", err.Error())
	} else {
		fmt.Printf("res: %v", results)
		for _, r := range results {
			line, err := mapobjects.NewObject(
				3,
				r.Value,
				`polyline {
	           fill: none;
	           stroke: red;
	         }`)
			if err == nil {
				fmt.Printf("line: %v", line)
				objects = append(objects, *line)
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
