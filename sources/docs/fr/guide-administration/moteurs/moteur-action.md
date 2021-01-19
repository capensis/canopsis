# Moteur `engine-action` (Go, Core)

Le moteur `engine-action` permet de déclencher conditionnellement des actions sur des alarmes.

Les actions sont définies dans la [collection MongoDB](#collection-mongodb-associee) `default_action` et peuvent être ajoutées et modifiées avec l'[API Action](../../guide-developpement/api/api-v2-action.md).

## Utilisation

En édition Core, la file du moteur est placée juste après le moteur [`engine-axe`](moteur-axe.md).

En édition CAT, la file du moteur est placée juste après le moteur [`engine-webhook`](moteur-webhook.md).

Les actions ack, ackremove, assocticket, declareticket, et cancel sont disponibles. Par ailleurs, il existe aussi la possibilité de déclencher les actions après un délai paramétré.

### Options du moteur

La commande `engine-action -help` liste toutes les options acceptées par le moteur.

## Fonctionnement

### Types d'action

Les types d'actions disponibles sont :

* `changestate`, qui correspond à un évènement [`changestate`](../../guide-developpement/struct-event.md#event-changestate-structure) : change et verrouille la criticité de l'alarme jusqu'à sa résolution.
* `pbehavior`, met en place un [comportement périodique](moteur-pbehavior.md).
* `snooze`, qui correspond à un évènement [`snooze`](../../guide-developpement/struct-event.md#event-snooze-structure) : pose une mise en veille automatique sur l'alarme.
* `ack`, qui correspond à un événement [`ack`](../../guide-developpement/struct-event.md#event-acknowledgment-structure) : pose un acquittement sur l'alarme.
* `ackremove`, qui correspond à un événement [`ackremove`](../../guide-developpement/struct-event.md#event-ackremove-structure) : supprime l'acquittement sur l'alarme.
* `assocticket`, qui correspond à un événement [`assocticket`](../../guide-developpement/struct-event.md#event-assocticket-structure) : associe un ticket à l'alarme.
* `declareticket`, qui correspond à un événement [`declareticket`](../../guide-developpement/struct-event.md#event-declareticket-structure) : déclarer un ticket pour l'alarme.
* `cancel`, qui correspond à un événement [`ackremove`](../../guide-developpement/struct-event.md#event-cancel-structure) : annule l'alarme.

### Paramètres généraux

Une action est composée d'un JSON contenant les paramètres suivants :

* `_id` : optionnel. Identifiant de l'action. S'il n'est pas fourni par l'utilisateur il sera généré automatiquement. Le champ est de type `string`.
* `type` : obligatoire. Type d'action (voir [section précédente](#types-daction)). Ce champ est de type `string`.
* `parameters` : obligatoire. [Paramétrage spécifique](#parametres-specifiques) à chaque type d'action.
* `delay` : optionnel. Délai avant l'exécution de l'action. Les unités acceptées sont celles utilisées par le langage [Golang](https://golang.org/pkg/time/#ParseDuration) soit `s`, `m`, `h` pour secondes, minutes et heures respectivement. Le champ est de type `string`.
* `hook` : obligatoire. Il est composé des paramètres suivants :
    * [`patterns`](moteur-che-event_filter.md#patterns) : optionnel. Conditions sur les champs des alarmes (`alarm_patterns`), des entités (`entity_patterns`) ou des évènements (`event_patterns`) dans lesquelles l'action doit être appelée.
    * [`triggers`](../architecture-interne/triggers.md) : obligatoire. Ils servent comme point de déclenchement pour les actions automatisées, en général lors de la réception d'un évènement.

!!! attention
    Les [`triggers`](../architecture-interne/triggers.md) `declareticketwebhook`, `resolve` et `unsnooze` n'étant pas déclenchés par des [évènements](../../guide-developpement/struct-event.md), ils ne sont pas utilisables avec les `event_patterns`.

### Paramètres spécifiques

#### Changestate

* `state` : obligatoire. Criticité dans laquelle sera verrouillée l'alarme (0 : Info, 1 : Minor, 2 : Major, 3 : Critical). Le champ est de type `integer`.
* `output` : optionel. Commentaire du changestate. Le champ est de type `string`.
* `author` : optionel. Auteur du changestate. Le champ est de type `string`.

#### PBehavior

* `rrule` : optionnel. Règle de récurrence pour le pbehavior. Le champ est de type `string`.
* `enabled` : obligatoire. Détermine si le pbehavior est actif ou non. Le champ est de type `boolean`.
* `author` : obligatoire. Auteur du pbehavior. Le champ est de type `string`.
* `name` : obligatoire. Nom du pbehavior. Le champ est de type `string`.
* `tstart` : obligatoire. Date de début du pbehavior. Le champ est de type `integer` en UNIX timestamp.
* `tstop` : obligatoire. Date de fin du pbehavior. Le champ est de type `integer` en UNIX timestamp.
* `type_` : obligatoire. Type du pbehavior. Le champ est de type `string`.
* `reason` : obligatoire. Raison du pbehavior. Le champ est de type `string`.
* `timezone` : optionnel. Timezone du pbehavior. Le champ est de type `string`.
* `comments` : optionnel. Commentaire du pbehavior. Le champ est de type `array`.
* `author` : optionnel. Auteur du commentaire. Le champ est de type `string`.
* `message` : optionnel. Commentaire du pbehavior. Le champ est de type `string`.
* `exdate` : optionnel. Date d'expiration, peut être composée de plusieurs valeurs. Le champ est de type `integer` en UNIX timestamp.

#### Snooze

* `message` : optionnel. Commentaire de la mise en veille. Le champ est de type `string`.
* `duration` : optionnel. Durée de la mise en veille en secondes. Le champ est de type `integer`.
* `author`: optionnel. Auteur de la mise en veille. Le champ est de type `string`.

#### Ack

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Ackremove

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Assocticket

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Declareticket

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Cancel

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

## Collection MongoDB associée

Les actions sont stockées dans la collection MongoDB `default_action` (voir [API Action](../../guide-developpement/api/api-v2-action.md) pour la création d'actions). Le champ `type` de l'objet définit le type d'action. Par exemple, avec un comportement périodique, le champ `type` vaut `pbehavior` :

```json
{
    "_id" : "xyz",
    "type": "pbehavior",
    "hook": {
        "event_patterns": [
            {
                "resource": "CPU_2"
            },
            {
                "resource": "HDD_2"
            }
        ],
        "triggers": [
            "create"
        ]
    },
    "delay" : "1m",
    "parameters" : {
        "author" : "whalefact",
        "name" : "Big",
        "reason" : "Most whales are legally unemployed",
        "type" : "Pause",
        "rrule" : "",
        "tstart" : 0,
        "tstop" : 253402297199,
    }
}
```

Un exemple d'action concernant la mise en veille automatique (le `type` d'action est donc `snooze`). Il a lieu à la création de l'alarme et si le champ `resource` de l'évènement contient les termes `CPU`  ou `HDD`.

Dans les `parameters`, on définit la durée de la mise en veille (600 secondes, soit 10 minutes dans cet exemple), l'auteur et le message accompagnant la mise en veille.

```json
{
	"_id": "temporisation-10m",
	"type": "snooze",
	"hook": {
		"event_patterns": [{
				"resource": {
					"regex_match": "CPU"
				}
			},
			{
				"resource": {
					"regex_match": "HDD"
				}
			}
		],
		"triggers": [
			"create"
		]
	},
	"parameters": {
		"author": "action",
		"message": "Temporisation de l'alarme pendant 10 minutes",
		"duration": 600
	}
}
```
