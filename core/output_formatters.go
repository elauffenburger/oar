package core

import (
	"encoding/json"
	"fmt"
)

type OutputFormatter interface {
	Format(results *Results) string
}

type JsonOutputFormatter struct {
}

func (formatter *JsonOutputFormatter) Format(results *Results) string {
	jsonArray := formatter.ToJsonArray(results)
	marshalledbytes, err := json.Marshal(&jsonArray)
	if err != nil {
		panic(fmt.Sprintf("Error marshalling results: '%s", err))
	}

	return string(marshalledbytes)
}

type JsonObject map[string]string
type JsonArray []*JsonObject

func (obj JsonObject) Keys() []string {
	keys := make([]string, len(obj))

	i := 0
	for key, _ := range obj {
		keys[i] = key

		i++
	}

	return keys
}

func (formatter *JsonOutputFormatter) ToJsonArray(results *Results) JsonArray {
	result := make(JsonArray, results.NumRows())

	for i, set := range results.ResultSets {
		object := make(JsonObject)

		for _, entry := range set.Values {
			object[entry.Name] = entry.Value
		}

		result[i] = &object
	}

	return result
}

func GetOutputFormatter(config *Configuration) OutputFormatter {
	return &JsonOutputFormatter{}
}
