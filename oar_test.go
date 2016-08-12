package main

import (
	"testing"

	oarconfig "github.com/elauffenburger/oar/configuration"
)

func TestCanLoadFromFile(t *testing.T) {
	_, err := oarconfig.LoadConfigurationFromFile("./data/test.json")

	if err != nil {
		t.Fail()
	}
}
