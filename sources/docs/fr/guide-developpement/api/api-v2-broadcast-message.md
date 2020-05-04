# Messages d'information

!!! note
    Disponible depuis Canopsis 3.40.0.

L'API broadcast-message permet de consulter, créer, modifier et supprimer des messages d'informations qui seront affichés sur l'interface graphique de Canopsis.

Pour plus d'informations sur ce qu'est un `message d'informations`, consulter la [documentation sur les messages](../../guide-utilisation/interface/broadcast-message.md).

### Création d'un message d'information

Crée un nouveau message d'information à partir du corps de la requête.

**URL** : `/api/v2/broadcast-message`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
	"message": "Maintenance en prévision",
	"color": "#e75e40",
	"start": 1588601154,
	"end": 1588601400
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
	"message": "Maintenance en prévision",
	"color": "#e75e40",
	"start": 1588601154,
	"end": 1588601400
}' 'http://<Canopsis_URL>/api/v2/broadcast-message'
```

#### Réponse en cas de réussite

**Condition** : le message d'information est créé

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"_id": "4be6eb22-f173-4ac1-8352-1f37cf5caf48"
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

### Modification de message d'information

Modifie un message d'information à partir du corps de la requête.

**URL** : `/api/v2/broadcast-message/<broadcast-message_id>`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
	"message": "Maintenance en prévision",
	"color": "#0062B1",
	"start": 1588601154,
	"end": 1588601400
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut envoyer le JSON ci-dessus pour modifier la règle dont l'`_id` vaut `3590174b-ea11-4d32-bdcd-72335b41b4fc` :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
	"message": "Maintenance en prévision",
	"color": "#0062B1",
	"start": 1588601154,
	"end": 1588601400
}' 'http://<Canopsis_URL>/api/v2/broadcast-message/3590174b-ea11-4d32-bdcd-72335b41b4fc'
```

#### Réponse en cas de réussite

**Condition** : le message d'information est modifié

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
}
```


### Suppression de message d'information

Supprime un message d'information en fonction de son `id`.

**URL** : `/api/v2/broadcast-message/<broadcast-message_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer le message d'information avec l'`id` `3590174b-ea11-4d32-bdcd-72335b41b4fc` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/broadcast-message/3590174b-ea11-4d32-bdcd-72335b41b4fc'
```

#### Réponse en cas de réussite

**Condition** : La suppression du message d'information a réussi.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "status": true
}
```


### Récupération des messages d'information

Récupère un ou plusieurs messages d'information créés en base.

#### Récupération d'un message d'information par id

**URL** : `/api/v2/broadcast-message/<broadcast-message_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer le message d'information avec l'`id` `4be6eb22-f173-4ac1-8352-1f37cf5caf48` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/broadcast-message/4be6eb22-f173-4ac1-8352-1f37cf5caf48'
```

##### Réponse en cas de réussite

**Condition** : Un message d'information correspondant à l'`id` est trouvé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "color": "#e75e40", 
    "start": 1588601154, 
    "message": "Maintenance en pr\u00e9vision", 
    "_id": "4be6eb22-f173-4ac1-8352-1f37cf5caf48", 
    "end": 1588601400
}
```

##### Réponse en cas d'erreur

**Condition** : Aucun message correspondant à l'`id` n'est trouvé.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "No message found with ID 4be6eb22-f173-4ac1-8352-1f37cf5caf4"}
}
```

#### Récupération de tous les messages d'information en base de données

Récupère tous les messages d'informations stockés en base

**URL** : `/api/v2/broadcast-message`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les messages d'information :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/broadcast-message'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
      "_id": "aa481acfb2d6d932c0654e5a23e20019",
      // ...
    },
    {
      "_id": "yet-another-service",
      // ...
    }
]
```

#### Récupération de tous les messages d'information actifs

Récupère tous les messages d'informations actifs

**URL** : `/api/v2/broadcast-message/active`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les messages d'information :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/broadcast-message/active'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
      "_id": "aa481acfb2d6d932c0654e5a23e20019",
      // ...
    },
    {
      "_id": "yet-another-service",
      // ...
    }
]
```
