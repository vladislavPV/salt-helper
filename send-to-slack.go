package main

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"log"
)

func SendToSlack(config *Config, vm Vm) {
	webhookUrl := config.Slackwebhook

	attachment := slack.Attachment{}

	color := "good"
	msg := "Accepted in salt"
	if vm.Account == "None" {
		color = "danger"
		msg = "Rejected in salt"
	}

	attachment.Color = &color
	attachment.AddField(slack.Field{Title: "Name", Value: vm.Name, Short: true})
	attachment.AddField(slack.Field{Title: "Region/Account", Value: vm.Region + "/" + vm.Account, Short: true})
	attachment.AddField(slack.Field{Title: "ID", Value: vm.Id, Short: true})
	attachment.AddField(slack.Field{Title: "Status", Value: msg, Short: true})

	payload := slack.Payload{
		Text:        "",
		Username:    "Salt bot",
		Channel:     "@myname",
		IconEmoji:   ":monkey:",
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		log.Printf("error: %s\n", err)
	}
}
