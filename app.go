package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/TerraFactory/tilegenerator/svg"
	"github.com/TerraFactory/tilegenerator/mapobjects"
	"strconv"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1) // Temporary workaround
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiles/{z}/{x}/{y}.svg", GetTile)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetTile(writer http.ResponseWriter, req *http.Request) {
	point, _ := mapobjects.NewObject(
		1,
		"POINT (0 0)",
		`circle {
		   fill: red;
		 }`)

	multipoint, _ := mapobjects.NewObject(
		2,
		"MULTIPOINT ((37.617778 55.755833), (30.316667 59.95), (33.533333 44.6))",
		`circle {
		   fill: blue;
		 }`)

	line, _ := mapobjects.NewObject(
		3,
		"LINESTRING (36.6 50.6, 36.183333 51.716667, 36.083333 52.966667)",
		`polyline {
	           fill: none;
	           stroke: red;
	         }`)

	multiline, _ := mapobjects.NewObject(
		4,
		"MULTILINESTRING ((10 10, 20 20, 10 40),(40 40, 30 30, 40 20, 30 10))",
		`polyline {
		  fill: none;
		  stroke: red
		}`)

	poly, _ := mapobjects.NewObject(
		5,
		"POLYGON ((-30 -10, -40 -40, -20 -40, -10 -20, -30 -10))",
		`polygon {
		  fill: rgba(100, 100, 100, .1);
		  stroke: black
		}`)

	multipoly, _ := mapobjects.NewObject(
		6,
		"MULTIPOLYGON (((10 -10, 40 -10, 40 -40, 10 -40, 10 -10)),((15 -5, 40 -10, 10 -20, -5 10, 15 -5)))",
		`polygon {
		  fill: rgba(100, 100, 100, .1);
		  stroke: black
		}`)

	objects := []mapobjects.MapObject{*point, *multipoint, *line, *multiline, *poly, *multipoly}

	vars := mux.Vars(req)
	x, errX := strconv.Atoi(vars["x"])
	y, errY := strconv.Atoi(vars["y"])
	z, errZ := strconv.Atoi(vars["z"])
	if (errX != nil || errY != nil || errZ != nil) {
		writer.WriteHeader(400)
		return
	}
	tile := mapobjects.NewTile(x, y, z)
	writer.Header().Set("Content-Type", "image/svg+xml")
	svg.RenderTile(tile, &objects, writer)
}