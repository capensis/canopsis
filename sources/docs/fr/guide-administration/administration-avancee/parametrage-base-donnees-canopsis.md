# Paramétrage en base de données de Canopsis

Certaines configurations de Canopsis peuvent nécessiter d'effectuer des modifications directement dans la base de données MongoDB.

## Activation du champ `last_event_date`

!!! attention
    Activer cette option entraîne une action supplémentaire systématique au [moteur axe](../moteurs/moteur-axe.md) et a un impact négatif sur ses performances.

Les alarmes dans Canopsis incluent un champ `alarm.v.last_event_date`.

Cependant, la mise-à-jour de ce champ n'est pas activée par défaut. Sa valeur est celle de `alarm.v.creation_date`, soit la date de création de l'alarme par le [moteur axe](../moteurs/moteur-axe.md)

### Modification du fichier de configuration

Dans la collection `configuration` de MongoDB, éditer le fichier dont l'`_id` est `global_config`.

Modifier la valeur du champ de `alarm.enablelasteventdate` pour `true` (en booléen).

Par exemple, on passe de :

```json
{
    "_id" : "global_config",
    "global" : {
        "prefetchcount" : 10000,
        "prefetchsize" : 0
    },
    "alarm" : {
        "flappingfreqlimit" : 0,
        "flappinginterval" : 0,
        "stealthyinterval" : 0,
        "baggottime" : "60s",
        "enablelasteventdate" : false
    }
}
```

à

```json
{
    "_id" : "global_config",
    "global" : {
        "prefetchcount" : 10000,
        "prefetchsize" : 0
    },
    "alarm" : {
        "flappingfreqlimit" : 0,
        "flappinginterval" : 0,
        "stealthyinterval" : 0,
        "baggottime" : "60s",
        "enablelasteventdate" : true
    }
}
```

### Redémarrage du moteur Axe

Redémarrer le moteur [moteur axe](../moteurs/moteur-axe.md) pour lui faire prendre en compte le changement de configuration.

Après avoir effectué ces modifications, la colonne dont la valeur est `alarm.v.last_event_date` affiche bien la date de dernière réception de l'événement.
