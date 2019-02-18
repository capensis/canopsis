# Accès et identifiants Canopsis par défaut

Après une installation de Canopsis, ses différents composants utiliseront les adresses et identifiants par défaut suivants.

## Accès à Canopsis

### Interface web Canopsis

Par défaut, l'interface web de Canopsis est disponible sur : [http://localhost:8082/](http://localhost:8082/).

Si elle n'est pas activée par défaut, la nouvelle interface « UIv3 » peut être chargée explicitement à l'adresse [http://localhost:8082/en/static/canopsis-next/dist/index.html#](http://localhost:8082/en/static/canopsis-next/dist/index.html#).

Identifiants par défaut : `root` / `root`.

Parcourez ensuite [le guide d'utilisation](../../../guide-utilisation/) pour en apprendre davantage sur l'interface web de Canopsis.

## Accès aux composants internes de Canopsis

### Interface web RabbitMQ

Par défaut, l'interface web d'administration de RabbitMQ est disponible sur : [http://localhost:15672/](http://localhost:15672/).

Identifiants par défaut : `cpsrabbit` / `canopsis`.

### Bus AMQP RabbitMQ

Le bus AMQP RabbitMQ par défaut est : `amqp://cpsrabbit@canopsis:localhost:5672/canopsis`.

### MongoDB

En ligne de commande, la base de données MongoDB est accessible avec la commande `mongo localhost`.

Identifiants par défaut : `cpsmongo` / `canopsis`.

### InfluxDB

En ligne de commande, la base de métriques InfluxDB est accessible avec la commande `influx -username cpsinflux -password canopsis -database canopsis`.

Identifiants par défaut : `cpsinflux` / `canopsis`.
