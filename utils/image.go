package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
)

var loadedImages = make(map[string][]byte)

func GetImgByURL(url string) ([]byte, error) {
	if val, ok := loadedImages[url]; ok {
		return val, nil
	} else {
		if resp, err := http.Get(url); err == nil {
			loadedImages[url], _ = ioutil.ReadAll(resp.Body)
			return loadedImages[url], nil
		} else {
			return nil, errors.New("Can't load img")
		}
	}
}
