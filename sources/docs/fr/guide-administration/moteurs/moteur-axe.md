# Axe

Le moteur axe permet de créer et d'enrichir les alarmes. Il permet également d'appliquer les actions entrées depuis le bac à alarmes.

Dans la version CAT, il permet aussi d'appliquer des [`webhooks`](moteur-axe-webhooks.md).

## Fonctionnement

La file du moteur est placée juste après le moteur [che](moteur-che.md).

À l'arrivée dans sa file, le moteur axe va transformer les événements en alarmes qu'il va créer et enrichir.

Si l'événement correspond à une alarme en cours, l'alarme va alors être mise à jour.

Si l'événement en correspond à aucune alarme en cours, l'alarme va alors être créée.

Si l'événement correspond à une action (comme la mise d'un ACK), l'alarme va être mise à jour en prenant en compte l'action.

### Options de l'engine-axe

```
  -autoDeclareTickets
        Déclare les tickets automatiquement pour chaque alarme. DÉPRÉCIÉ, remplacé par les webhooks.
  -d    debug
  -featureHideResources
        Active les features de gestion de ressources cachées.
  -featureStatEvents
        Envoie les évènements de statistiques
  -postProcessorsDirectory
        Le répetoire contenant les plugins de post-traitement (par défaut ".")
  -printEventOnError
        Afficher les évènements sur les erreurs de traitement.
  -publishQueue
        Publie les événements sur cette queue. (par défaut "Engine_watcher")
  -version
        version infos
```

## Collection

Les alarmes sont stockées dans la collection Mongo `periodical_alarm`.

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
