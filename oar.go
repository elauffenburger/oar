package main

import (
	"flag"

	oarconfig "github.com/elauffenburger/oar/configuration"
	fields "github.com/elauffenburger/oar/configuration/fields"
)

type Results struct {
	ResultSets []*ResultSet
	NumRows    int
}

type ResultSet struct {
	Entries []*ResultSetEntry
}

type ResultSetEntry struct {
	fields.ConfigurationField
	Value string
}

var configFlag = flag.String("config", "", "json file to load configuration from")
var firstNamesFlag = flag.String("firstnames", "", "csv of first names")
var lastNamesFlag = flag.String("lastnames", "", "csv of last names")

func GetValueForFieldType(fieldType fields.ConfigurationFieldType) string {
	return ""
}

func GenerateResults(config oarconfig.Configuration) (*Results, error) {
	results := &Results{ResultSets: make([]*ResultSet, config.NumRows)}

	for i := 0; i < config.NumRows; i++ {
		set := &ResultSet{}

		for _, field := range config.Fields {
			entry := &ResultSetEntry{ConfigurationField: *field}
			entry.Value = GetValueForFieldType(entry.FieldType)

			set.Entries = append(set.Entries, entry)
		}

		results.ResultSets = append(results.ResultSets, set)
	}

	return results, nil
}

func main() {
	configFlagValue := *configFlag

	if len(configFlagValue) == 0 {
		panic("No configuration file provided!")
	}

	config, err := oarconfig.LoadConfigurationFromFile(configFlagValue)
	if err != nil {
		panic("Error loading configuration file")
	}

	GenerateResults(*config)
}
