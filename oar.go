package main

import (
    "flag"
    "fmt"

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

    // print results as json to stdout
    fmt.Fprint(os.Stdout, results.ToJson())
}
