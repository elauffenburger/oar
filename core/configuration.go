package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
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

type OutputType string

const (
	JSON OutputType = "json"
	SQL  OutputType = "sql"
)

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

func applyOptions(config *Configuration) {

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
	results := &Results{ResultSets: make(ResultsRowList, numRows)}

	for i := 0; i < numRows; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			set := &ResultsRow{}
			for _, field := range config.Fields {
				entry := &ResultsRowValue{ConfigurationField: *field}

				value, err := config.GetValueForFieldTypeInSet(entry.ConfigurationField, set)
				if err != nil {
					// todo handle error
					continue
				}

				entry.Value = value
				set.Values = append(set.Values, entry)
			}

			results.ResultSets[index] = set
		}(i)
	}

	wg.Wait()
	return results, nil
}

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomValue(list *[]string) string {
	return (*list)[rand.Intn(len(*list))]
}

func init() {
	addDefaultLoaderFactories()
}
