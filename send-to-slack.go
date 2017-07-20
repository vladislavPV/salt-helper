package main

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	log "github.com/sirupsen/logrus"
)

func SendToSlack(config *Config, vm Vm) {
	webhookUrl := config.Slackwebhook

	attachment := slack.Attachment{}

	attachment.Color = &vm.Color
	attachment.AddField(slack.Field{Title: "Name", Value: vm.Name, Short: true})
	attachment.AddField(slack.Field{Title: "Region/Account", Value: vm.Region + "/" + vm.Account, Short: true})
	attachment.AddField(slack.Field{Title: "ID", Value: vm.Id, Short: true})
	attachment.AddField(slack.Field{Title: "Status", Value: vm.Status, Short: true})

	payload := slack.Payload{
		Text:        "",
		Username:    config.Slackbotname,
		Channel:     config.Slackchannel,
		IconEmoji:   config.Slackimoji,
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		log.Error(err)
	}
}
