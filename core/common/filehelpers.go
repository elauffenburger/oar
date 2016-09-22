package common

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadContentFromFile(path string) (*string, error) {
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

func ReadContentFromFileAndSplit(path string, sep string) ([]string, error) {
	content, err := ReadContentFromFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(*content, sep), nil
}
