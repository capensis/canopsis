# Healthcheck

Voici les commandes afin de connaître l'état des services autour de Canopsis.

## MongoDB

```bash
#connaître l'état du service
systemctl status mongod
```

## RabbitMQ

```bash
#connaître l'état du service
systemctl status rabbitmq-server
```

## InfluxDB

```bash
#connaître l'état du service
systemctl status influxdb
```

## Redis

```bash
#connaître l'état du service
systemctl status redis
```

## Canopsis

### Gestion de l'hyperviseur

```bash
#connaître l'état du service
/opt/canopsis/bin/canopsis-systemd status
```

## Gestion des services avancé

[Install, arrêt et relance](../gestion-services/installation-arret-relance.md)

## A venir

Healtheck va bientôt arriver sur votre interface Canopsis ! 