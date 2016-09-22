package main

import (
	"fmt"
	"testing"

	"github.com/elauffenburger/oar/core"
	"github.com/elauffenburger/oar/core/output"
	res "github.com/elauffenburger/oar/core/results"
)

func TestCanLoadFromFile(t *testing.T) {
	config, err := core.LoadConfigurationFromFile("./test/test.json")

	if err != nil {
		t.Fail()
	}

	results, err := core.GenerateResults(config)

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
	config, _ := core.LoadConfigurationFromFile("./test/test.json")

	results, _ := core.GenerateResults(config)
	entries := results.ResultSets[0].Values

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
	config, _ := core.LoadConfigurationFromFile("./test/test.json")
	results, _ := core.GenerateResults(config)

	numentries := results.NumEntries()
	allEntries := make(map[*res.ResultsRowValue]bool)

	count := 0
	for _, set := range results.ResultSets {
		for _, entry := range set.Values {
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
	config, _ := core.LoadConfigurationFromFile("./test/test.json")

	results1, _ := core.GenerateResults(config)
	results2, _ := core.GenerateResults(config)

	n, n2 := len(results1.ResultSets[0].Values)-1, len(results2.ResultSets[0].Values)-1
	if n != n2 {
		t.Fail()
	}

	if results1.ResultSets[0].Values[n].Value == results2.ResultSets[0].Values[n].Value {
		t.Fail()
	}

	if results1.ResultSets[0].Values[n].Value == results1.ResultSets[1].Values[n].Value {
		t.Fail()
	}
}

func TestCanConvertToJsonArray(t *testing.T) {
	config, _ := core.LoadConfigurationFromFile("./test/test.json")

	results, _ := core.GenerateResults(config)

	jsonformatter := &output.JsonOutputFormatter{}
	jsonarray := jsonformatter.ToJsonArray(results)

	keys := jsonarray[0].Keys()
	key := keys[len(keys)-1]

	val1, _ := (*jsonarray[0])[key]
	val2, _ := (*jsonarray[1])[key]

	if val1 == val2 {
		t.Fail()
	}
}
