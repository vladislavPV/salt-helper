package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

func GetAwsVms(config *Config, c chan *[]Vm) {
	log.Println("Init aws")
	vms := <-c

	for _, account := range config.AwsAccounts {
		for _, region := range account.Regions {
			log.Println("Get vms from", account.Name, region)

			sess := session.New()
			sess.Config.Credentials = credentials.NewStaticCredentials(account.Id, account.Secret, "")
			sess.Config.Region = &region

			ec2svc := ec2.New(sess)
			params := &ec2.DescribeInstancesInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("instance-state-name"),
						Values: []*string{aws.String("running"), aws.String("pending")},
					},
				},
			}
			resp, err := ec2svc.DescribeInstances(params)
			if err != nil {
				log.Fatal(err.Error())
			}

			for _, reservation := range resp.Reservations {
				for _, instance := range reservation.Instances {
					for _, tag := range instance.Tags {
						if *tag.Key == "Name" {
							*vms = append(*vms, Vm{
								*tag.Value,
								region,
								account.Name,
								*instance.InstanceId,
							})
						}
					}
				}
			}
		}
	}
	c <- vms
}
