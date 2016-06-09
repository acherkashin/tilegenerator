package primitives

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"encoding/base64"
	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/tilegenerator/database/entities"
	"github.com/TerraFactory/tilegenerator/utils"
	"github.com/TerraFactory/wktparser/geometry"
)

type ImagePrimitive struct {
	Width  int64
	Height int64
	Href   string
	Rotate float64
	Format string
	bytes  []byte
}

func (img ImagePrimitive) Render(svg *svg.SVG, geo geometry.Geometry, object *entities.MapObject) {
	point, _ := geo.AsPoint()
	resultHref := strings.Replace(img.Href, "${ID}", strconv.Itoa(object.ID), 1)
	if result, err := utils.GetImgByURL(resultHref); err == nil {
		img.bytes = result
		inlineBase64Img := base64.StdEncoding.EncodeToString(img.bytes)
		svg.TranslateRotate(
			int(math.Floor(point.Coordinates.X+.5)),
			int(math.Floor(point.Coordinates.Y+0.5)),
			img.Rotate)
		svg.Image(-int(img.Width)/2, -int(img.Height)/2, int(img.Width), int(img.Height), "data:"+img.Format+";base64,"+inlineBase64Img)
		svg.Gend()
	} else {
		fmt.Printf("Can't render %s because of err: '%s'", resultHref, err.Error())
	}
}

func NewImagePrimitive(params *map[string]interface{}) (ImagePrimitive, error) {
	img := ImagePrimitive{}
	for key, value := range *params {
		switch strings.ToUpper(key) { // Switch here is temporary workaround. I should use reflect instead.
		case "WIDTH":
			img.Width = value.(int64)
		case "HEIGHT":
			img.Height = value.(int64)
		case "HREF":
			img.Href = value.(string)
		case "ROTATE":
			img.Rotate = value.(float64)
		case "FORMAT":
			img.Format = value.(string)
		}
	}

	return img, nil
}
