package tiles

import (
	"io"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/settings/styling"
)

// RenderTile takes a tile struct, map objects and then draws these objects on the tile
func RenderTile(tile *Tile, objects *[]entities.MapObject, styles *map[string]styling.Style, writer io.Writer) {
	f := func(x, y float64) (float64, float64) {
		nx, ny := tile.Degrees2Pixels(y, x)
		return float64(nx), float64(ny)
	}

	canvas := svg.New(writer)
	canvas.Start(TileSize, TileSize)
	for _, object := range *objects {
		object.Geometry.ConvertCoords(f)
		for _, style := range *styles {
			if style.ShouldRender(&object) {
				style.Render(&object, canvas)
			}
		}
	}
	canvas.End()
}
