package core

import (
	"errors"
	"fmt"
)

func (config *Configuration) GetValueForFieldTypeInSet(field ConfigurationField, set *ResultsRow) (string, error) {
	if field.TypeLoader == nil {
		return "", errors.New("Failed to find a loader")
	}

	val, _ := field.TypeLoader.GenerateSingleValue(config, set)
	return fmt.Sprintf("%s", val), nil
}

type ConfigurationField struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	TypeLoader TypeLoader `json:"-"`
}

type ConfigurationFields []*ConfigurationField

func NewConfigurationFields() ConfigurationFields {
	return make(ConfigurationFields, 0)
}

func loadFields(config *Configuration, types map[string]TypeLoader) {
	// load fields from config
	for _, field := range config.Fields {
		loader, ok := types[field.Type]

		if !ok {
			// todo handle failure
			continue
		}

		field.TypeLoader = loader
	}
}
