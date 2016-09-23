package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	conf "github.com/elauffenburger/oar/core/configuration"
	"github.com/elauffenburger/oar/core/loaders"
	"github.com/elauffenburger/oar/core/output"
	res "github.com/elauffenburger/oar/core/results"
)

func LoadConfigurationFromJson(content string) (*conf.Configuration, error) {
	empty := &conf.Configuration{}
	config := conf.NewConfiguration()

	bytes := []byte(content)
	err := json.Unmarshal(bytes, config)

	if err != nil {
		return empty, errors.New("Failed to parse json")
	}

	// hook for options to modify config
	applyOptions(config)

	return config, nil
}

func applyOptions(config *conf.Configuration) {

}

func LoadConfigurationFromFile(path string) (*conf.Configuration, error) {
	empty := &conf.Configuration{}

	if configPath, err := filepath.Abs(path); err == nil {
		bytes, _ := ioutil.ReadFile(configPath)

		return LoadConfigurationFromJson(string(bytes))
	} else {
		return empty, errors.New(fmt.Sprintf("Could not load file at path '%s'", path))
	}
}

func GenerateResults(config *conf.Configuration) (*res.Results, error) {
	return GenerateResultsWithTypeLoaderContext(config, NewTypeLoaderFactoryContext())
}

func GenerateResultsWithTypeLoaderContext(config *conf.Configuration, loaderFactoryContext *loaders.TypeLoaderFactoryContext) (*res.Results, error) {
	var wg sync.WaitGroup

	numRows := config.NumRows
	results := &res.Results{Rows: make(res.ResultsRowList, numRows)}

	// generate type loaders to fulfill this config
	types := BuildTypeLoadersForConfig(config, loaderFactoryContext)

	for i := 0; i < numRows; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			set := &res.ResultsRow{}
			for _, field := range config.Fields {
				entry := &res.ResultsRowValue{ConfigurationField: *field}

				value, err := GenerateValueForField(config, entry.ConfigurationField, set, types)
				if err != nil {
					// todo handle error
					continue
				}

				entry.Value = value
				set.Values = append(set.Values, entry)
			}

			results.Rows[index] = set
		}(i)
	}

	wg.Wait()
	return results, nil
}

func GenerateValueForField(config *conf.Configuration, field conf.ConfigurationField, set *res.ResultsRow, loaders map[string]loaders.TypeLoader) (string, error) {
	loader, ok := loaders[field.Type]

	if !ok {
		return "", errors.New("Failed to find a loader")
	}

	val, _ := loader.GenerateSingleValue(config, set)
	return fmt.Sprintf("%s", val), nil
}

func GetOutputFormatter(config *conf.Configuration) output.OutputFormatter {
	outputtype := config.OutputType

	switch outputtype {
	case conf.JSON:
		return &output.JsonOutputFormatter{}
	case conf.SQL:
		return output.NewSqlOutputFormatter(config.Name)
	}

	panic("Couldn't figure out which output formatter to use")
}

func BuildTypeLoadersForConfig(config *conf.Configuration, typeLoaderFactoryCtx *loaders.TypeLoaderFactoryContext) map[string]loaders.TypeLoader {
	types := make(map[string]loaders.TypeLoader)

	// load custom types
	for typename, t := range config.Types {
		if _, exists := types[typename]; exists {
			continue
		}

		loadername := t.LoaderArgs.Name
		factory, ok := (*typeLoaderFactoryCtx)[loadername]

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

func NewTypeLoaderFactoryContext() *loaders.TypeLoaderFactoryContext {
	ctx := make(loaders.TypeLoaderFactoryContext)

	loaders.AddDefaultLoaderFactories(&ctx)

	return &ctx
}
