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
	point, _ := geos.FromWKT("POINT(0 0)")
	geometries := []geos.Geometry{*point}
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