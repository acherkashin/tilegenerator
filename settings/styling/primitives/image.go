package primitives

import (
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/TerraFactory/svgo"
	"github.com/TerraFactory/wktparser/geometry"
	"encoding/base64"
)

type ImagePrimitive struct {
	Width  int64
	Height int64
	Href   string
	Rotate float64
	bytes  []byte
}

func (img ImagePrimitive) Render(svg *svg.SVG, geo geometry.Geometry) {
	point, _ := geo.AsPoint()
	inlineBase64Img := base64.StdEncoding.EncodeToString(img.bytes)
	svg.TranslateRotate(
		int(math.Floor(point.Coordinates.X + .5)),
		int(math.Floor(point.Coordinates.Y + 0.5)),
		img.Rotate)
	svg.Image(0, 0, int(img.Width), int(img.Height), "data:image/png;base64," + inlineBase64Img)
	svg.Gend()
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
			resp, _ := http.Get(img.Href)
			img.bytes, _ = ioutil.ReadAll(resp.Body) // Must implement reading bytes from files.
		case "ROTATE":
			img.Rotate = value.(float64)
		}
	}

	return img, nil
}
