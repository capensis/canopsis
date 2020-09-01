# Healthcheck

Voici les commandes afin de connaître l'état des services autour de Canopsis.

## MongoDB

```sh
systemctl status mongod.service
```

## RabbitMQ

```sh
systemctl status rabbitmq-server.service
```

## InfluxDB

```sh
systemctl status influxdb.service
```

## Redis

```sh
systemctl status redis.service
```

## Canopsis

### Gestion de l'hyperviseur

```sh
canoctl status
```

## Gestion avancées des services

[Arrêt et relance des services](../gestion-services/arret-relance-services.md)

## À venir

Healtheck va bientôt arriver sur votre interface Canopsis ! 
