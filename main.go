package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
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
  helper [--log-level=<level>] [--config=<path>] [--fastaccept] [--nocleanup] [--noscheduler] [--allow-known]
  helper -h | --help
  helper --version

Options:
  -h --help     		Show this screen.
  --version     		Show version.
  --log-level=<level>   Log level (debug, info, warn, error, fatal) [default: info]
  --config=<path>   	Path to config file 			[default: ./config.yaml].`

	arguments, _ := docopt.Parse(usage, nil, true, "Salt helper 0.1", false)

	log.SetOutput(os.Stderr)
	logLevel := arguments["--log-level"].(string)
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.SetLevel(level)
	}
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.Info("Starting salt-helper...")

	configPath := fmt.Sprint(arguments["--config"])
	config := GetConfig(configPath)

	done := make(chan bool)
	go Listener(config)
	if !arguments["--noscheduler"].(bool) {
		go Scheduler(config)
	}
	<-done
}
