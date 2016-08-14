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

	"strconv"

	"sync"

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
				entry.Value = config.GetValueForFieldTypeInSet(entry.FieldType, set)

				set.Entries = append(set.Entries, entry)
			}

			results.ResultSets[index] = set
		}(i)
	}

	wg.Wait()
	return results, nil
}

func (config *Configuration) GetValueForFieldTypeInSet(fieldType fields.ConfigurationFieldType, set *ResultSet) string {
	switch fieldType {
	case fields.FirstName:
		return config.NewFirstName()
	case fields.LastName:
		return config.NewLastName()
	case fields.FullName:
		entries := set.Entries

		if entries.HasFirstAndLastNames() {
			return config.NewFullNameFromNames(entries.GetFirstAndLastNames())
		}
	case fields.String:
		return ""
	case fields.Number:
		return strconv.FormatInt(config.NewNumber(), 10)
	case fields.DateTime:
		return config.NewDateTime().Format(time.RFC1123Z)
	}

	return ""
}

func (config *Configuration) NewFirstName() string {
	return config.FieldsData.FirstNames.GetRandomValue()
}

func (config *Configuration) NewLastName() string {
	return config.FieldsData.LastNames.GetRandomValue()
}

func (config *Configuration) NewFullName() string {
	return config.NewFullNameFromNames(config.NewFirstName(), config.NewLastName())
}

func (config *Configuration) NewNumber() int64 {
	return rand.Int63()
}

func (config *Configuration) NewDateTime() time.Time {
	return time.Unix(rand.Int63(), rand.Int63())
}

func (config *Configuration) NewFullNameFromNames(firstName string, lastName string) string {
	return fmt.Sprintf("%s %s", firstName, lastName)
}
