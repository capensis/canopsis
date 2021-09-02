# Moteur `engine-service` (Community)

Les moteur `engine-service` permet de surveiller et de répercuter les états d'alarmes ouvertes sur des entités surveillées.

À partir de Canopsis 4.3.0, il remplace l'ancien moteur `engine-watcher`.

## Utilisation

### Options du moteur

La commande `engine-service -help` liste toutes les options acceptées par le moteur.

### Multi-instanciation

Il est possible, à partir de Canopsis 4.3.0, de lancer plusieurs instances du moteur `engine-service`, afin d'améliorer sa performance de traitement et sa résilience.

En environnement Docker, il vous suffit par exemple de lancer Docker Compose avec `docker-compose up -d --scale service=2` pour que le moteur `engine-service` soit lancé avec 2 instances.

Cette fonctionnalité sera aussi disponible en installation par paquets lors d'une prochaine mise à jour.

## Fonctionnement

### Concept d'un service

Un service (ou *service*) représente un groupe de surveillance : c'est-à-dire que la criticité d'une entité de type service dépendra de la criticité des entités surveillées, et des alarmes ouvertes sur ces entités.

Le but d'un service est de donner une visibilité accrue et claire sur l'état d'un groupe d'entités, afin de détecter un changement de criticité positif ou négatif sur les alarmes liées aux entités du groupe surveillé.

### Templates

L'`output_template` est un [template](https://golang.org/pkg/text/template/) permettant d'afficher diverses informations dans l'output de l'alarme correspondant à l'service.

Les informations disponibles sont :

 - `{{.Alarms}}` : le nombre d'alarmes en cours sur les entités observées par l'service.
 - `{{.State.Info}}` : le nombre d'entités observées n'ayant pas d'alarmes, ou une alarme en criticité `Info`.
 - `{{.State.Minor}}` : le nombre d'alarmes mineures sur les entités observées.
 - `{{.State.Major}}` : le nombre d'alarmes majeures sur les entités observées.
 - `{{.State.Critical}}` : le nombre d'alarmes critiques sur les entités observées.
 - `{{.Acknowledged}}` : le nombre d'alarmes acquittées sur les entités observées.
 - `{{.NotAcknowledged}}` : le nombre d'alarmes non-acquittées sur les entités observées.

Par exemple, l'output d'un service avec l'`output_template` suivant :

```
Crit : {{.State.Critical}} / Total : {{.Alarms}}
```

sera

```
Crit : 12 / Total : 60
```

s'il y a 60 alarmes en cours dont 12 critiques sur les entités observées.
