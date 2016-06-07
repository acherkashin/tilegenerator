package listeners

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TerraFactory/tilegenerator/database"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings"
	"github.com/TerraFactory/tilegenerator/settings/styling"
	"github.com/TerraFactory/tilegenerator/tiles"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var db database.GeometryDB
var styles *map[string]styling.Style

func printStartingMsg(config *settings.Settings) {
	fmt.Printf("Starting with the following settings:\n")
	fmt.Printf("\tGeometry table: %s\n", color.CyanString(config.DBGeometryTable))
	fmt.Printf("\tGeometry column: %s\n", color.CyanString(config.DBGeometryColumn))
	fmt.Printf("\tHTTP port: %s\n", color.CyanString(config.HTTPPort))
	color.Green("\n Started!\n")
}

func getTile(writer http.ResponseWriter, req *http.Request) {
	objects := []entities.MapObject{}
	/* hardcode for test */
	o1, err1 := entities.NewObject(1, "POINT (0 0)")
	o2, err2 := entities.NewObject(2, "POINT (10 10)")
	if err1 == nil && err2 == nil {
		o1.StyleName = "home"
		o2.StyleName = "mil/airbase"
		objects = append(objects, *o1)
		objects = append(objects, *o2)
		/* hardcode for test end*/

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
		tiles.RenderTile(tile, &objects, styles, writer)
	} else {
		writer.WriteHeader(400)
	}
}

func StartApplication(conf *settings.Settings) {
	/* connect to DB */
	/* pool of connections needed here later. */
	db = database.GeometryDB{}
	db.InitConnection(conf.DBInstanceName, conf.DBConnectionString, conf.DBGeometryTable, conf.DBGeometryColumn)

	/* Read styles from file system */
	styles, _ = styling.GetStyles(conf)

	/* Create router and start listening */
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiles/{z}/{x}/{y}.svg", getTile)
	printStartingMsg(conf)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.HTTPPort), router))
}
