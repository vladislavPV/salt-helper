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
	if err != nil {
		log.Error(err)
	}
	err = yaml.Unmarshal([]byte(out.Stdout), &minions)
	if err != nil {
		log.Error(err)
	}
	log.Debug("Salt minions loaded", minions)
	c <- minions
}
