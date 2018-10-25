# Première connexion, identifiants par défaut

## Interface web Canopsis

Par défaut, l'interface web de Canopsis se lance sur : http://localhost:8082/

Identifiants par défaut : `root` / `root`.

## Interface web RabbitMQ

Par défaut, l'interface web d'administration de RabbitMQ se lance sur : http://localhost:15672/

Identifiants par défaut : `cpsrabbit` / `canopsis`.

## Bus AMQP RabbitMQ

Le bus AMQP RabbitMQ par défaut est : `amqp://cpsrabbit@canopsis:localhost:5672/canopsis`.

## MongoDB

En ligne de commande, la base de données MongoDB est accessible avec la commande `mongo localhost`.

Identifiants par défaut : `cpsmongo` / `canopsis`.

**TODO (DWU) :** lien vers manip Robo3T ?

## InfluxDB

En ligne de commande, la base de métriques InfluxDB est accessible avec la commande `influx -username cpsinflux -password canopsis -database canopsis`.

Identifiants par défaut : `cpsinflux` / `canopsis`.