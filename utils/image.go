package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetImgByURL(url string) ([]byte, error) {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if result, readErr := ioutil.ReadAll(resp.Body); readErr == nil {
			return result, nil
		} else {
			return nil, errors.New("Can't read bytes from image loading response")
		}
	} else {
		return nil, errors.New("Can't load img")
	}
}

func SaveImageToFile(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("Can't save image, path: %v", path))
		return err
	}
	defer file.Close()

	file.Write(content)

	return nil
}

func GetImgFromFile(path string) ([]byte, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
