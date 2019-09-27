# Alerts - Alarm filter

!!! note
    Cette page concerne l'alarm-filter disponible dans le moteur Python `alerts`. Le moteur [action](moteur-action.md) propose des fonctionnalités similaires pour une stack Go.

L'alarm-filter est une fonctionnalité du moteur alerts permettant de déclencher conditionnellement des actions lors de la création d'alarmes.

Les actions sont définies dans la collection MongoDB `default_alarmfilter`, et peuvent être ajoutées et modifiées avec l'[API alarm-filter](../../guide-developpement/alarm-filter/api_v2_alarm_filter.md).

## Fonctionnement

Une action est un document JSON contenant les paramètres suivants :

 - `entity_filter` (requis) : un filtre MongoDB appliqué à la collection `periodical_alarm` pour sélectionner les alarmes sur lesquelles l'action doit être exécutée.
 - `condition` (requis): un filtre MongoDB appliqué à la collection `periodical_alarm` pour sélectionner les alarmes sur lesquelles l'action doit être exécutée. Il est recommandé d'utiliser l'`entity_filter` pour filtrer les alarmes, et d'utiliser la valeur `{}` pour le champ `condition`.
 - `tasks` (requis) : une liste de tâches à appliquer à l'alarme.
 - `output_format` (optionnel) : le message à afficher dans la timeline des alarmes.
 - `limit` (requis) : la durée (en secondes) entre la création de l'alarme et l'exécution de l'action, et entre deux exécutions consécutives de l'action.
 - `postpone_if_active_pbehavior` (optionnel, `false` par défaut): `true` pour que l'action ne soit pas exécutée si un pbehavior est actif sur l'alarme, et pour que le décompte du délai `limit` soit réinitialisé en sortie de pbehavior.
 - `repeat` (optionnel, 1 par défaut) : le nombre d'exécutions de l'action.


Les tâches utilisables dans le champ `tasks` sont :

 - `alerts.useraction.ack`
 - `alerts.useraction.ackremove`
 - `alerts.useraction.cancel`
 - `alerts.useraction.comment`
 - `alerts.useraction.uncancel`
 - `alerts.useraction.declareticket`
 - `alerts.useraction.done`
 - `alerts.useraction.assocticket`
 - `alerts.useraction.changestate`
 - `alerts.useraction.keepstate`
 - `alerts.useraction.snooze`
 - `alerts.systemaction.state_increase`
 - `alerts.systemaction.state_decrease`
 - `alerts.systemaction.status_increase`
 - `alerts.systemaction.status_decrease`

## Exemple

L'action suivante augmente automatiquement l'état des alarmes au bout d'une heure.

```json
{
    "entity_filter": "{\"v.state.val\": 1}",
    "condition": {},
    "tasks": [
        "alerts.useraction.state_increase"
    ],
    "output_format": "L'alarme a été créée il y a une heure. Augmentation automatique de son état.",
    "limit": 3600,
    "postpone_if_active_pbehavior": true
}
```