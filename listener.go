package main

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

func Listener(config *Config, fastaccept bool, allowknown bool) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	log.Info("Listening dir ", config.Watchdir)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					fullpath := event.Name
					filename := filepath.Base(fullpath)
					log.Info("Created file ", filename)

					if !allowknown {
						if _, err := os.Stat(config.DstDir + "/" + filename); !os.IsNotExist(err) {
							err := os.Remove(fullpath)
							if err != nil {
								log.Fatal(err)
							} else {
								vm := Vm{Name: filename, Region: "None", Account: "None", Id: "None", Status: "Already exist", Color: "danger"}
								SendToSlack(config, vm)
								log.Info("Vm ", filename, " already exists, new key removed")
							}
							break
						}
					}
					if fastaccept {
							err := os.Rename(fullpath, config.DstDir+"/"+filename)
							if err != nil {
								log.Fatal(err)
							} else {
								log.Info("Adding vm to salt ", filename)
							}
					}
					found := false
					c := make(chan *[]Vm)
					vms := make([]Vm, 0, 10000)

					go GetAwsVms(config, c)
					go GetOsVms(config, c)

					c <- &vms
					result := <-c

					for _, vm := range *result {
						if vm.Name == filename {
							found = true
							if !fastaccept {
								err := os.Rename(fullpath, config.DstDir+"/"+filename)
								if err != nil {
									log.Fatal(err)
								} else {
									log.Info("Adding vm to salt ", vm)
								}
							}
							vm := Vm{Name: filename, Region: vm.Region, Account: vm.Account, Id: vm.Id, Status: "Added to salt", Color: "good"}
							SendToSlack(config, vm)
							log.Info("Check is done for ", filename)
						}
					}

					if !found {
						if fastaccept {
							fullpath = config.DstDir+"/"+filename
						}
						err := os.Rename(fullpath, config.RejectedDir+"/"+filename)
						if err != nil {
							log.Fatal(err)
						} else {
							vm := Vm{Name: filename, Region: "None", Account: "None", Id: "None", Status: "Not found in clouds", Color: "danger"}
							log.Info("Rejecting vm from salt ", vm)
							SendToSlack(config, vm)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(config.Watchdir)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
