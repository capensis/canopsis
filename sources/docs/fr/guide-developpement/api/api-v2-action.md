# Action

L'API Action permet de consulter, créer, modifier et supprimer des Actions sur les alarmes.

Pour plus d'informations sur ce qu'est une action, consulter la [documentation du moteur Action](../../guide-administration/moteurs/moteur-action.md).

### Création d'Actions

Crée une nouvelle Action à partir du corps de la requête.

**URL** : `/api/v2/actions`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
    "_id": "action_id_pbehavior",
    "type": "pbehavior",
    "hook": {
        "event_patterns": [
            {
                "resource": "CPU"
            },
            {
                "resource": "HDD"
            }
        ],
        "triggers": [
            "create"
        ]
    },
    "parameters": {
        "name": "pbehavior_name",
        "author": "System",
        "type": "Pause",
        "rrule": "",
        "reason": "Problème d\'habilitation",
        "tstart": 0,
        "tstop": 253402297199
    }
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
    "_id": "action_id_pbehavior",
    "type": "pbehavior",
    "hook": {
        "event_patterns": [
            {
                "resource": "CPU"
            },
            {
                "resource": "HDD"
            }
        ],
        "triggers": [
            "create"
        ]
    },
    "parameters": {
        "name": "pbehavior_test_action",
        "author": "System",
        "type": "Pause",
        "rrule": "",
        "reason": "Problème d\'habilitation",
        "tstart": 0,
        "tstop": 4170912120
    }
}' 'http://<Canopsis_URL>/api/v2/actions'
```

#### Réponse en cas de réussite

**Condition** : l'Action est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"_id": "action_id_pbehavior"
}
```

#### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas valide.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "Invalid JSON"
}
```

---

**Condition** : Si une action similaire existe déjà en base.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "error while creating an action"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "error while creating an action"
}
```

### Modification d'Actions

Modifie une action à partir du corps de la requête.

**URL** : `/api/v2/actions/<action_id>`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
    "_id": "action_id",
    "type": "pbehavior",
    "hook": {
        "event_patterns": [
            {
                "resource": "RAM"
            },
            {
                "resource": "SWAP"
            }
        ],
        "triggers": [
            "create"
        ]
    },
    "parameters": {
        "name": "pbehavior_name",
        "author": "System",
        "type": "Pause",
        "rrule": "",
        "reason": "",
        "tstart": 0,
        "tstop": 253402297199
    }
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut envoyer le JSON ci-dessus pour modifier l'Action dont l'`_id` vaut `action_id_pbehavior` :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
    "_id": "action_id",
    "type": "pbehavior",
    "hook": {
        "event_patterns": [
            {
                "resource": "RAM"
            },
            {
                "resource": "SWAP"
            }
        ],
        "triggers": [
            "create"
        ]
    },
    "parameters": {
        "name": "pbehavior_name",
        "author": "System",
        "type": "Pause",
        "rrule": "",
        "reason": "",
        "tstart": 0,
        "tstop": 253402297199
    }
}' 'http://<Canopsis_URL>/api/v2/actions/action_id_pbehavior'
```

#### Réponse en cas de réussite

**Condition** : la règle est modifiée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
}
```

---

### Suppression d'Actions

Supprime une Action en fonction de son `id`.

**URL** : `/api/v2/actions/<action_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer l'action avec l'`id` `action_id_pbehavior` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/actions/action_id_pbehavior'
```

#### Réponse en cas de réussite

**Condition** : La suppression de l'Action a réussi.

Renvoie un booléen.

### Récupération des Actions

Récupère une ou plusieurs Actions créées en base.

#### Récupération d'une Action par id

**URL** : `/api/v2/actions/<action_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer l'action avec l'`id` `action_id_pbehavior` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/actions/action_id_pbehavior'
```

##### Réponse en cas de réussite

**Condition** : Une Action correspondant à l'`id` est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"regex": null,
	"parameters": {
		"name": "pbehavior_test_action",
		"author": "System",
		"reason": "Probl\u00e8me d'habilitation",
		"rrule": "",
		"tstart": 0,
		"tstop": 4170912120,
		"type": "Pause"
	},
	"fields": [],
	"hook": {
		"event_patterns": [{
			"resource": "CPU"
		}, {
			"resource": "HDD"
		}],
		"triggers": ["create"]
	},
	"_id": "action_id_pbehavior",
	"type": "pbehavior"
}
```

##### Réponse en cas d'erreur

**Condition** : Aucune Action correspondant à l'`id` n'est trouvée.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "failed to get action"
}
```

---

#### Récupération de toutes les Actions en base de données

Récupère toutes les Actions stockées en base

**URL** : `/api/v2/actions`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer toutes les Actions :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/actions'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[{
	"hook": {
		"event_patterns": [{
			"resource": "CPU"
		}, {
			"resource": "HDD"
		}],
		"triggers": ["create"]
	},
	"_id": "action_id_pbehavior",
	"type": "pbehavior",
	"parameters": {
		"name": "pbehavior_test_action",
		"author": "System",
		"reason": "Probl\u00e8me d'habilitation",
		"rrule": "",
		"tstart": 0,
		"tstop": 4170912120,
		"type": "Pause"
	}
}, {
	"hook": {
		"event_patterns": [{
			"resource": "CPU_2"
		}, {
			"resource": "HDD_2"
		}],
		"triggers": ["create"]
	},
	"_id": "action_test_snooze_100m",
	"type": "snooze",
	"parameters": {
		"duration": 6000,
		"message": null,
		"author": "action"
	}
}]
```
