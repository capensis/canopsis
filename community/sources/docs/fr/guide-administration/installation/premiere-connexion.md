# Connexion à Canopsis et à ses composants

Après une installation de Canopsis, ses différents composants utiliseront les adresses et identifiants par défaut suivants.

## Accès à l'interface web de Canopsis

Voir [le Guide de premier accès](../../guide-utilisation/premier-acces/index.md).

## Accès aux composants internes de Canopsis

De manière générale, les URI de connexion aux différents composants de Canopsis se trouve dans le fichier :

* `/opt/canopsis/etc/go-engines-vars.conf` pour une [installation paquets](installation-paquets.md)
* `compose.env` pour une [installation avec Docker Compose](installation-conteneurs.md)

### Interface web RabbitMQ

Par défaut, l'interface web d'administration de RabbitMQ est disponible depuis votre navigateur à l'adresse suivante : <http://localhost:15672/>.

Identifiants par défaut : `cpsrabbit` / `canopsis`.

### Bus AMQP RabbitMQ

Le bus AMQP RabbitMQ par défaut est : `amqp://cpsrabbit@canopsis:localhost:5672/canopsis`.

### MongoDB

En ligne de commande, la base de données MongoDB est accessible avec la commande  
`mongo mongodb://cpsmongo:canopsis@mongodb/cnopsis?replicaSet=rs0`

Identifiants par défaut : `cpsmongo` / `canopsis`.

### PostgreSQL

La base de données PostgreSQL est accessible avec la commande `psql -U cpspostgres -W -d canopsis -h localhost`.

Identifiants par défaut : `cpspostgres` / `canopsis`.
