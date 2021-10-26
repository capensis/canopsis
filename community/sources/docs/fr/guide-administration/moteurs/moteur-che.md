# Moteur `engine-che` (Community)

Le moteur `engine-che` permet d'enrichir les événements (via son [`event-filter`](moteur-che-event_filter.md)), de créer et d'enrichir les entités et de créer le context-graph.

## Utilisation

### Options du moteur

La commande `engine-che -help` liste toutes les options acceptées par le moteur.

### Multi-instanciation

Il est possible, à partir de Canopsis 3.39.0, de lancer plusieurs instances du moteur `engine-che`, afin d'améliorer sa performance de traitement et sa résilience.

En environnement Docker, il vous suffit par exemple de lancer Docker Compose avec `docker-compose up -d --scale che=2` pour que le moteur `engine-che` soit lancé avec 2 instances.

Cette fonctionnalité sera aussi disponible en installation par paquets lors d'une prochaine mise à jour.

## Fonctionnement

À l'arrivée dans sa file, le moteur `engine-che` va leur appliquer les règles d'enrichissement de son [`event-filter`](moteur-che-event_filter.md).

Si l'événement est de type `check` ou `declareticket` : au prochain battement (*beat*) du moteur, il va ensuite créer, enrichir ou mettre à jour les entités, puis il va mettre à jour le context-graph qui gère les liens entre les entités.

## Activation des plugins d'enrichissement externe (`datasource`)

La fonctionnalité d'[event-filter](moteur-che-event_filter.md) peut utiliser des sources de données externes (à l'exception de `entity`) afin d'enrichir les évènements, à l'aide de *plugins*. Des plugins `datasource` sont disponibles pour cela, dans Canopsis Pro.

Si vous bénéficiez d'une souscription à Canopsis Pro, la procédure suivante vous permet d'activer ces plugins d'enrichissement externe. Il est néanmoins conseillé de n'appliquer cette procédure qu'en cas de réel besoin des fonctionnalités `datasource`.

=== "En installation Docker"

    Le conteneur `che` doit être modifié pour utiliser l'image `canopsis/engine-che-cat` à la place de l'image `canopsis/engine-che` proposée par défaut. L'accès à l'image `canopsis/engine-che-cat` nécessite une souscription Pro.

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

    Vous devez tout d'abord avoir installé le paquet système `canopsis-engines-go-cat`, qui est uniquement disponible lors d'une souscription à Canopsis Pro.

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
