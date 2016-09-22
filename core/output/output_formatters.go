package output

import (
	"encoding/json"
	"fmt"

	res "github.com/elauffenburger/oar/core/results"
)

type OutputFormatter interface {
	Format(results *res.Results) string
}

type JsonOutputFormatter struct {
}

func (formatter *JsonOutputFormatter) Format(results *res.Results) string {
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

func (formatter *JsonOutputFormatter) ToJsonArray(results *res.Results) JsonArray {
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
