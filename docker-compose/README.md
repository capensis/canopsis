# Quickstart

## Requirements

To use Canopsis, you first have to install :

* Git
* Docker Engine ( >= 19.x )
* docker-compose ( >= 1.2x )



## Get this repository

```shell
$ git clone https://gitlab.canopsis.net/canopsis/canopsis.git
$ cd canopsis
```



## Update with latest available Images

```shell
$ docker-compose pull
```



## Start Canopsis Stack

```shell
$ docker-compose up -d
```



## Access to Web Interface

Browse url http://localhost:8082 or http://localhost (:warning: be sure that no other service is running on the `TCP/80` port to avoid port conflict )

With credentials: 

* user: `root`
* password: `root`

## Stop Canopsis Stack

```shell
$ docker-compose down -v
```

