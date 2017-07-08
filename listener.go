package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
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

	log.Printf("Listening dir: %#v\n", config.Watchdir)

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
					log.Println("Created file:", filename)

					c := make(chan *[]Vm)
					vms := make([]Vm, 0, 10000)

					go GetAwsVms(config, c)
					go GetOsVms(config, c)

					c <- &vms
					result := <-c

					found := false
					for _, vm := range *result {
						if vm.Name == filename {
							found = true
							err := os.Rename(fullpath, config.DstDir+"/"+filename)
							if err != nil {
								log.Fatal(err)
							} else {
								log.Println("Adding vm to salt", vm)
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
							log.Println("Rejecting vm from salt", vm)
							SendToSlack(config, vm)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Printf("error:", err)
			}
		}
	}()

	err = watcher.Add(config.Watchdir)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
