# Pré-requis SELinux

Ce document présente les pré-requis nécessaires au bon fonctionnement de Canopsis dans un environnement confiné par SELinux


## Ports en écoute pour Canopsis

| Port/protocole | role                            | Version minimale |
|----------------|---------------------------------|------------------|
| 8082/tcp       | webserver                       | toutes           |
| 4444/udp       | influxdb - écriture             | toutes           |
| 8083/tcp       | influxdb - réplication          | toutes           |
| 8086/tcp       | influxdb - lecture (apis rest)  | toutes           |
| 27017/tcp      | mongodb                         | toutes           |
| 15672/tcp      | rabbitmq - WebUI administration | toutes           |
| 5672/tcp       | rabbitmq - messages             | toutes           |
| 15671/tcp      | rabbitmq - ?                    | toutes           |
| 25672/tcp      | rabbitmq - ?                    | toutes           |
| 5671/tcp       | rabbitmq - ?                    | toutes           |
| 6379/tcp       | redis                           | 2.7.0            |


## Ports en écoute - Connecteurs

| Port/protocole | role                            |
| 162/udp        | snmp                            |


## Accès filesystem

| Dossier      | Droits           |
|--------------|------------------|
|/opt/canopsis | Lecture/écriture |
