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
)

const (
	left      = "LEFT"
	leftTop   = "LEFTTOP"
	top       = "TOP"
	rightTop  = "RIGHTTOP"
	right     = "RIGHT"
	rightDown = "RIGHTDOWN"
	down      = "DOWN"
	leftDown  = "LEFTDOWN"
	center    = "CENTER"
)

type ImagePrimitive struct {
	Width  int64
	Height int64
	Scale  float64
	Href   string
	Rotate float64
	Format string
	bytes  []byte
}

func (img ImagePrimitive) Render(svg *svg.SVG, object *entities.MapObject) {
	point, _ := object.Geometry.AsPoint()
	resultHref := strings.Replace(img.Href, "${ID}", strconv.Itoa(object.ID), 1)

	img.Rotate = object.AzimuthalGrid.Azimut
	img.Scale = object.View.Scale
	tmpH := float64(img.Height)
	tmpW := float64(img.Width)
	img.Width, img.Height = int64(math.Floor(tmpW*img.Scale)), int64(math.Floor(tmpH*img.Scale))

	if result, err := utils.GetImgByURL(resultHref); err == nil {
		img.bytes = result
		inlineBase64Img := base64.StdEncoding.EncodeToString(img.bytes)
		svg.TranslateRotate(
			int(math.Floor(point.Coordinates.X+.5)),
			int(math.Floor(point.Coordinates.Y+0.5)),
			img.Rotate)

		if object.View.NeedMirrorReflection {
			svg.CSS(fmt.Sprintf("#id%v { transform: scale(-1, 1) }", object.ID))
		}

		svg.Gid(fmt.Sprintf("id%v", strconv.Itoa(object.ID)))

		shiftX := getHorizontalShift(img, object.MarkerPosition)
		shiftY := getVerticalShift(img, object.MarkerPosition)

		svg.Image(int(shiftX), int(shiftY), int(img.Width), int(img.Height), "data:"+img.Format+";base64,"+inlineBase64Img)
		svg.Gend()
		svg.Gend()
		svg.Circle(int(point.Coordinates.X), int(point.Coordinates.Y), int(6), "fill:red;stroke:red")
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

func getHorizontalShift(img ImagePrimitive, position string) (shift int64) {
	switch strings.ToUpper(position) {
	case leftDown, left, leftTop:
		shift = -img.Width
	case rightDown, right, rightTop:
		shift = 0
	case down, top:
		shift = -img.Width / 2
	default:
		shift = -img.Width / 2
	}

	return shift
}

func getVerticalShift(img ImagePrimitive, position string) (shift int64) {
	switch strings.ToUpper(position) {
	case leftTop, top, rightTop:
		shift = -img.Height
	case leftDown, down, rightDown:
		shift = 0
	case left, right:
		shift = -img.Height / 2
	default:
		shift = -img.Height / 2
	}

	return shift
}

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}
