package svg

import (
	"io"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/mapobjects"
)

// RenderTile takes a tile struct, map objects and then draws these objects on the tile
func RenderTile(tile *mapobjects.Tile, objects *[]mapobjects.MapObject, writer io.Writer) {
	canvas := svg.New(writer)
	canvas.Start(mapobjects.TileSize, mapobjects.TileSize)
	//for _, _ := range *objects {
	//}
	canvas.End()
}
