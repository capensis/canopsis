# API Watchers

L'API Watchers permet de consulter, créer et supprimer des Watchers.

### Creation de watcher

Crée un nouveau watcher à partir du corps de la requête.

**URL** : `/api/v2/watcherng`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de corps de requête** :
```json
{
    "_id": "h4z25rzg6rt-64rge354-5re4g",
    "name": "Client Capensis",
    "type": "watcher",
    "entities": [{
        "infos": {
            "customer": {
                "value": "capensis"
            }
        }
    }, {
        "_id": {"regex_match": ".+/comp"}
    }],
    "state": {
        "method": "worst"
    },
    "output_template": "Alarmes critiques : {{.State.Critical}}"
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le Json ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{

}' 'http://<Canopsis_URL>/api/v2/watcherng'
```

#### Réponse en cas de réussite

**Condition** : le watcher est créé

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
"h4z25rzg6rt-64rge354-5re4g"
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

**Condition** : Si un watcher similaire existe déjà en base.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "Error while creating a watcher"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "Error while creating a watcher"
}
```


### Suppression de watcher

Supprime un watcher en fonction de son `id`.

**URL** : `/api/v2/watcherng/<watcher>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer le watcher avec l'`id` `h4z25rzg6rt-64rge354-5re4g` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/watcherng/h4z25rzg6rt-64rge354-5re4g'
```

#### Réponse en cas de réussite

**Condition** : La suppresion du watcher a réussi.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "status": true
}
```

#### Réponse en cas d'erreur

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "Can not retrieve the watcher data from database, contact your administrator."
}
```


### Récupération des watchers

Récupère un ou plusieurs watchers présent en base.

#### Récupération d'un watcher par id

**URL** : `/api/v2/watcherng/<watcher_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer le watcher avec l'`id` `h4z25rzg6rt-64rge354-5re4g` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/watcherng/h4z25rzg6rt-64rge354-5re4g'
```

##### Réponse en cas de réussite

**Condition** : Un watcher correspondant à l'`id` est trouvé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "_id": "h4z25rzg6rt-64rge354-5re4g",
    "name": "Client Capensis",
    "type": "watcher",
    "entities": [{
        "infos": {
            "customer": {
                "value": "capensis"
            }
        }
    }, {
        "_id": {"regex_match": ".+/comp"}
    }],
    "state": {
        "method": "worst"
    },
    "output_template": "Alarmes critiques : {{.State.Critical}}"
}
```

##### Réponse en cas d'erreur

**Condition** : Aucun watcher correspondant à l'`id` n'est trouvé.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "No watcher found with ID declare_external_ticket"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "Can not retrieve the watcher data from database, contact your administrator."
}
```

#### Récupération de tous les watchers en base de données.

Récupère tous les watchers stocké en base

**URL** : `/api/v2/watcherng`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les webhooks :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/watcherng'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
      "_id": "h4z25rzg6rt-64rge354-5re4g",
      // ...
    },
    {
      "_id": "aa481acfb2d6d932c0654e5a23e20019",
      // ...
    },
    {
      "_id": "yet-another-watcher",
      // ...
    }
]
```

##### Réponse en cas d'erreur

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "Can not retrieve the watchers list from database, contact your administrator."
}
```
