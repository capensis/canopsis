# État des composants de Canopsis

Voici les commandes afin de connaître l'état des composants de Canopsis.

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

## Gestion avancée des composants

[Arrêt et relance des composants](../guide-administration/gestion-composants/arret-relance-composants.md)
