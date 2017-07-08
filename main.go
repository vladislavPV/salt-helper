package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

type Vm struct {
	Name    string
	Region  string
	Account string
	Id      string
}

func main() {
	usage := `Salt helper.

Usage:
  helper [--config=<path>] [--fastaccept] [--nocleanup]
  helper -h | --help
  helper --version

Options:
  -h --help     	Show this screen.
  --version     	Show version.
  --config=<path>   Path to config file 			[default: ./config.yaml].`

	arguments, _ := docopt.Parse(usage, nil, true, "Salt helper 0.1", false)

	configPath := fmt.Sprint(arguments["--config"])
	config := GetConfig(configPath)

	done := make(chan bool)
	go Listener(config)
	// Scheduler(config)
	<-done
}
