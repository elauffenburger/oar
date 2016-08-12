package main

import (
	"testing"

	oarconfig "github.com/elauffenburger/oar/configuration"
)

func TestCanLoadFromFile(t *testing.T) {
	config, err := oarconfig.LoadConfigurationFromFile("./data/test.json")

	if err != nil {
		t.Fail()
	}

	results, err := config.GenerateResults()

	if err != nil {
		t.Fail()
	}

	if results.NumRows() != len(results.ResultSets) {
		t.Fail()
	}

	if results.NumRows() != config.NumRows {
		t.Fail()
	}
}

func TestCanGenerateNames(t *testing.T) {
	config, _ := oarconfig.LoadConfigurationFromFile("./data/test.json")

	if len(config.FieldsData.FirstNames) == 0 {
		t.Fail()
	}

	if len(config.FieldsData.LastNames) == 0 {
		t.Fail()
	}

	results, _ := config.GenerateResults()
	entries := results.ResultSets[0].Entries

	firstName, err := entries.GetEntryWithName("First Name")
	if err != nil {
		t.Fail()
	}

	lastName, _ := entries.GetEntryWithName("Last Name")
	fullName, _ := entries.GetEntryWithName("Full Name")

	if config.NewFullNameFromNames(firstName.Value, lastName.Value) != fullName.Value {
		t.Fail()
	}
}

func TestCanConvertToRows(t *testing.T) {
	config, _ := oarconfig.LoadConfigurationFromFile("./data/test.json")
	results, _ := config.GenerateResults()

	numentries := results.NumEntries()
	allEntries := make(map[*oarconfig.ResultSetEntry]bool)

	count := 0
	for _, set := range results.ResultSets {
		for _, entry := range set.Entries {
			allEntries[entry] = false

			count++
		}
	}

	if count != numentries {
		t.Fail()
	}

	rows := results.AsRows()
	for _, row := range rows {
		if row == nil {
			t.Fail()
		}

		for _, entry := range row {
			if entry == nil {
				t.Fail()
			}

			// if we've already seen this entry or the entry isn't in the set of all entries, fail
			seenEntry, ok := allEntries[entry]
			if !ok || seenEntry {
				t.Fail()
			}
		}
	}
}

func TestWillGenerateRandomValues(t *testing.T) {
	config, _ := oarconfig.LoadConfigurationFromFile("./data/test.json")

	results1, _ := config.GenerateResults()
	results2, _ := config.GenerateResults()

	if results1.ResultSets[0].Entries[0].Value == results2.ResultSets[0].Entries[0].Value {
		t.Fail()
	}
}
