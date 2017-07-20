package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func GetSaltVms(c chan *SaltVms) {
	log.Debug("Starting salt get minions")
	minions := <-c

	cmd := "salt-run"

	args := []string{
		"manage.status",
		"--out",
		"yaml",
	}
	out, err := ExecuteCommand(cmd, args)

	err = yaml.Unmarshal([]byte(out.Stdout), &minions)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Out loaded", minions)
	c <- minions
}
