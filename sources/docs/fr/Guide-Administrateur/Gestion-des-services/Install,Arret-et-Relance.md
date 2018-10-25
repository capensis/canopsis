# Guide Administrateur

## Section : Gestion des services / Install, arrêt et relance

### MongoDB

[Installer MongDB](https://docs.mongodb.com/manual/installation/#tutorials) 

```bash
#démarrer le service
systemctl start mongod
#stopper le service
systemctl stop mongod
#redémarrer le service
systemctl restart mongod
```

### RabbitMQ

[Installer RabbitMQ](https://www.rabbitmq.com/download.html)

```bash
#démarrer le service
systemctl start rabbitmq-server
#stopper le service
systemctl stop rabbitmq-server
#redémarrer le service
systemctl restart rabbitmq-server
```

## InfluxDB

[Installer InfluxDB 0.10.x](https://docs.influxdata.com/influxdb/v1.6/introduction/installation/)

```bash
#démarrer le service
systemctl start influxdb
#stopper le service
systemctl stop influxdb
#redémarrer le service
systemctl restart influxdb
```

### Redis

[Installer Redis](https://redis.io/topics/quickstart) 

```bash
#démarrer le service
systemctl start redis
#stopper le service
systemctl stop redis
#redémarrer le service
systemctl restart redis
```

## Canopsis

### Gestion de l'hyperviseur

```bash
#démarrer le service
/opt/canopsis/bin/canopsis-systemd start
#stopper le service
/opt/canopsis/bin/canopsis-systemd stop
#redémarrer le service
/opt/canopsis/bin/canopsis-systemd restart
```

## Aller plus loin 

Pour connaître l'état de votre service, [rendez-vous ici](/doc-ce/Guide Administrateur/Troubleshooting/HealthCheck.md)
