package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// type VmSalt struct {
// 	Name    string
// 	State   string
// }

func Scheduler(config *Config) {
	log.Info("Starting Scheduler ", config.Schedule)
	for {
		time.Sleep(time.Duration(config.Schedule) * time.Millisecond)
		log.Debug("Scheduler next run in ", config.Schedule)

		// get aws vm
		// get os vm
		// get salt vm
		// c := make(chan *[]Vm)
		// vms := make([]Vm, 0, 10000)

		// go GetAwsVms(config, c)
		// go GetOsVms(config, c)
		// c <- &vms
		// result := <-c

		// log.Info(result)
		// c2 := make(chan *[]VmSalt)
		// vms := make([]VmSalt, 0, 10000)
		// go GetSaltVms(config, c2)
		// compare
	}
}
