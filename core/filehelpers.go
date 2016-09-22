package core

import (
	"io/ioutil"
	"os"
	"strings"
)

func readContentFromFile(path string) (*string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	str := string(bytes)
	return &str, nil
}

func readContentFromFileAndSplit(path string, sep string) ([]string, error) {
	content, err := readContentFromFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(*content, sep), nil
}
