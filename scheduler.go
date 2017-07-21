package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type SaltVms struct {
	Up   []string
	Down []string
}

func Scheduler(config *Config, nocleanup bool) {
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

		log.Debug("Scheduler got cloud vms ", cloudVms)

		var minions SaltVms
		c2 := make(chan *SaltVms)
		go GetSaltVms(c2)
		c2 <- &minions
		saltVms := <-c2
		log.Debug("Scheduler got salt vms ", saltVms)

		for _, saltVm := range saltVms.Down {
			found := false
			for _, cloudVm := range *cloudVms {
				if saltVm == cloudVm.Name {
					found = true
				}
			}
			if !found {
				if nocleanup {
					log.Debug("Instance, ", saltVm, " can be removed manualy from salt")
				} else {
					err := os.Remove(config.DstDir + "/" + saltVm)
					if err != nil {
						log.Fatal(err)
					} else {
						log.Debug("Instance ", saltVm, " was removed from salt")
						vm := Vm{Name: saltVm, Region: "None", Account: "None", Id: "None", Status: "Vm not found in cloud", Color: "warning"}
						SendToSlack(config, vm)
					}
				}
			}
		}
		for _, cloudVm := range *cloudVms {
			found := false
			for _, saltVm := range saltVms.Down {
				if saltVm == cloudVm.Name {
					found = true
					log.Debug("Instance is down ", cloudVm)
					vm := Vm{Name: cloudVm.Name, Region: cloudVm.Region, Account: cloudVm.Account, Id: cloudVm.Id, Status: "Is down", Color: "warning"}
					SendToSlack(config, vm)
				}
			}
			for _, saltVm := range saltVms.Up {
				if saltVm == cloudVm.Name {
					found = true
				}
			}
			if !found {
				log.Debug("Instance not registred in salt ", cloudVm)
				vm := Vm{Name: cloudVm.Name, Region: cloudVm.Region, Account: cloudVm.Account, Id: cloudVm.Id, Status: "Not registered in salt", Color: "warning"}
				SendToSlack(config, vm)
			}
		}
	}
}
