package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)
type SaltVms struct {
    Up		[]string
    Down	[]string
}

func Scheduler(config *Config) {
	log.Info("Starting Scheduler ", config.Schedule)
	for {
		time.Sleep(time.Duration(config.Schedule) * time.Second)
		log.Debug("Scheduler next run in ", config.Schedule)

		cloudChan := make(chan *[]Vm)
		vms := make([]Vm, 0, 10000)
		go GetAwsVms(config, cloudChan)
		go GetOsVms(config, cloudChan)
		cloudChan <- &vms
		cloudVms := <-cloudChan

		log.Debug("Scheduler got cloud vms ",cloudVms)

		var minions SaltVms
		c2 := make(chan *SaltVms)
		go GetSaltVms(c2)
		c2 <- &minions
		saltVms := <-c2
		log.Debug("Scheduler got salt vms ",saltVms)

		for _, saltVm := range saltVms.Down {
			found := false
			for _, cloudVm := range *cloudVms {
				if saltVm == cloudVm.Name {
					found = true
				}
			}
			if !found {
				log.Debug("send instance can be removed from salt ",saltVm)
			}
		}
		for _, cloudVm := range *cloudVms {
			found := false
			for _, saltVm := range saltVms.Down {
				if saltVm == cloudVm.Name {
					found = true
					log.Debug("send instance is down ",cloudVm)
				}
			}
			for _, saltVm := range saltVms.Up {
				if saltVm == cloudVm.Name {
					found = true
				}
			}
			if !found {
				log.Debug("send instance not registred in salt ",cloudVm)
			}
		}
	}
}
