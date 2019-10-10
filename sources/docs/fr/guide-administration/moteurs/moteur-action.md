# Action

Le moteur action permet de déclencher conditionnellement des actions sur des alarmes.

Les actions sont définies dans la collection MongoDB `default_action`, et peuvent être ajoutées et modifiées avec l'[API Action](../../guide-developpement/action/api_v2_action.md).

## Fonctionnement

La file du moteur est placée juste après le moteur [Axe](moteur-axe.md).

Les types d'actions disponibles sont :

* `pbehavior`, qui va poser un [PBehavior](moteur-pbehavior.md)
* `snooze`, qui va poser des snooze automatiques sur les alarmes lors de leur création.

Une action est composée d'un JSON contenant les paramètres suivants :

- `_id` (optionnel) : l'identifiant du webhook (généré automatiquement ou choisi par l'utilisateur).
- `type` : `[pbehavior](../../guide-developpement/moteurs/moteur-pbehavior.md)` ou `snooze`.
- `hook` (requis) : les conditions dans lesquelles le webhook doit être appelé, dont :
    - `alarm_patterns` (optionnel) : Liste de patterns permettant de filtrer les alarmes.
    - `entity_patterns` (optionnel) : Liste de patterns permettant de filtrer les entités.
    - `event_patterns` (optionnel) : Liste de patterns permettant de filtrer les évènements. Le format des patterns est le même que pour l'[event-filter](moteur-che-event_filter.md).
    - [`triggers`](../architecture-interne/triggers.md) (requis) : Liste de [triggers](../architecture-interne/triggers.md). Au moins un de ces [triggers](../architecture-interne/triggers.md) doit avoir eu lieu pour que le webhook soit appelé.
- `parameters` (requis) : les informations nécessaires correspondant au type d'action.

## Collection

Les actions sont stockées dans la collection MongoDB `default_action` (voir [API Action](../../guide-developpement/action/api_v2_action.md) pour la création d'actions). Le champ `type` de l'objet définit le type d'action. Par exemple, avec un pbehavior, le champ `type` vaut `pbehavior` :

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

Un exemple d'action concernant le snooze automatique (le `type` d'action est donc `snooze`). Il a lieu à la création de l'alarme et si le champ `resource` de l'événement contient les termes `CPU`  ou `HDD`.

Dans les `parameters`, on définit la durée du snooze (600 secondes, soit 10 minutes dans cet exemple), l'auteur et le message accompagnant le snooze.

```JSON
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
