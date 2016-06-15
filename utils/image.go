package utils

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

var loadedImages = make(map[string][]byte)

func GetImgByURL(url string) ([]byte, error) {
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

func GetImgFromFile(path string) ([]byte, error) {
	imgFile, err := os.Open(path)

	if err == nil {

		defer imgFile.Close()

		// create a new buffer base on file size
		fInfo, _ := imgFile.Stat()
		size := fInfo.Size()
		buf := make([]byte, size)

		// read file content into buffer
		fReader := bufio.NewReader(imgFile)
		fReader.Read(buf)

		return buf, err
	}

	return nil, err
}
