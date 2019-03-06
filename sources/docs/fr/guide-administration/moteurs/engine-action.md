# Engine Action

L'engine action permet de déclencher des actions conditionnellement à des alarmes.

## Fonctionnement

La queue de l'engine est placée juste après l'engine Axe.

Une action se déclenche si un des champs d'une alarm (désigné dans `fields`) va matcher une regex (`regex`, voir [re2 syntax](https://github.com/google/re2/wiki/Syntax)).

Pour un [PBehavior](../../guide-developpement/PBehavior/index.md), il va être posé en fonction du paramétrage (`parameters`) de l'action.

C'est également ce moteur qui va poser des snooze automatiques sur les alarmes, en fonction des champs `fields` et `regex`. Si l'alarme correspond à ces champs, qu'elle n'a pas déjà été snoozed et qu'elle a été créée il y a moins de `duration` secondes, alors un snooze sera posé sur l'alarme de manière automatique.

Attention, les valeurs dans le tableau des `fields` est sensible à la casse (il faut une majuscule).

## Collection

Les actions sont stockées dans la collection `default_action`. Par exemple, avec un pbehavior:

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

Un exemple d'action concernant le snooze automatique (le `type` d'action est donc `snooze`). Les champs `fields` et `regex` indiquent sur quelles actions on veut peut poser le snooze. Dans les `parameters`, on définit la durée du snooze (10 minutes dans cet exemple), l'auteur et le message accompgnant le snooze.

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

!!! note
    Si le système de [webhooks](../webhooks/index.md) est utilisé, il est vivement conseillé de créer un webhook contenant le trigger `unsnooze`. Ce webhook va gérer les alarmes à la fin de leur snooze.

