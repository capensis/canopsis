# Heartbeat

L'API HeartBeat permet de consulter, créer et supprimer des HeartBeats.

Pour plus d'informations sur ce qu'est un heartbeat, consulter la [documentation du moteur `engine-heartbeat`](../../guide-administration/moteurs/moteur-heartbeat.md).

### Création de Heartbeat

Crée un nouveau HeartBeat à partir du corps de la requête.

**URL** : `/api/v2/heartbeat`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de corps de requête** :
```json
{
  "author": "créateur",
  "name": "Heartbeat",
  "description": "Ceci est une description",
  "expected_interval": "1m",
  "output": "Une note personnalisable",
  "pattern": {
    "connector": "connect",
    "connector_name": "c01"
  }
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
  "author": "créateur",
  "name": "Heartbeat",
  "description": "Ceci est une description",
  "expected_interval": "1m",
  "output": "Une note personnalisable",
  "pattern": {
    "connector": "connect",
    "connector_name": "c01"
  }
}' 'http://localhost:8082/api/v2/heartbeat'
```

#### Réponse en cas de réussite

**Condition** : l'HeartBeat est créé

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "name": "heartbeat created",
  "description": "eaadc64a-a0ef-4c5c-8453-a45e574fd5ca"
}
```

#### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas valide.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "invalid heartbeat payload."
}
```

---

**Condition** : Si un HeartBeat similaire existe déjà en base.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "heartbeat pattern already exists"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "database error, please contact your administrator."
}
```


### Suppression de Heartbeat

Supprime un HeartBeat en fonction de son `id`.

**URL** : `/api/v2/heartbeat/<heartbeat_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer le heartbeat avec l'`id` `eaadc64a-a0ef-4c5c-8453-a45e574fd5ca` :

```sh
curl -X DELETE -u root:root 'http://localhost:8082/api/v2/heartbeat/eaadc64a-a0ef-4c5c-8453-a45e574fd5ca'
```

#### Réponse en cas de réussite

**Condition** : L'Heartbeat à bien été supprimé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "name": "heartbeat removed",
  "description": "eaadc64a-a0ef-4c5c-8453-a45e574fd5ca"
}
```

#### Réponse en cas d'erreur

**Condition** : Si l'`id` ne correspond à aucun HeartBeat.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
  "description" : "heartbeat not found",
  "name" : "eaadc64a-a0ef-4c5c-8453-a45e574fd5ca"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "database error, please contact your administrator."
}
```


### Récupération des HeartBeats

Récupère un ou plusieurs HeartBeats créés en base.

#### Récupération d'un HeartBeat par id

**URL** : `/api/v2/heartbeat/<heartbeat_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer le heartbeat avec l'`id` `cf7097b6-71ba-48db-a749-e9474aa70c93` :

```sh
curl -X GET -u root:root 'http://localhost:8082/api/v2/heartbeat/cf7097b6-71ba-48db-a749-e9474aa70c93'
```

##### Réponse en cas de réussite

**Condition** : Un HeartBeat correspondant à l'`id` est trouvé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "updated": 1607680814,
  "description": "Ceci est une description",
  "created": 1607680814,
  "pattern": {
    "connector": "connect",
    "connector_name": "c01"
  },
  "author": "créateur",
  "expected_interval": "1m",
  "output": "Une chaîne de caractères",
  "_id": "cf7097b6-71ba-48db-a749-e9474aa70c93",
  "name": "Heartbeat"
}
```

##### Réponse en cas d'erreur

**Condition** : Aucun HeartBeat correspondant à l'`id` n'est trouvé.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
  "description" : "heartbeat not found",
  "name" : "cf7097b6-71ba-48db-a749-e9474aa70c93"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "database error, please contact your administrator."
}
```

#### Récupération de tous les HeartBeats en base de données

Récupère tous les HeartBeats stockés en base

**URL** : `/api/v2/heartbeat`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les heartbeats :

```sh
curl -X GET -u root:root 'http://localhost:8082/api/v2/heartbeat'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
      "updated": 1607680882,
      "description": "Ceci est une description",
      "created": 1607680882,
      "pattern": {
        "connector": "connect",
        "connector_name": "c02"
      },
      "author": "root",
      "expected_interval": "1m",
      "output": "Une chaîne de caractères",
      "_id": "619e1f78-cb2e-4f15-afbe-1921aef8db8e",
      "name": "Heartbeat 2"
    },
    {
      "updated": 1607680814,
      "description": "Ceci est une description",
      "created": 1607680814,
      "pattern": {
        "connector": "connect",
        "connector_name": "c01"
      },
      "author": "créateur",
      "expected_interval": "1m",
      "output": "Une chaîne de caractères",
      "_id": "cf7097b6-71ba-48db-a749-e9474aa70c93",
      "name": "Heartbeat"
    }
]
```

##### Réponse en cas d'erreur

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "database error, please contact your administrator."
}
```
