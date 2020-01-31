# Action

Le moteur action permet de déclencher conditionnellement des actions sur des alarmes.

Les actions sont définies dans la collection MongoDB `default_action`, et peuvent être ajoutées et modifiées avec l'[API Action](../../guide-developpement/api/api-v2-action.md).

## Utilisation

En édition `core`, la file du moteur est placée juste après le moteur [Axe](moteur-axe.md).

En édition `CAT`, la file du moteur est placée juste après le moteur [Webhook](moteur-webhook.md).

### Options de l'engine-action

```
  -d	debug
  -version
      version infos
```

## Fonctionnement

### Types d'action

Les types d'actions disponibles sont :

* `changestate`, qui correspond à un évènement [`changestate`](../../guide-developpement/struct-event.md#event-changestate-structure) : change et verrouille l'état de l'alarme dans une criticité donnée jusqu'à sa résolution
* `pbehavior` : pose un [PBehavior](moteur-pbehavior.md)
* `snooze`, qui correspond à un évènement [`snooze`](../../guide-developpement/struct-event.md#event-snooze-structure) : pose un snooze automatique sur l'alarme

### Paramètres généraux

Une action est composée d'un JSON contenant les paramètres suivants :

```javascript
{
"_id"        // identifiant de l'action , optionnel, s'il n'est pas fourni par l'utilisateur il sera généré automatiquement - le champ est de type `string`
"type"       // type du pbehavior (`changestate`, `pbehavior` ou `snooze`), obligatoire - le champ est de type `string`
"hook"       // conditions sur les champs des alarmes (`alarm_patterns`), des entités (`entity_patterns`) ou des évènements (`event_patterns`) dans lesquelles l'action doit être appelée, optionnel
"triggers"   // conditions de déclenchement sur la vie de l'alarme, si plusieurs triggers sont indiqués, au moins un de ces triggers doit avoir eu lieu pour que l'action soit appelée
"parameters" // paramétrage spécifique à chaque type d'action.
}
```

!!! attention
    Les [`triggers`](../architecture-interne/triggers.md) `declareticketwebhook`, `resolve` et `unsnooze` n'étant pas déclenchés par des [évènements](../../guide-developpement/struct-event.md), ils ne sont pas utilisables avec les `event_patterns`.

### Paramètres spécifiques

#### Changestate

```javascript
{
"state":   // état dans lequel sera verrouillée l'alarme (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL), le champ est de type `integer`
"output":  // commentaire du changestate, optionnel - le champ est de type `string`
"author":  // auteur du changestate, optionnel - le champ est de type `string`
}
```

#### PBehavior

```javascript
{
"rrule":     // règle de récurrence pour le pbehavior, optionnel - le champ est de type `string`
"enabled":   // détermine si le pbehavior est actif ou non, obligatoire - le champ est de type `boolean`
"author":    // auteur du pbehavior, obligatoire - le champ est de type `string`
"name":      // nom du pbehavior, obligatoire - le champ est de type `string`
"tstart":    // date de début du pbehavior, obligatoire - le champ est de type `integer` en UNIX timestamp
"tstop":     // date de fin du pbehavior, obligatoire - le champ est de type `integer` en UNIX timestamp
"type_":     // type du pbehavior, obligatoire - le champ est de type `string`
"reason":    // raison du pbehavior, obligatoire - le champ est de type `string`
"timezone":  // timezone du pbehavior, optionnel - le champ est de type `string`
"comments":  // commentaire du pbehavior, optionnel - le champ est de type `array`
[{           // début de l'array du commentaire
"author":  // auteur du commentaire, optionnel - le champ est de type `string`
"message": // commentaire du pbehavior, optionnel - le champ est de type `string`
}],
"exdate": []   // date d'expiration, peut être composé de plusieurs valeurs, optionnel - le champ est de type integer` en UNIX timestamp
}
```

#### Snooze

```javascript
{
"message":   // commentaire du snooze, optionnel - le champ est de type `string`
"duration":  // durée du snooze en secondes, optionnel - le champ est de type `number`
"author":    // auteur du snooze, optionnel - le champ est de type `string`
}
```

## Collection

Les actions sont stockées dans la collection MongoDB `default_action` (voir [API Action](../../guide-developpement/api/api-v2-action.md) pour la création d'actions). Le champ `type` de l'objet définit le type d'action. Par exemple, avec un pbehavior, le champ `type` vaut `pbehavior` :

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

Un exemple d'action concernant le snooze automatique (le `type` d'action est donc `snooze`). Il a lieu à la création de l'alarme et si le champ `resource` de l'évènement contient les termes `CPU`  ou `HDD`.

Dans les `parameters`, on définit la durée du snooze (600 secondes, soit 10 minutes dans cet exemple), l'auteur et le message accompagnant le snooze.

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
