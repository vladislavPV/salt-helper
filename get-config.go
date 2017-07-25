package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Watchdir     string
	DstDir       string
	RejectedDir  string
	Slackwebhook string
	Slackchannel string
	Slackimoji   string
	Slackbotname string
	Schedule     int
	AwsAccounts  []AwsAccounts
	OsAccounts   []OsAccounts
}
type AwsAccounts struct {
	Name    string
	Regions []string
	Id      string
	Secret  string
}
type OsAccounts struct {
	Name         string
	Regions      []string
	Username     string
	Password     string
	Version      string
	Auth_url     string
	Project_name string
}

func GetConfig(filename string) *Config {
	log.Info("Loading config ", filename)
	fullpath, _ := filepath.Abs(filename)
	yamlFile, err := ioutil.ReadFile(fullpath)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Config loaded")

	// some sane defaults
	if len(config.Watchdir) == 0{
		config.Watchdir = "/etc/salt/pki/master/minions_pre/"
	}
	if len(config.DstDir) == 0 {
		config.DstDir = "/etc/salt/pki/master/minions/"
	}
	if len(config.RejectedDir) == 0 {
		config.RejectedDir = "/etc/salt/pki/master/minions_rejected/"
	}
	if len(config.Slackchannel) == 0 {
		config.Slackchannel = "#general"
	}
	if len(config.Slackimoji) == 0 {
		config.Slackimoji = ":b:"
	}
	if len(config.Slackbotname) == 0 {
		config.Slackbotname = "Salt Bot"
	}
	if config.Schedule == 0 {
		config.Schedule = 300
	}

	log.Debug("Config: ", config)
	return &config
}
