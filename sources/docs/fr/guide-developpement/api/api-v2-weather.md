# Météo des services

L'API météo des services permet de consulter l'état des
[watchers](../../guide-administration/moteurs/moteur-watcher.md). Elle est utilisée par le widget [météo des services](../../guide-utilisation/interface/widgets/index.md#meteo-de-services).

### Récupération d'une liste de watchers

**URL** : `/api/v2/weather/watchers/<filtre>`

`<filtre>` est un filtre MongoDB, utilisé sur la collection `default_entities`. Il permet en particulier de filtrer sur le nom du watcher (`name`) et ses informations enrichies (`infos.*.value`).

**Méthode** : `GET`

**Paramètres** :

 - `limit` : le nombre de watchers à renvoyer.
 - `start` : le nombre de watchers à passer.
 - `orderby` : le nom du champ à utiliser pour le tri. Les champs utilisables
   sont les mêmes que pour le filtre.
 - `direction` : l'ordre dans lequel les watchers doivent être trié (`ASC` ou
   `DESC`)

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les watchers :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/weather/watchers/\{\}'
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer les watchers dont l'information `customer` vaut `Capensis`.

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/weather/watchers/\{"infos.customer.value":"Capensis"\}'
```

### Réponse en cas de réussite

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
[
    {
        "entity_id": "watcher_id",
        "infos": {
            "customer": {
                "description": "Nom du client",
                "value": "Capensis"
            },
            // ...
        },
        "display_name": "Nom du watcher",
        "mfilter": "...",
        "watcher_pbehavior": [],  // Liste des pbehaviors actifs sur le watcher
        "automatic_action_timer": null,

        // Derniers steps state, status, snooze et ack de l'alarme du watcher
        "state": {...},   // contient {"val": 0} s'il n'y a pas d'alarme
        "status": {...},  // non défini s'il n'y a pas d'alarme
        "snooze": {...},  // non défini s'il n'y a pas d'alarme
        "ack": {...},     // non défini s'il n'y a pas d'alarme

        // Champs de l'alarme (ces champs ne sont pas définis s'il n'y a pas d'alarme)
        "connector": "...",
        "connector_name": "...",
        "component": "...",
        "resource": "...",
        "last_update_date": ...,

        "isActionRequired": false,     // true si le watcher est impacté par une alarme non-acquittée
        "isAllEntitiesPaused": false,  // true si toutes les dépendances ont un pbehavior actif
        "isWatcherPaused": false,      // true si le watcher a un pbehavior actif
        "tileColor": "ok",             // le nom de la couleur de la tuile : pause, ok, minor, major ou critical
        "tileIcon": "ok"               // le nom de l'icône de la tuile : pause, maintenance, unmonitored, ok, minor, major ou critical
        "tileSecondaryIcon": null,     // le nom de l'icône secondaire de la tuile : pause, maintenance ou unmonitored
    },
    // ...
]
```

### Récupération des dépendances d'un watcher

**URL** : `/api/v2/weather/watchers/<watcher_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer le watcher `watcher_id` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/weather/watchers/watcher_id'
```

### Réponse en cas de réussite

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
[
    {
        "entity_id": "...",
        "name": "Nom de la dépendance",
        "source_type": "resource",
        "infos": {
            "customer":{
                "description": "Nom du client",
                "value": "Capensis"
            },
            // ...
        }

        "stats": {  // Statistiques sur l'entité (si le moteur go engine-stat est activé)
            "ko":0,
            "ok":0,
            "last_ko":1569919009,
            "last_event":1569923420
        },

        // Derniers steps state, status, snooze et ack de l'alarme de la dépendance
        "state": {...},   // contient {"val": 0} s'il n'y a pas d'alarme
        "status": {...},  // non défini s'il n'y a pas d'alarme
        "snooze": {...},  // non défini s'il n'y a pas d'alarme
        "ack": {...},     // non défini s'il n'y a pas d'alarme
        "ticket": {...},  // non défini s'il n'y a pas d'alarme

        // Champs de l'alarme (ces champs ne sont pas définis s'il n'y a pas d'alarme)
        "connector": "...",
        "connector_name": "...",
        "component": "...",
        "resource": "...",
        "last_update_date": ...,
        "alarm_creation_date": ...,
        "alarm_display_name": "...",
        "automatic_action_timer": ...,

        "pbehavior":[],  // Liste des pbehaviors actifs sur l'entité
        "linklist":[],   // Liste de lien générés par les linkbuilders
    },
    // ...
]
```
