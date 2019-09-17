# Che

Le moteur che permet d'enrichir les événements (via son [`event-filter`](moteur-che-event_filter.md)), de créer et d'enrichir les entités et de créer le context-graph.

## Fonctionnement

La file du moteur est placée juste après l'exchange `canopsis.events`.

À l'arrivée dans sa file, le moteur che va leur appliquer les règles d'enrichissement de son [`event-filter`](moteur-che-event_filter.md).

Il va ensuite créer, enrichir ou mettre à jour les entités, puis il va mettre à jour le context-graph qui gère les liens entre les entités.

### Options de l'engine-che

```
  -consumeQueue string
        Consomme les évènements venant de cette queue. (default "Engine_che").
  -createContext
        Active la création de context graph. Activé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorse que vous l'utilisez. (default true)
  -d    debug
  -dataSourceDirectory
        The path of the directory containing the event filter's data source plugins. (default ".")
  -enrichContext
        Active l'enrichissment de context graph à partir d'un event. Désactivé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorse que vous l'utilisez. (default true)
  -enrichExclude string
        Liste de champs séparés par des virgules ne faisant pas partie de l'enrichissement du contexte
  -enrichInclude string
        Coma separated list of the only fields that will be part of context enrichment. If present, -enrichExclude is ignored.
  -printEventOnError
        Print event on processing error
  -processEvent
        enable event processing. enabled by default. (default true)
  -publishQueue
        Publie les événements sur cette queue. (default "Engine_event_filter")
  -purge
        purge consumer queue(s) before work
  -version
        version infos
```

## Collection

Les entités sont stockées dans la collection Mongo `default_entities`.

Le champ `type` de l'objet définit le type de l'entité. Par exemple, avec une ressource, son champ `type` vaut `resource` :

```json
{
    "_id" : "disk2/pbehavior_test_1",
    "name" : "disk2",
    "impact" : [
        "pbehavior_test_1"
    ],
    "depends" : [
        "superviseur1/superviseur1"
    ],
    "enable_history" : [
        NumberLong(1567437797)
    ],
    "measurements" : null,
    "enabled" : true,
    "infos" : {},
    "type" : "resource"
}
```
