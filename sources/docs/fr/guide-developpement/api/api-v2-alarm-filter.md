# Alarm-Filter

L'API alarm-filter permet de consulter, créer, modifier et supprimer des actions d'alarm-filter.

Pour plus d'informations sur les actions, consulter la [documentation de l'alarm-filter](../../guide-administration/moteurs/moteur-alerts-alarm-filter.md).

### Création d'actions

Crée une nouvelle action à partir du corps de la requête.

**URL** : `/api/v2/alerts/filters`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
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

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
{
    "entity_filter": "{\"v.state.val\": 1}",
    "condition": {},
    "tasks": [
        "alerts.useraction.state_increase"
    ],
    "output_format": "L'"'"'alarme a été créée il y a une heure. Augmentation automatique de son état.",
    "limit": 3600,
    "postpone_if_active_pbehavior": true
}' 'http://<Canopsis_URL>/api/v2/alerts/filters'
```

#### Réponse en cas de réussite

**Condition** : l'action est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "_id": "ad8457d0-d593-492e-bb41-c701d396b0b2",
    "entity_filter": "{\"v.state.val\": 1}",
    "condition": {},
    "tasks": [
        "alerts.useraction.state_increase"
    ],
    "output_format": "L'"'"'alarme a été créée il y a une heure. Augmentation automatique de son état.",
    "limit": 3600,
    "postpone_if_active_pbehavior": true,
    "repeat": 1
}
```

#### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas valide.

**Code** : `500 Internal Server Error`

**Corps de la réponse** : Page HTML. Une trace d'erreur est disponible dans les logs du webserver.

### Modification d'actions

Modifie une action à partir du corps de la requête.

**URL** : `/api/v2/alerts/filters/<action_id>`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
    "_id": "ad8457d0-d593-492e-bb41-c701d396b0b2",
    "entity_filter": "{\"v.state.val\": 2}",
    "condition": {},
    "tasks": [
        "alerts.useraction.state_increase"
    ],
    "output_format": "L'alarme a été créée il y a une heure. Augmentation automatique de son état.",
    "limit": 3600,
    "postpone_if_active_pbehavior": true
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
{
    "_id": "ad8457d0-d593-492e-bb41-c701d396b0b2",
    "entity_filter": "{\"v.state.val\": 2}",
    "condition": {},
    "tasks": [
        "alerts.useraction.state_increase"
    ],
    "output_format": "L'"'"'alarme a été créée il y a une heure. Augmentation automatique de son état.",
    "limit": 3600,
    "postpone_if_active_pbehavior": true
}' 'http://<Canopsis_URL>/api/v2/alerts/filters/ad8457d0-d593-492e-bb41-c701d396b0b2'
```

#### Réponse en cas de réussite

**Condition** : l'action est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "_id": "ad8457d0-d593-492e-bb41-c701d396b0b2",
    "entity_filter": "{\"v.state.val\": 2}",
    "condition": {},
    "tasks": [
        "alerts.useraction.state_increase"
    ],
    "output_format": "L'"'"'alarme a été créée il y a une heure. Augmentation automatique de son état.",
    "limit": 3600,
    "postpone_if_active_pbehavior": true,
    "repeat": 1
}
```

#### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas valide.

**Code** : `500 Internal Server Error`

**Corps de la réponse** : Page HTML. Une trace d'erreur est disponible dans les logs du webserver.

### Suppression d'actions

Supprime une action en fonction de son `_id`.

**URL** : `/api/v2/alerts/filters/<action_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer l'action avec l'`_id` `ad8457d0-d593-492e-bb41-c701d396b0b2` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/alerts/filters/ad8457d0-d593-492e-bb41-c701d396b0b2'
```

#### Réponse en cas de réussite

**Condition** : La suppression de l'action a réussi.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "n": 1,
    "ok": 1.0
}
```

### Récupération d'une action par id

**URL** : `/api/v2/alerts/filters/<action_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer l'action avec l'`id` `ad8457d0-d593-492e-bb41-c701d396b0b2` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/alerts/filters/ad8457d0-d593-492e-bb41-c701d396b0b2'
```

#### Réponse en cas de réussite

**Condition** : Une action correspondant à l'`id` est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
        "_id": "ad8457d0-d593-492e-bb41-c701d396b0b2",
        "entity_filter": "{\"v.state.val\": 2}",
        "condition": {},
        "tasks": [
            "alerts.useraction.state_increase"
        ],
        "output_format": "L'"'"'alarme a été créée il y a une heure. Augmentation automatique de son état.",
        "limit": 3600,
        "postpone_if_active_pbehavior": true,
        "repeat": 1
    }
]
```

#### Réponse en cas d'erreur

**Condition** : Aucune action ne correspondant à l'`id` est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[]
```
