package results

import (
	"fmt"

	conf "github.com/elauffenburger/oar/core/configuration"
)

type ResultsRowList []*ResultsRow
type Results struct {
	ResultSets ResultsRowList
}

type ResultRowValueList []*ResultsRowValue
type ResultsRow struct {
	Values ResultRowValueList
}

type ResultsRowValue struct {
	conf.ConfigurationField
	Value string
}

func (entries *ResultRowValueList) GetEntryWithName(name string) (*ResultsRowValue, error) {
	for _, entry := range *entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("Failed to find an entry with name '%s'", name)
}

func (entries *ResultRowValueList) GetEntryWithType(entryType string) (*ResultsRowValue, error) {
	for _, entry := range *entries {
		if entry.Type == entryType {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("Failed to find an entry with type '%s'", entryType)
}

func NewRowsOfSize(n int) []ResultRowValueList {
	return make([]ResultRowValueList, n)
}

func NewRows() []ResultRowValueList {
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

	return len(results.ResultSets[0].Values) * numrows
}

func (results *Results) ToRows() []ResultRowValueList {
	rows := NewRowsOfSize(results.NumRows())

	for i, set := range results.ResultSets {
		rows[i] = set.Values
	}

	return rows
}
