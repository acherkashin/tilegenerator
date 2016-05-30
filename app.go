package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TerraFactory/tilegenerator/database"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/TerraFactory/tilegenerator/tiles"
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


func getTile(writer http.ResponseWriter, req *http.Request) {
	var objects []entities.MapObject

	vars := mux.Vars(req)
	x, errX := strconv.Atoi(vars["x"])
	y, errY := strconv.Atoi(vars["y"])
	z, errZ := strconv.Atoi(vars["z"])
	if errX != nil || errY != nil || errZ != nil {
		writer.WriteHeader(400)
		return
	}

	tile := tiles.NewTile(x, y, z)
	writer.Header().Set("Content-Type", "image/svg+xml")
	tiles.RenderTile(tile, &objects, writer)
}
