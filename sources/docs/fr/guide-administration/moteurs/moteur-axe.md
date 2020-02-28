# Axe

Le moteur axe permet de créer et d'enrichir les alarmes. Il permet également d'appliquer les actions entrées depuis le bac à alarmes.

Jusqu'en `3.33.0`, dans la version CAT, il permettait aussi d'appliquer des [`webhooks`](moteur-webhook.md).

Depuis la `3.34.0`, les [`webhooks`](moteur-webhook.md) sont devenus leur propre moteur (disponible uniquement en version CAT).

## Utilisation

La file du moteur est placée juste après le moteur [che](moteur-che.md).

### Options de l'engine-axe

```
  -d	debug
  -featureHideResources
      Enable Hide Resources Management
  -featureStatEvents
      Send statistic events
  -ignoreDefaultTomlConfig
    	load toml file values into database
  -postProcessorsDirectory string
      The path of the directory containing the post-processing plugins. (default ".")
  -printEventOnError
      Print event on processing error
  -publishQueue string
      Publish event to this queue. (default "Engine_watcher")
  -version
      version infos
```


### Fichier de configuration

Au premier démarrage, le moteur `axe` lit le fichier de configuration `default_configuration.toml`.  
Il peut relire ce fichier si le flag `-ignoreDefaultTomlConfig` est positionné et ainsi écraser les informations de configuration en base de données.


````
[global]
PrefetchCount = 10000
PrefetchSize = 0

[alarm]
FlappingFreqLimit = 0
FlappingInterval = 0
StealthyInterval = 0
BaggotTime = "60s"
EnableLastEventDate = false
CancelAutosolveDelay = "1h"
````

* Activation du champ `last_event_date`

!!! attention
    Activer cette option entraîne une action supplémentaire systématique au [moteur axe](../moteurs/moteur-axe.md) et a un impact négatif sur ses performances.

Les alarmes dans Canopsis incluent un champ `alarm.v.last_event_date`.

Cependant, la mise-à-jour de ce champ n'est pas activée par défaut. Sa valeur est celle de `alarm.v.creation_date`, soit la date de création de l'alarme par le [moteur axe](../moteurs/moteur-axe.md).  
Pour l'activer, passez le paramètre `EnableLastEventDate` à `true`.  


* Délai de résolution d'une alarme annulée manuellement

Lorsqu'une alarme est annulée manuellement, via l'interface web par exemple, elle n'est marquée `résolue` qu'après un délai d'une heure par défaut.  
Vous pouvez agir sur ce délai en modifiant le paramètre `CancelAutosolveDelay`.


## Fonctionnement

À l'arrivée dans sa file, le moteur axe va transformer les événements en alarmes qu'il va créer et mettre à jour.

### Événements de type check

3 possibilités pour un événement de type [`check`](../../guide-developpement/struct-event.md#event-check-structure) :

* Il ne correspond à aucune alarme en cours : l'alarme va alors être créée
* Il correspond à une alarme en cours et son champ `state` ne vaut pas `0` : l'alarme va alors être mise à jour
* Il correspond à une alarme en cours et son champ `state` vaut `0` : l'alarme va alors passer en état `OK`. Au 2° [battement (beat)](../../guide-utilisation/vocabulaire/index.md#battement) suivant, si l'alarme n'a pas été rouverte par un nouvel événement de type [`check`](../../guide-developpement/struct-event.md#event-check-structure), elle est considérée comme résolue. Un champ `v.resolved` lui est alors ajouté avec le timestamp courant.

### Autres types d'événements

Si l'événement correspond à une action (comme la mise d'un [`ACK`](../../guide-developpement/struct-event.md#event-acknowledgment-structure)), l'alarme va être mise à jour en appliquant l'action.

## Collection

Les alarmes sont stockées dans la collection MongoDB `periodical_alarm`.

L'`_id` est générée automatiquement.

Le champ `d` correspond à l'`_id` de l'entité à laquelle l'alarme est rattachée.

```json
{
    "_id" : "aad73d0b-2e0e-453d-90c5-1c843cd196b2",
    "t" : NumberLong(1567498879),
    "d" : "disk2/serveur_de_salle_machine_DHCP",
    "v" : {
        "state" : {
            "_t" : "stateinc",
            "t" : NumberLong(1567498879),
            "a" : "superviseur1.superviseur1",
            "m" : "Disque plein a 98%, 50GO occupe",
            "val" : NumberLong(2)
        },
        "status" : {
            "_t" : "statusinc",
            "t" : NumberLong(1567498879),
            "a" : "superviseur1.superviseur1",
            "m" : "Disque plein a 98%, 50GO occupe",
            "val" : NumberLong(1)
        },
        "steps" : [
            {
                "_t" : "stateinc",
                "t" : NumberLong(1567498879),
                "a" : "superviseur1.superviseur1",
                "m" : "Disque plein a 98%, 50GO occupe",
                "val" : NumberLong(2)
            },
            {
                "_t" : "statusinc",
                "t" : NumberLong(1567498879),
                "a" : "superviseur1.superviseur1",
                "m" : "Disque plein a 98%, 50GO occupe",
                "val" : NumberLong(1)
            }
        ],
        "component" : "serveur_de_salle_machine_DHCP",
        "connector" : "superviseur1",
        "connector_name" : "superviseur1",
        "creation_date" : NumberLong(1567498879),
        "display_name" : "XA-KU-AQ",
        "extra" : {},
        "initial_output" : "Disque plein a 98%, 50GO occupe",
        "output" : "Disque plein a 98%, 50GO occupe",
        "initial_long_output" : "",
        "long_output" : "",
        "long_output_history" : [
            ""
        ],
        "last_update_date" : NumberLong(1567498879),
        "last_event_date" : NumberLong(1567498879),
        "resource" : "disk2",
        "state_changes_since_status_update" : NumberLong(0),
        "tags" : [],
        "total_state_changes" : NumberLong(1)
    }
}
```
