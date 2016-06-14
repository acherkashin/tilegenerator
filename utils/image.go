package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
)

var loadedImages = make(map[string][]byte)
var mutex sync.Mutex

func GetImgByURL(url string) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if val, ok := loadedImages[url]; ok {
		return val, nil
	} else {
		if resp, err := http.Get(url); err == nil {
			if result, readErr := ioutil.ReadAll(resp.Body); readErr == nil {
				loadedImages[url] = result
				return result, nil
			} else {
				return nil, errors.New("Can't read bytes from image loading response")
			}
		} else {
			return nil, errors.New("Can't load img")
		}
	}
}
