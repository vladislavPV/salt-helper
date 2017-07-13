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

	return &config

}
