package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	fields "github.com/elauffenburger/oar/configuration/fields"
)

func LoadConfigurationFromJson(content string) (*Configuration, error) {
	empty := &Configuration{}
	config := &Configuration{}

	bytes := []byte(content)
	err := json.Unmarshal(bytes, config)

	if err != nil {
		return empty, errors.New("Failed to parse json")
	}

	for i, _ := range config.Fields {
		field := config.Fields[i]
		fieldType, err := GetFieldTypeForField(*field)

		if err != nil {
			// TODO handle case where can't find type for field
			continue
		}

		field.FieldType = fieldType
	}

	return config, nil
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

type Configuration struct {
	Fields  []*fields.ConfigurationField `json:"fields"`
	NumRows int                          `json:"numRows"`
}
