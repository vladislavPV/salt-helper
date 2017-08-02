# salt-helper. [![Build Status](https://travis-ci.org/vladislavPV/salt-helper.svg?branch=master)](https://travis-ci.org/vladislavPV/salt-helper)
### Auto accept salt minions with checks

`salt-helper` is a Go application which runs on the same instance with salt-master and auto accepts minions
Before accepting, helper checks if hostname(minion key filename) exists in configured AWS or
OpenStack regions , and sends alert in slack.
This allows you to have dynamic infra in multiple clouds. Also helper can keep your salt clean of old removed minions with scheduled checks.

### Install
Download binary for linux from https://github.com/vladislavPV/salt-helper/releases
You can instantiate it by using supervisord, systemd, upstart or any other init system.

### Config
By default config file(config.yaml) should be in the same dir as salt-helper-linux
or you can use --config option. Example config is here:
https://github.com/vladislavPV/salt-helper/blob/master/config-example.yaml

Also few other options available
	--log-level=debug|info	allows you to set verbosity
	--fastaccept			do cloud check after accepting minion. Could be usefull for autoscaled instances
	--nocleanup     		disable cleanup of dead salt minions
	--noscheduler   		disable scheduled checks. will not send you if minion is down
	--allow-known   		force accept minions already existing in salt. Could be usefull for autoscaled instances


NOTE!	In Aws your ApiKey should be able to read ec2 metadata, so in IAM you have to allow DescribeInstances.

### Run
```
$ sudo ./salt-helper-linux --config /path/to/config.yaml
```

### Behavior

When new minion is trying to connect, salt master creates file with minion_id in /etc/salt/pki/master/minions_pre/
Salt-helper listens on such events and checks the minion_id in all known accounts/regions.
If minion_id found in clouds salt-helper will move file in /etc/salt/pki/master/minions/
minion_id == instance Name

