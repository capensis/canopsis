# API Heartbeat

L'API HeartBeat permet de consulter, créer et supprimer des HeartBeats.

### Creation de Heartbeat

Crée un nouveau HeartBeat à partir du corps de la requête.

**URL** : `/api/v2/heartbeat`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de corps de requête** :
```json
{
  "pattern": {
      "connector": "c1",
      "connector_name": "connector1"
  },
  "expected_interval": "10s"
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le Json ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
  "pattern": {
      "connector": "c1",
      "connector_name": "connector1"
  },
  "expected_interval": "10s"
}' 'http://<Canopsis_URL>/api/v2/heartbeat'
```

#### Réponse en cas de réussite

**Condition** : l'HeartBeat est crée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "name": "heartbeat created",
  "description": "cd92421e77f48435d38b3682beb62f07"
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

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer le heartbeat avec l'`id` `cd92421e77f48435d38b3682beb62f07` :

```sh
curl -X DEL -u root:root 'http://<Canopsis_URL>/api/v2/heartbeat/cd92421e77f48435d38b3682beb62f07'
```

#### Réponse en cas de réussite

**Condition** : L'Heartbeat à bien été supprimé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "name": "heartbeat removed",
  "description": "cd92421e77f48435d38b3682beb62f07"
}
```

#### Réponse en cas d'erreur

**Condition** : Si l'`id` ne correspond à aucun HeartBeat.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
  "description" : "heartbeat not found",
  "name" : "cd92421e77f48435d38b3682beb62f07"
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

Récupère un ou plusieurs HeartBeat crée en base.

#### Récupération d'un HeartBeat par id

**URL** : `/api/v2/heartbeat/<heartbeat_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer le heartbeat avec l'`id` `cd92421e77f48435d38b3682beb62f07` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/heartbeat/cd92421e77f48435d38b3682beb62f07'
```

##### Réponse en cas de réussite

**Condition** : Un HeartBeat correspondant à l'`id est trouvé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "_id": "cd92421e77f48435d38b3682beb62f07",
  "pattern": {
      "connector": "c1",
      "connector_name": "connector1"
  },
  "expected_interval": "10s"
}
```

##### Réponse en cas d'erreur

**Condition** : Aucun HeartBeat correspondant à l'`id n'est trouvé.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
  "description" : "heartbeat not found",
  "name" : "cd92421e77f48435d38b3682beb62f07"
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

#### Récupération de tous les HeartBeat en base de données.

Récupèr tous les HeartBeats stocké en base

**URL** : `/api/v2/heartbeat/`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les heartbeats :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/heartbeat/'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
      "_id": "cd92421e77f48435d38b3682beb62f07",
      "pattern": {
          "connector": "c1",
          "connector_name": "connector1"
      },
      "expected_interval": "10s"
    },
    {
      "_id": "3d071cb49acce44040b35ed8c6714ef1",
      "pattern": {
          "connector": "c2",
          "connector_name": "connector2"
      },
      "expected_interval": "15m"
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
