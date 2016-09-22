package configuration

import (
	"encoding/json"
	"fmt"
)

type ResultSetList []*ResultSet
type Results struct {
	ResultSets ResultSetList
}

type EntryRow []*ResultSetEntry
type ResultSet struct {
	Entries EntryRow
}

type ResultSetEntry struct {
	ConfigurationField
	Value string
}

func (entries *EntryRow) GetEntryWithName(name string) (*ResultSetEntry, error) {
	for _, entry := range *entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("Failed to find an entry with name '%s'", name)
}

func (entries *EntryRow) GetEntryWithType(entryType string) (*ResultSetEntry, error) {
	for _, entry := range *entries {
		if entry.Type == entryType {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("Failed to find an entry with type '%s'", entryType)
}

func NewRowsOfSize(n int) []EntryRow {
	return make([]EntryRow, n)
}

func NewRows() []EntryRow {
	return NewRowsOfSize(0)
}

func (results *Results) NumRows() int {
	return len(results.ResultSets)
}

func (results *Results) NumEntries() int {
	numrows := results.NumRows()

	if numrows == 0 {
		return 0
	}

	return len(results.ResultSets[0].Entries) * numrows
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

func (results *Results) ToJsonArray() JsonArray {
	result := make(JsonArray, results.NumRows())

	for i, set := range results.ResultSets {
		object := make(JsonObject)

		for _, entry := range set.Entries {
			object[entry.Name] = entry.Value
		}

		result[i] = &object
	}

	return result
}

func (results *Results) ToJson() string {
	jsonArray := results.ToJsonArray()
	marshalledbytes, err := json.Marshal(&jsonArray)
	if err != nil {
		panic(fmt.Sprintf("Error marshalling results: '%s", err))
	}

	return string(marshalledbytes)
}

func (results *Results) ToRows() []EntryRow {
	rows := NewRowsOfSize(results.NumRows())

	for i, set := range results.ResultSets {
		rows[i] = set.Entries
	}

	return rows
}
