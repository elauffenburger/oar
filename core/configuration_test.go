package core

import (
	"testing"
)

func TestLoadJsonConfiguration(t *testing.T) {
	config, err := LoadConfigurationFromJson("{\"fields\":[{\"name\":\"test\", \"value\":\"12345\", \"type\": \"number\"}], \"numRows\":50}")

	if err != nil {
		t.Errorf("Error loading configuration: %s", err)
		t.Fail()
	}

	if config.NumRows != 50 {
		t.Errorf("Number of rows didn't match")
		t.Fail()
	}

	if len(config.Fields) != 1 {
		t.Errorf("Number of fields didn't match")
		t.Fail()
	}

	field := config.Fields[0]
	if field.Name != "test" || field.Type != "number" {
		t.Errorf("Something else weird happened; field: %v", field)
		t.Fail()
	}
}
