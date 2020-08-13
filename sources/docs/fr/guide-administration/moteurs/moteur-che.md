# Moteur `engine-che` (Go, Core)

Le moteur `engine-che` permet d'enrichir les [événements](../../guide-developpement/struct-event.md) (via son [`event-filter`](moteur-che-event_filter.md)), de créer et d'enrichir les entités et de créer le context-graph.

## Utilisation

La file du moteur est placée juste après l'exchange `canopsis.events`.

### Options du moteur

```
  -consumeQueue string
        Consomme les évènements venant de cette file. (default "Engine_che").
  -createContext
        Active la création de context graph. Activé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorsque vous l'utilisez. (default true)
  -d    debug
  -dataSourceDirectory
        The path of the directory containing the event filter's data source plugins. (default ".")
  -enrichContext
        Active l'enrichissment de context graph à partir d'un event. Désactivé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorsque vous l'utilisez. (default true)
  -enrichExclude string
        Liste de champs séparés par des virgules ne faisant pas partie de l'enrichissement du contexte
  -enrichInclude string
        Coma separated list of the only fields that will be part of context enrichment. If present, -enrichExclude is ignored.
  -printEventOnError
        Print event on processing error
  -processEvent
        enable event processing. enabled by default. (default true)
  -publishQueue
        Publie les événements sur cette file. (default "Engine_event_filter")
  -purge
        purge consumer queue(s) before work
  -version
        version infos
```

### Multi-instanciation

!!! note
    Cette fonctionnalité est disponible à partir de Canopsis 3.39.0. Elle ne doit pas être utilisée sur les versions antérieures.

Il est possible, à partir de **Canopsis 3.39.0**, de lancer plusieurs instances du moteur `engine-che`, afin d'améliorer sa performance de traitement et sa résilience.

En environnement Docker, il vous suffit par exemple de lancer Docker Compose avec `docker-compose up -d --scale che=2` pour que le moteur `engine-che` soit lancé avec 2 instances.

Cette fonctionnalité sera aussi disponible en installation par paquets lors d'une prochaine mise à jour.

## Fonctionnement

À l'arrivée dans sa file, le moteur `engine-che` va leur appliquer les règles d'enrichissement de son [`event-filter`](moteur-che-event_filter.md).

Si l'événement est de type [`check`](../../guide-developpement/struct-event.md#event-check-structure) ou [`declareticket`](../../guide-developpement/struct-event.md#event-declareticket-structure) : au prochain battement (*beat*) du moteur, il va ensuite créer, enrichir ou mettre à jour les entités, puis il va mettre à jour le context-graph qui gère les liens entre les entités.

## Activation des plugins d'enrichissement externe (`datasource`)

La fonctionnalité d'[event-filter](moteur-che-event_filter.md) peut utiliser des sources de données externes (à l'exception de `entity`) afin d'enrichir les évènements, à l'aide de *plugins*. Des plugins `datasource` sont disponibles pour cela, dans Canopsis CAT.

Si vous bénéficiez d'une souscription à Canopsis CAT, la procédure suivante vous permet d'activer ces plugins d'enrichissement externe.

=== "En installation Docker"

    Le conteneur `che` doit être modifié pour utiliser l'image `canopsis/engine-che-cat` à la place de l'image `canopsis/engine-che` proposée par défaut. L'accès à l'image `canopsis/engine-che-cat` nécessite une souscription CAT.

    Si vous utilisez Docker Compose, adaptez votre section `che` à l'exemple suivant :

    ``` yaml hl_lines="2 6"
      che:
        image: canopsis/engine-che-cat:${CANOPSIS_IMAGE_TAG}
        env_file:
          - compose.env
        restart: unless-stopped
        command: /engine-che -dataSourceDirectory /plugins/datasource
    ```

=== "En installation paquets"

    Vous devez tout d'abord avoir installé le paquet système `canopsis-engines-go-cat`, qui est uniquement disponible lors d'une souscription à Canopsis CAT.

    Le moteur `engine-che` doit ensuite être lancé avec l'option `-dataSourceDirectory /opt/canopsis/lib/go/plugins/datasource`, afin que le plugin `datasource` soit chargé.

    Vous pouvez, pour cela, exécuter les commandes suivantes :

    ``` sh
    mkdir -p /etc/systemd/system/canopsis-engine-go@engine-che.service.d
    cat > /etc/systemd/system/canopsis-engine-go@engine-che.service.d/datasource.conf << EOF
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/%i -dataSourceDirectory /opt/canopsis/lib/go/plugins/datasource
    EOF

    systemctl daemon-reload
    systemctl restart canopsis-engine-go@engine-che
    ```

Voyez ensuite [la documentation de l'event-filter du moteur `engine-che`](moteur-che-event_filter.md#donnees-externes) afin d'en savoir plus sur l'utilisation de cette fonctionnalité.

## Collection MongoDB associée

Les entités sont stockées dans la collection MongoDB `default_entities`.

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
