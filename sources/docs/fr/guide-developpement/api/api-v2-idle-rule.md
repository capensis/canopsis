# API Idle Rule
L'API Idle Rule permet de consulter, créer, modifier et supprimer des règles de détection d'inactivité sur les alarmes.

### Création d'Idle rules
**URL** : `/api/v2/idle-rule`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
  "name" : "close 3m after last event",
  "type": "last_event",
  "duration": "3m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
   ]
}
```
**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :
```sh
curl -H "Content-Type: application/json" -u root:root -X POST -d '{
  "name" : "close 3m after last event",
  "type": "last_event",
  "duration": "3m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
   ]
}' 'http://localhost:8082/api/v2/idle-rule'
```

#### Réponse en cas de réussite

**Condition** : l'Idle rule est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "description": "",
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
  ],
  "last_modified_date": 1589980727,
  "creation_date": 1589980727,
  "duration": "3m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "name": "close 3m after last event",
  "author": "root",
  "entity_patterns": null,
  "_id": "1dd3f61a-3eb6-4fbd-829e-7d6cb3f1a4a9",
  "type": "last_event"
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

**Condition** : Si une Idle rule similaire existe déjà en base.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "failed to create idle rule: duplicate id idltest"
}
```

### Modification d'Idle rules

Modifie une idle rule à partir du corps de la requête.

**URL** : `/api/v2/idle-rule/<idle-rule_id>`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
  "name" : "close 3m after last event",
  "_id":"idltest",
  "type": "last_event",
  "duration": "150m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
   ]
}'
```
**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut envoyer le JSON ci-dessus pour modifier l'Idle rule dont l'`_id` vaut `idltest` :
```sh
curl -v -H "Content-Type: application/json" -u root:root -X PUT -d '{
  "name" : "close 3m after last event",
  "_id":"idltest",
  "type": "last_event",
  "duration": "150m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
   ]
}' 'http://localhost:8082/api/v2/idle-rule/idltest'
```

#### Réponse en cas de réussite

**Condition** : la règle est modifiée

**Code** : `200 OK`

**Exemple du corps de la réponse** :
```json 
{
  "description": "",
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
  ],
  "last_modified_date": 1589981972,
  "creation_date": 1589980983,
  "duration": "150m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "name": "close 3m after last event",
  "author": "root",
  "entity_patterns": null,
  "_id": "idltest",
  "type": "last_event"
}
```


### Suppression d'Idle rules

Supprime une Idle rule en fonction de son `id`.

**URL** : `/api/v2/idle-rule/<idle-rule_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer l'idle rule avec l'`id` `idltest` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/idle-rule/idltest'
```

#### Réponse en cas de réussite

**Condition** : La suppression de l'Idle rule a réussi.

**Exemple du corps de la réponse** :
```json
{
  "deletedCount": 1,
  "acknowledged": true
}
```

##### Réponse en cas d'erreur

**Condition** : Aucune Idle rule correspondant à l'`id` n'est trouvée.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "no idle rule with id e6d5add2-8e02-4c8a-bc0d-d1d2bf36b755"
}
```

### Récupération des Idle rules

#### Récupération d'une Idle rule par id

**URL** : `/api/v2/idle-rule/<idle-rule_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer l'idle rule avec l'`id` `idltest` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/idle-rule/idltest'
```

##### Réponse en cas de réussite

**Condition** : Une Idle rule correspondant à l'`id` est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :
```json
{
  "description": "",
  "author": "root",
  "alarm_patterns": [
    {
      "v": {
        "resource": {
          "regex_match": "CPU"
        }
      }
    }
  ],
  "last_modified_date": 1589981972,
  "creation_date": 1589980983,
  "entity_patterns": null,
  "duration": "150m",
  "operation": {
    "type": "cancel",
    "parameters": {
      "author": "idle rule"
    }
  },
  "_id": "idltest",
  "type": "last_event",
  "name": "close 3m after last event"
}
```

##### Réponse en cas d'erreur

**Condition** : Aucune Idle rule correspondant à l'`id` n'est trouvée.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
  "name": "idltest",
  "description": "Rule not found"
}
```

#### Récupération de toutes les Idle rules en base de données

Récupère toutes les Idle rules stockées en base.

**URL** : `/api/v2/idle-rule`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer toutes les Idle rules :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/idle-rule'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
  {
    "description": "",
    "author": "root",
    "alarm_patterns": [
      {
        "v": {
          "resource": {
            "regex_match": "CPU"
          }
        }
      }
    ],
    "last_modified_date": 1589981662,
    "creation_date": 1589981662,
    "entity_patterns": null,
    "duration": "3m",
    "operation": {
      "type": "cancel",
      "parameters": {
        "author": "idle rule"
      }
    },
    "_id": "4a8f6cfd-14ec-4bc8-b1cf-4be19b18fb29",
    "type": "last_event",
    "name": "close 3m after last event"
  },
  {
    "description": "",
    "author": "root",
    "alarm_patterns": [
      {
        "v": {
          "resource": {
            "regex_match": "CPU"
          }
        }
      }
    ],
    "last_modified_date": 1589982546,
    "creation_date": 1589982546,
    "entity_patterns": null,
    "duration": "3m",
    "operation": {
      "type": "cancel",
      "parameters": {
        "author": "idle rule"
      }
    },
    "_id": "idle-rule-delete-test",
    "type": "last_event",
    "name": "close 3m after last event"
  }
]
```