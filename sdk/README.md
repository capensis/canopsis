# Canopsis SDK Usage

## Requirements for CAT projet

* :warning: You must have 
  * access to gitlab `git.canopsis.net` with a created user and you must have `Developer` role into some projects into namespace(s) :
    * https://git.canopsis.net/cat
  * an SSH key configured for your user inside gitlab `git.canopsis.net` and into your machine ( SSH agent, ssh_config )


### `RECOMMENDED METHOD` - Installation inside a Vagrant box ( Virtualbox provider )

#### First you have to install on your current host

* `VirtualBox` ( recommended version == `6.0.x` )

https://www.virtualbox.org/wiki/Downloads

* And `Vagrant` ( recommended version >= `2.2.x` )

https://www.vagrantup.com/downloads.html

#### Once these softwares are installed, you can retrieve Vagrantfile from Canopsis SDK

```shell
$ mkdir -p ~/vagrant/canopsis-sdk
$ cd ~/vagrant/canopsis-sdk
$ curl https://gitlab.capensis.fr/capensis-opensource/sdk-canopsis-core/raw/go-system/vagrant/Vagrantfile -o Vagrantfile
```

An start Vagrant box

```sh
$ cd ~/vagrant/canopsis-sdk
$ vagrant up
```

If you want to update sdk use `provision` method

```sh
$ cd ~/vagrant/canopsis-sdk
$ vagrant provision
``

If you want to use some options to sdk, go inside the box and use

* `debug` mode

```shell
$ cd ~/vagrant/canopsis-sdk
$ vagrant ssh
$ vagrant@ubuntu-bionic:~$ /bin/sh ./install.sh -d
```
* `CAT` installation

```shell
$ cd ~/vagrant/canopsis-sdk
$ vagrant ssh
$ vagrant@ubuntu-bionic:~$ /bin/sh ./install.sh -c
```

:warning: Be sure to logout at the end of the install !

**Restart the VM if you use ssh ControlMaster sockets or if you're not sure to
don't use them**
```
$ vagrant reload
```

Vagrant installation is sharing local vagrant box directory `/opt/canopsis-sdk` whith a symlink to `~/canopsis-sdk`, between local machine and vagrant box.

:rocket: So you can now work with you preferred IDE with this sources location ( `~/canopsis-sdk` ) :up:


### `EXPERIMENTAL METHOS` - Local Installation inside a supported OS

#### System Requirements for local installation

* `For Ubuntu and Debian,be sure to use 64bit system`
* Ubuntu 18.04 LTS
* :warning: `curl` must be installed first

#### Local Installation

:warning: You must have `sudo` configured to install some packages and create some links

* Get Canopsis SDK Install Script

Using wget

```sh
$ cd /tmp && wget https://gitlab.capensis.fr/capensis-opensource/sdk-canopsis-core/raw/go-system/install.sh
```

Or Curl

```sh
$ cd /tmp && curl https://gitlab.capensis.fr/capensis-opensource/sdk-canopsis-core/raw/go-system/install.sh -o install.sh
```

* Install Canopsis OpenCore SDK

```sh
$ /bin/sh /tmp/install.sh
```

* Using `debug` mode

```sh
$ /bin/sh /tmp/install.sh -d
```

* Install Canopsis CAT SDK

```sh
$ sudo /bin/sh /tmp/install.sh -c
```


:warning: Be sure to reload your terminal at the end of the install !



## How to Launch Canopsis stack from SDK

:warning: If you are within the Vagrant Box, don't forget to `vagrant ssh` first, before executing these commands

### Requirements

If your Canopsis's version < `3.31.0` ,
In ```/opt/canopsis-sdk/src/canopsis/sources/canopsis/etc/common/redis_store.conf``` 

Fist Update host var :

```
host = localhost
```

### Usage

* For compiling Canopsis

```sh
$ canopsis-sdk --build_all | --build_go | --build_python | --build_ui
```

* To start Canopsis

default `Go` stack

```sh
$ canopsis-sdk --start_go
```

or old `Python` stack

```sh
$ canopsis-sdk --start_python
```

:information_source: you can also build and start at the same time

```sh
$ canopsis-sdk --build_all --start_go
```

* Testing your installation

```sh
curl -u root:root http://localhost:8082/api/v2/event -H 'Content-Type: application/json' -d ' {
    "author": "root",
    "component": "component test",
    "connector": "connector test",
    "connector_name": "connector_name",
    "event_type": "check",
    "output": "",
    "resource": "resource",
    "source_type": "resource",
    "state": 3
}'
```

### conflict between engine 

| Engine      | Go stack   | Python stack |
| --- | --- | --- |
| canopsis-webserver | yes | yes |
| canopsis-engine@dynamic-alerts | no | yes |
| canopsis-engine@cleaner-cleaner_events | no | yes |
| canopsis-engine@dynamic-context-graph | no | yes |
| canopsis-engine@event_filter-event_filter | no | yes |
| canopsis-engine@dynamic-pbehavior | yes | yes |
|-canopsis-engine@dynamic-watcher | no | yes |
| canopsis-engine@event_filter-event_filter | no | yes |
| canopsis-engine@metric-metric | no | yes |
| canopsis-engine@scheduler-scheduler | no | yes |
| canopsis-engine@task_importctx-task_importctx | no | yes |
| canopsis-engine-go@action | yes | no |
| canopsis-engine-go@axe | yes | no |
| canopsis-engine-go@che | yes | no |
| canopsis-engine-go@heartbeat | yes | no |
| canopsis-engine-go@watcher | yes | no |

### Systemd Management

#### Managing specific engine

```sh 
$ systemctl --user start/restart/stop/status  canopsis-engine-go@<engine-name>.service
or
$ systemctl --user start/restart/stop/status  canopsis-engine@<engine-name>.service
```

#### To start Canopsis stack from systemctl

```sh
$ systemctl --user start/restart/stop/status canopsis-{go|python}-systemd.service
```

#### To stop all Canopsis running services

```sh
$ systemctl --user start canopsis-down.service
```

#### During the installation, SDK will create the following system units

* canopsis-go-systemd.service
* canopsis-python-systemd.service 
* canopsis-down.service 
