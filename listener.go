package main

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

func Listener(config *Config) {

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

					if _, err := os.Stat(config.DstDir + "/" + filename); !os.IsNotExist(err) {
						err := os.Rename(fullpath, config.DstDir+"/"+filename)
						if err != nil {
							log.Fatal(err)
						} else {
							log.Info("Vm ", filename, " already exists, new key removed")
							vm := Vm{Name: filename, Region: "None", Account: "None", Id: "None"}
							log.Info("Rejecting vm from salt ", vm)
							SendToSlack(config, vm)
						}
						break
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
							err := os.Rename(fullpath, config.DstDir+"/"+filename)
							if err != nil {
								log.Fatal(err)
							} else {
								log.Info("Adding vm to salt ", vm)
								SendToSlack(config, vm)
							}
						}
					}

					if !found {
						err := os.Rename(fullpath, config.RejectedDir+"/"+filename)
						if err != nil {
							log.Fatal(err)
						} else {
							vm := Vm{Name: filename, Region: "None", Account: "None", Id: "None"}
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
