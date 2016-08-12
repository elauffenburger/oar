package configuration

import (
	"fmt"
	"math/rand"
	"time"

	"strconv"

	"github.com/elauffenburger/oar/configuration/fields"
)

type Results struct {
	ResultSets []*ResultSet
}

type EntryRow []*ResultSetEntry
type ResultSet struct {
	Entries EntryRow
}

type ResultSetEntry struct {
	fields.ConfigurationField
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

func (entries *EntryRow) GetEntryWithType(entryType fields.ConfigurationFieldType) (*ResultSetEntry, error) {
	for _, entry := range *entries {
		if entry.FieldType == entryType {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("Failed to find an entry with type '%s'", entryType)
}

func (entries *EntryRow) HasFirstAndLastNames() bool {
	_, firsterror := entries.GetEntryWithType(fields.FirstName)
	_, lasterror := entries.GetEntryWithType(fields.LastName)

	return firsterror == nil && lasterror == nil
}

func (entries *EntryRow) GetFirstAndLastNames() (string, string) {
	first, _ := entries.GetEntryWithType(fields.FirstName)
	last, _ := entries.GetEntryWithType(fields.LastName)

	return first.Value, last.Value
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

func (results *Results) AsRows() []EntryRow {
	rows := NewRowsOfSize(results.NumRows())

	for i, set := range results.ResultSets {
		rows[i] = set.Entries
	}

	return rows
}

func (config *Configuration) GenerateResults() (*Results, error) {
	results := &Results{ResultSets: make([]*ResultSet, config.NumRows)}

	for i := 0; i < config.NumRows; i++ {
		set := &ResultSet{}

		for _, field := range config.Fields {
			entry := &ResultSetEntry{ConfigurationField: *field}
			entry.Value = config.GetValueForFieldTypeInSet(entry.FieldType, set)

			set.Entries = append(set.Entries, entry)
		}

		results.ResultSets[i] = set
	}

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
