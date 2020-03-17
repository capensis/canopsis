# Canopsis SDK Usage

## Repositories dependencies

Some additional repositories needed to build the canopsis project will be download under `[canopsis-sources]/deps` by the `install.sh` script provided by the `sdk`.

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

#### Once these softwares are installed, you can use the Vagrantfile provided with canopsis sources under the `sdk`  directory

```shell
[canopsis-sources] $ cd sdk 
[canopsis-sources]/sdk $ vagrant up
```

If you want to update sdk use `provision` method

```sh
[canopsis-sources] $ cd sdk 
[canopsis-sources]/sdk $ vagrant provision
```

If you want to use some options to sdk, go inside the box and use

* `debug` mode

```shell
[canopsis-sources] $ cd sdk 
[canopsis-sources]/sdk $ vagrant ssh
vagrant@ubuntu-bionic:~ $ cd canopsis/sdk
vagrant@ubuntu-bionic:~/canopsis/sdk $ install.sh -d
```

* `CAT` installation

```shell
[canopsis-sources] $ cd sdk 
[canopsis-sources]/sdk $ vagrant ssh
vagrant@ubuntu-bionic:~ $ cd canopsis/sdk
vagrant@ubuntu-bionic:~/canopsis/sdk $ install.sh -c
```

:warning: Be sure to logout at the end of the install !

**Restart the VM if you use ssh ControlMaster sockets or if you're not sure to
don't use them**

```shell
canopsis-sources/sdk $ vagrant reload
```

Vagrant installation is sharing current canopsis sources with the directory `~/canopsis` inside the VagrantBox.

:rocket: So you can now work with you preferred IDE with the canopsis sources you previously cloned :up:


### `EXPERIMENTAL METHOS` - Local Installation inside a supported OS

#### System Requirements for local installation

* `For Ubuntu and Debian,be sure to use 64bit system`
* Ubuntu 18.04 LTS
* :warning: `curl` must be installed first

#### Local Installation

:warning: You must have `sudo` configured to install some packages and create some links

* Install Canopsis OpenCore SDK

```sh
[canopsis-sources] $ sdk/install.sh
```

* Using `debug` mode

```sh
[canopsis-sources] $ sdk/install.sh -d
```

* Install Canopsis CAT SDK

```sh
[canopsis-sources] $ sdk/install.sh -c
```


:warning: Be sure to reload your terminal at the end of the install !



## How to Launch Canopsis stack from SDK

:warning: If you are outside the Vagrant Box, don't forget to `vagrant ssh` first, before executing these commands

### Requirements

If your Canopsis's version < `3.31.0` ,
edit ```[canopsis-sources]/sources/canopsis/etc/common/redis_store.conf``` 

Fist Update host var :

```
host = localhost
```

### Usage

* For compiling Canopsis

```sh
$ canopsis-sdk --build_all | --build_go | --build_python | --build_ui | --cat
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
