package main

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
)

func GetOsVms(config *Config, c chan *[]Vm) {
	log.Debug("Init os")
	vms := <-c

	for _, account := range config.OsAccounts {
		for _, region := range account.Regions {
			log.Debug("Get vms from ", account.Name, region)

			authOpts := gophercloud.AuthOptions{
				IdentityEndpoint: account.Auth_url,
				Username:         account.Username,
				Password:         account.Password,
				TenantName:       account.Project_name,
			}

			provider, err := openstack.AuthenticatedClient(authOpts)
			if err != nil {
				log.Error(err)
			}

			client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
				Region: region,
			})
			if err != nil {
				log.Error(err)
			}

			opts := servers.ListOpts{}
			pager := servers.List(client, opts)

			err = pager.EachPage(func(page pagination.Page) (bool, error) {
				serverList, _ := servers.ExtractServers(page)

				for _, server := range serverList {
					*vms = append(*vms, Vm{
						server.Name,
						region,
						account.Name,
						server.ID,
					})
				}
				return true, nil
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
	c <- vms
}
