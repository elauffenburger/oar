package main

import (
	"flag"
	"fmt"

	"encoding/json"

	"os"

	oarconfig "github.com/elauffenburger/oar/configuration"
)

var configFlag = flag.String("config", "", "json file to load configuration from")

func main() {
	flag.Parse()

	configFlagValue := *configFlag

	if len(configFlagValue) == 0 {
		panic("No configuration file provided!")
	}

	config, err := oarconfig.LoadConfigurationFromFile(configFlagValue)
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration file: %s", err))
	}

	results, err := config.GenerateResults()
	if err != nil {
		panic(fmt.Sprintf("Error generating results: '%s'", err))
	}

	rows := results.AsRows()
	marshalledbytes, err := json.Marshal(&rows)
	if err != nil {
		panic(fmt.Sprintf("Error marshalling results: '%s", err))
	}

	// print results as json to stdout
	fmt.Fprint(os.Stdout, string(marshalledbytes))
}
