# Moteur `engine-axe`

Le moteur `engine-axe` permet de créer et d'enrichir les alarmes. Il permet également d'appliquer les actions entrées depuis le Bac à alarmes. Il fait partie des moteurs Go nouvelle génération.

Jusqu'en 3.33.0, `engine-axe` permettait aussi d'appliquer des Webhooks, dans la version CAT. Depuis Canopsis 3.34.0, les Webhooks sont gérés par un moteur [`engine-webhook`](moteur-webhook.md) dédié (toujours en édition CAT).

## Utilisation

### Options du moteur

La commande `engine-axe -help` liste toutes les options acceptées par le moteur.

Les options acceptées par la dernière version de Canopsis sont les suivantes :

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

## Fichier de configuration

Lors de son tout premier démarrage, le moteur `engine-axe` lit le fichier de configuration `/opt/canopsis/etc/default_configuration.toml` (ou `/default_configuration.toml` en environnement Docker) et inscrit ces informations en base de données.

À partir de Canopsis 3.37.0, l'option `-ignoreDefaultTomlConfig` permet de forcer le moteur à prendre en compte toutes les nouvelles mises à jour de son fichier de configuration, après un redémarrage. Si cette option n'est pas précisée, `engine-axe` synchronisera sa configuration en base uniquement à son premier lancement.

Le contenu par défaut de ce fichier de configuration est le suivant :

```ini
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
```

### Option `EnableLastEventDate`

!!! attention
    Activer cette option entraîne une action supplémentaire systématique dans le moteur qui a une incidence négative sur ses performances.

Les alarmes dans Canopsis incluent un champ `alarm.v.last_event_date`.

Cependant, la mise à jour de ce champ n'est pas activée par défaut. Sa valeur est celle de `alarm.v.creation_date`, soit la date de création de l'alarme par `engine-axe`.

Pour l'activer, passez le paramètre `EnableLastEventDate` du fichier de configuration à `true`.

### Option `CancelAutosolveDelay`

Lorsqu'une alarme est annulée manuellement, via l'interface web par exemple, elle n'est marquée « résolue » qu'après un certain délai, d'une heure par défaut.  

Vous pouvez agir sur ce délai en modifiant le paramètre `CancelAutosolveDelay`.

## Fonctionnement du moteur

La file du moteur est placée juste après le moteur [`engine-che`](moteur-che.md).

À l'arrivée dans sa file, le moteur `engine-axe` va transformer les événements en alarmes qu'il va créer et mettre à jour.

### Événements de type check

3 possibilités pour un événement de type [`check`](../../guide-developpement/struct-event.md#event-check-structure) :

1. Il ne correspond à aucune alarme en cours : l'alarme va alors être créée
2. Il correspond à une alarme en cours et son champ `state` ne vaut pas `0` : l'alarme va alors être mise à jour
3. Il correspond à une alarme en cours et son champ `state` vaut `0` : l'alarme va alors passer en état `OK`. Au 2° [battement (beat)](../../guide-utilisation/vocabulaire/index.md#battement) suivant, si l'alarme n'a pas été rouverte par un nouvel événement de type [`check`](../../guide-developpement/struct-event.md#event-check-structure), elle est considérée comme résolue. Un champ `v.resolved` lui est alors ajouté avec le timestamp courant.

### Autres types d'événements

Si l'événement correspond à une action (comme la mise d'un [`ACK`](../../guide-developpement/struct-event.md#event-acknowledgment-structure)), l'alarme va être mise à jour en appliquant l'action.

## Collection MongoDB associée

Les alarmes sont stockées dans la collection MongoDB `periodical_alarm`.

Le champ `_id` est généré automatiquement.

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
