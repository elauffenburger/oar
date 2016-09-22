package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/elauffenburger/oar/core"
)

var configFlag = flag.String("config", "", "json file to load configuration from")
var rowsFlag = flag.Int("rows", 0, "number of rows to generate")

func main() {
	flag.Parse()

	configFlagValue := *configFlag
	rows := *rowsFlag

	if len(configFlagValue) == 0 {
		panic("No configuration file provided!")
	}

	config, err := core.LoadConfigurationFromFile(configFlagValue)
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration file: %s", err))
	}

	if rows != 0 {
		config.NumRows = rows
	}

	results, err := config.GenerateResults()
	if err != nil {
		panic(fmt.Sprintf("Error generating results: '%s'", err))
	}

	// print results to stdout
	formatter := core.GetOutputFormatter(config)
	fmt.Fprint(os.Stdout, formatter.Format(results))
}
