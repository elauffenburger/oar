package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	fields "github.com/elauffenburger/oar/configuration/fields"
)

type Configuration struct {
	Fields     fields.ConfigurationFields    `json:"fields"`
	FieldsData fields.ConfigurationFieldData `json:"-"`
	NumRows    int                           `json:"numRows"`
	Options    map[string]string             `json:"options"`
}

func NewConfiguration() *Configuration {
	return &Configuration{Options: make(map[string]string), Fields: fields.NewConfigurationFields()}
}

func LoadConfigurationFromJson(content string) (*Configuration, error) {
	empty := &Configuration{}
	config := NewConfiguration()

	bytes := []byte(content)
	err := json.Unmarshal(bytes, config)

	if err != nil {
		return empty, errors.New("Failed to parse json")
	}

	// load fields from config
	for _, field := range config.Fields {
		fieldType, err := GetFieldTypeForField(*field)

		if err != nil {
			// TODO handle case where can't find type for field
			continue
		}

		field.FieldType = fieldType
	}

	// hook for options to modify config
	readOptions(config)

	return config, nil
}

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

func readNamesDataFromRelativePath(rel string) ([]string, error) {
	abs, err := filepath.Abs(rel)
	if err != nil {
		return nil, err
	}

	names, err := readContentFromFileAndSplit(abs, "\n")
	if err != nil {
		return nil, err
	}

	return names, nil
}

func readOptions(config *Configuration) {
	options := config.Options

	// get first names
	if firstnamespath, ok := options["firstnames"]; ok {
		if names, err := readNamesDataFromRelativePath(firstnamespath); err == nil {
			config.FieldsData.FirstNames = names
		}
	}

	// get last names
	if lastnamespath, ok := options["lastnames"]; ok {
		if names, err := readNamesDataFromRelativePath(lastnamespath); err == nil {
			config.FieldsData.LastNames = names
		}
	}
}

func LoadConfigurationFromFile(path string) (*Configuration, error) {
	empty := &Configuration{}

	if configPath, err := filepath.Abs(path); err == nil {
		bytes, _ := ioutil.ReadFile(configPath)

		return LoadConfigurationFromJson(string(bytes))
	} else {
		return empty, errors.New(fmt.Sprintf("Could not load file at path '%s'", path))
	}
}

func GetFieldTypeForField(field fields.ConfigurationField) (fields.ConfigurationFieldType, error) {
	converted := fields.ConfigurationFieldType(strings.ToLower(field.Type))
	if len(converted) == 0 {
		return fields.Default, nil
	}

	return converted, nil
}
