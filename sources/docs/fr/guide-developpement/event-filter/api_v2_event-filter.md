# API Event-Filter

L'API Event-Filter permet de consulter, créer et supprimer des règles d'enrichissement.

Pour plus d'informations sur ce qu'est une règle d'enrichissement, consulter la [documentation du moteur Che](../../guide-administration/moteurs/event-filter/index.md).

### Création d'une règle

Crée une nouvelle règle à partir du corps de la requête.

**URL** : `/api/v2/eventfilter/rules`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
    "type": "enrichment",
    "pattern": {
        "component": "192.168.0.1"
    },
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "example.com"
        }
    ],
    "priority": 101,
    "on_success": "pass",
    "on_failure": "pass"
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
    "type": "enrichment",
    "pattern": {
        "component": "192.168.0.1"
    },
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "example.com"
        }
    ],
    "priority": 101,
    "on_success": "pass",
    "on_failure": "pass"
}' 'http://<Canopsis_URL>/api/v2/eventfilter/rules'
```

#### Réponse en cas de réussite

**Condition** : la règle est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "_id": "6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd"
}
```

#### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas valide.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
	"name": "",
	"description": "Malformed JSON: Extra data: line 1 column 66 - line 1 column 275 (char 65 - 274)"
}
```

---

**Condition** : Si une règle portant un `_id` similaire existe déjà en base.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
	"name": "",
	"description": "Operation failure while doing insert: E11000 duplicate key error collection: canopsis.eventfilter index: _id_ dup key: { : \"6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd\" }"
}
```

---

### Modification de Règle

Modifie une règle à partir du corps de la requête.

**URL** : `/api/v2/eventfilter/rules/<rule_id>

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
    "type": "enrichment",
    "pattern": {
        "component": "192.168.0.8"
    },
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "example.net"
        }
    ],
    "priority": 101,
    "on_success": "pass",
    "on_failure": "pass"
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut envoyer le JSON ci-dessus pour modifier la règle dont l'`_id` vaut `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
    "type": "enrichment",
    "pattern": {
        "component": "192.168.0.8"
    },
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "example.net"
        }
    ],
    "priority": 101,
    "on_success": "pass",
    "on_failure": "pass"
}' 'http://<Canopsis_URL>/api/v2/eventfilter/rules/6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd'
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


### Suppression de Règle

Supprime une règle en fonction de son `id`.

**URL** : `/api/v2/eventfilter/rules/<rule_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer la règle avec l'`id` `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/eventfilter/rules/6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd'
```

#### Réponse en cas de réussite

**Condition** : La suppression de la règle a réussi.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
}
```

#### Réponse en cas d'erreur

**Condition** : En cas d'absence de règle avec l'`_id` dans la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
	"name": "",
	"description": "No rule with id: 6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd"
}
```


### Récupération des règles

Récupère une ou plusieurs règles crée en base.

#### Récupération d'une règle par id

**URL** : `/api/v2/eventfilter/rules/<rule_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer la règle avec l'`id` `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/eventfilter/rules/6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd'
```

##### Réponse en cas de réussite

**Condition** : Une règle correspondant à l'`id` est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"pattern": {
		"component": "192.168.0.1"
	},
	"actions": [{
		"type": "set_field",
		"name": "Component",
		"value": "example.com"
	}],
	"priority": 101,
	"on_failure": "pass",
	"_id": "6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd",
	"type": "enrichment",
	"on_success": "pass"
}
```

##### Réponse en cas d'erreur

**Condition** : Aucune règle correspondant à l'`id` n'est trouvée.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
	"name": "",
	"description": "No rule with id: 6dacc239-59e8-4ba9-b1d0-e9c08ab8ea444cd"
}
```

---

#### Récupération de toutes les règles en base de données

Récupère toutes les règles stockées en base

**URL** : `/api/v2/eventfilter/rules`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer toutes les règles :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/eventfilter/rules'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[{
	"pattern": {
		"component": "192.168.0.1"
	},
	"actions": [{
		"type": "set_field",
		"name": "Component",
		"value": "example.com"
	}],
	"priority": 101,
	"on_failure": "pass",
	"_id": "4da620cd-5883-4952-8eac-b61f597fa622",
	"type": "enrichment",
	"on_success": "pass"
}, {
	"pattern": {
		"component": "192.168.0.2"
	},
	"actions": [{
		"type": "set_field",
		"name": "Component",
		"value": "example2.com"
	}],
	"priority": 101,
	"on_failure": "pass",
	"_id": "f62b37c5-e301-4bfc-ba52-8ff2454cc0aa",
	"type": "enrichment",
	"on_success": "pass"
}]
```
