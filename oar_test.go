package main

import (
	"fmt"
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

	results, _ := config.GenerateResults()
	entries := results.ResultSets[0].Entries

	firstName, err := entries.GetEntryWithName("FirstName")
	if err != nil {
		t.Fail()
	}

	lastName, _ := entries.GetEntryWithName("LastName")
	fullName, _ := entries.GetEntryWithName("FullName")

	if fmt.Sprintf("%s %s", firstName.Value, lastName.Value) != fullName.Value {
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

	rows := results.ToRows()
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

	n, n2 := len(results1.ResultSets[0].Entries)-1, len(results2.ResultSets[0].Entries)-1
	if n != n2 {
		t.Fail()
	}

	if results1.ResultSets[0].Entries[n].Value == results2.ResultSets[0].Entries[n].Value {
		t.Fail()
	}

	if results1.ResultSets[0].Entries[n].Value == results1.ResultSets[1].Entries[n].Value {
		t.Fail()
	}
}

func TestCanConvertToJsonArray(t *testing.T) {
	config, _ := oarconfig.LoadConfigurationFromFile("./data/test.json")

	results, _ := config.GenerateResults()
	jsonarray := results.ToJsonArray()

	keys := jsonarray[0].Keys()
	key := keys[len(keys)-1]

	val1, _ := (*jsonarray[0])[key]
	val2, _ := (*jsonarray[1])[key]

	if val1 == val2 {
		t.Fail()
	}
}
