# Action

Le moteur action permet de déclencher des actions conditionnellement sur des alarmes.

## Fonctionnement

La queue du moteur est placée juste après le moteur Axe.

Une action se déclenche si un des champs d'une alarme (désigné dans `fields`) va correspondre à une regex (`regex`, voir [re2 syntax](https://github.com/google/re2/wiki/Syntax)).

Les types d'actions disponibles sont :

* `pbehavior`, qui va poser un [PBehavior](../../guide-developpement/PBehavior/index.md) ;
* `snooze`, qui va poser des snooze automatiques sur les alarmes lors de leur création.

## Collection

Les actions sont stockées dans la collection Mongo `default_action` (voir [API Action](../../guide-developpement/action/api_v2_action.md) pour la création d'actions). Le champ `type` de l'objet définit le type d'action. Par exemple, avec un pbehavior, le champ `type` vaut `pbehavior` :

```json
{
    "_id" : "xyz",
    "type": "pbehavior",
    "fields" : ["Resource", "Component"],
    "regex" : ".*whale.*",
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

Un exemple d'action concernant le snooze automatique (le `type` d'action est donc `snooze`). Les champs `fields` et `regex` indiquent sur quelles actions on veut peut poser le snooze. Dans les `parameters`, on définit la durée du snooze (600 secondes, soit 10 minutes dans cet exemple), l'auteur et le message accompgnant le snooze.

```json
{
    "_id" : "temporisation-10m",
    "type": "snooze",
    "fields" : ["Resource", "Component"],
    "regex" : "(FS|HARDWARE)",
    "parameters" : {
        "author" : "action",
        "message" : "Temporisation de l'alarme pendant 10 minutes",
        "duration" : 600
    }
}
```

!!! attention
    Les valeurs dans le tableau des `fields` sont sensibles à la casse, il faut utiliser les majuscules (`Resource`, `Component`, etc).
