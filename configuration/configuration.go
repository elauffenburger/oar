package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sync"
)

type UseTypeDTO struct {
	LoaderArgs UseTypeLoaderArgsDTO `json:"loader"`
}

type UseTypeLoaderArgsDTO struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type Configuration struct {
	Fields  ConfigurationFields   `json:"fields"`
	NumRows int                   `json:"numRows"`
	Options map[string]string     `json:"options"`
	Types   map[string]UseTypeDTO `json:"types"`
}

func NewConfiguration() *Configuration {
	return &Configuration{Options: make(map[string]string), Fields: NewConfigurationFields()}
}

func LoadConfigurationFromJson(content string) (*Configuration, error) {
	empty := &Configuration{}
	config := NewConfiguration()

	bytes := []byte(content)
	err := json.Unmarshal(bytes, config)

	if err != nil {
		return empty, errors.New("Failed to parse json")
	}

	// hook for options to modify config
	applyOptions(config)

	// load custom types
	types := generateTypeLoaders(config)

	// load fields
	loadFields(config, types)

	return config, nil
}

func generateTypeLoaders(config *Configuration) map[string]*TypeLoader {
	types := make(map[string]*TypeLoader)

	// load custom types
	for typename, t := range config.Types {
		if _, exists := types[typename]; exists {
			continue
		}

		loadername := t.LoaderArgs.Name
		factory, ok := loaderFactories[loadername]

		if !ok {
			// todo handle unknown loader
			continue
		}

		// create a loader for this type
		loader := factory()
		loader.Load(&t)

		types[typename] = loader
	}

	return types
}

func applyOptions(config *Configuration) {

}

func loadFields(config *Configuration, types map[string]*TypeLoader) {
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

func LoadConfigurationFromFile(path string) (*Configuration, error) {
	empty := &Configuration{}

	if configPath, err := filepath.Abs(path); err == nil {
		bytes, _ := ioutil.ReadFile(configPath)

		return LoadConfigurationFromJson(string(bytes))
	} else {
		return empty, errors.New(fmt.Sprintf("Could not load file at path '%s'", path))
	}
}

func (config *Configuration) GenerateResults() (*Results, error) {
	var wg sync.WaitGroup

	numRows := config.NumRows
	results := &Results{ResultSets: make(ResultSetList, numRows)}

	for i := 0; i < numRows; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			set := &ResultSet{}
			for _, field := range config.Fields {
				entry := &ResultSetEntry{ConfigurationField: *field}

				value, err := config.GetValueForFieldTypeInSet(entry.ConfigurationField, set)
				if err != nil {
					// todo handle error
					continue
				}

				entry.Value = value
				set.Entries = append(set.Entries, entry)
			}

			results.ResultSets[index] = set
		}(i)
	}

	wg.Wait()
	return results, nil
}

func (config *Configuration) GetValueForFieldTypeInSet(field ConfigurationField, set *ResultSet) (string, error) {
	if field.TypeLoader == nil {
		return "", errors.New("Failed to find a loader")
	}

	val, _ := field.TypeLoader.GenerateSingleValue(config, set)
	return fmt.Sprintf("%s", val), nil
}

type ConfigurationField struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	TypeLoader *TypeLoader `json:"-"`
}

type ConfigurationFields []*ConfigurationField

func NewConfigurationFields() ConfigurationFields {
	return make(ConfigurationFields, 0)
}

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomValue(list *[]string) string {
	n := len(*list)
	randomint := src.Int()

	return (*list)[randomint%n]
}

func init() {
	addDefaultLoaderFactories()
}
