# Connexion à Canopsis et à ses composants

Après une installation de Canopsis, ses différents composants utiliseront les adresses et identifiants par défaut suivants.

## Accès à l'interface web de Canopsis

Par défaut, l'interface web de Canopsis est disponible depuis votre navigateur à l'adresse suivante : <http://localhost>.

Les identifiants de connexion par défaut sont `root` / `root`.

Parcourez ensuite [le Guide d'utilisation](../../guide-utilisation/index.md) pour en apprendre davantage sur l'interface web de Canopsis.

## Accès aux composants internes de Canopsis

### Interface web RabbitMQ

Par défaut, l'interface web d'administration de RabbitMQ est disponible depuis votre navigateur à l'adresse suivante : <http://localhost:15672/>.

Identifiants par défaut : `cpsrabbit` / `canopsis`.

### Bus AMQP RabbitMQ

Le bus AMQP RabbitMQ par défaut est : `amqp://cpsrabbit@canopsis:localhost:5672/canopsis`.

### MongoDB

En ligne de commande, la base de données MongoDB est accessible avec la commande `mongo -u cpsmongo -p canopsis canopsis`.

Identifiants par défaut : `cpsmongo` / `canopsis`.
