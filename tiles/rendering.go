package tiles

import (
	"io"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
)

// RenderTile takes a tile struct, map objects and then draws these objects on the tile
func RenderTile(tile *Tile, objects *[]entities.MapObject, writer io.Writer) {
	canvas := svg.New(writer)
	canvas.Start(TileSize, TileSize)
	//for _, _ := range *objects {
	//}
	canvas.End()
}
