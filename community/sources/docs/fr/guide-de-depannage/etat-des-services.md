# État des services Canopsis

Voici les commandes afin de connaître l'état des services autour de Canopsis.

## MongoDB

```sh
systemctl status mongod.service
```

## RabbitMQ

```sh
systemctl status rabbitmq-server.service
```

## Redis

```sh
systemctl status redis.service
```

## Nginx

```sh
systemctl status nginx.service
```

## Canopsis

### Gestion des moteurs et services internes

```sh
canoctl status
```

## Gestion avancée des services

[Arrêt et relance des services](../guide-administration/gestion-services/arret-relance-services.md)
