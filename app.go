package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/terrafactory/tilegenerator/svg"
	"github.com/terrafactory/tilegenerator/mapobjects"
	"strconv"
	"github.com/paulsmith/gogeos/geos"
	"log"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiles/{z}/{x}/{y}.svg", GetTile)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetTile(writer http.ResponseWriter, req *http.Request) {
	point, _ := geos.FromWKT("POINT (0 0)")
	multipoint, _ := geos.FromWKT("MULTIPOINT ((10 40), (40 30), (20 20), (30 10))")
	line, _ := geos.FromWKT("LINESTRING (0 0, 20 10, 10 10, 20 20)")
	multiline, _ := geos.FromWKT("MULTILINESTRING ((10 10, 20 20, 10 40),(40 40, 30 30, 40 20, 30 10))")
	poly, _ := geos.FromWKT("POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10))")
	multipoly, _ := geos.FromWKT("MULTIPOLYGON (((10 10, 40 10, 40 40, 10 40, 10 10)),((15 5, 40 10, 10 20, 5 10, 15 5)))")
	geometries := []geos.Geometry{*point, *multipoint, *line, *multiline, *poly, *multipoly}

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
	svg.RenderTile(tile, &geometries, writer)
}