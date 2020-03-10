# PBehavior

L'API PBehavior permet de consulter, créer et supprimer des plages de maintenance.

Pour plus d'informations sur ce qu'est un PBehavior, consulter la [documentation du moteur PBehavior](../../guide-administration/moteurs/moteur-pbehavior.md).

### Création d'un pbehavior

Crée un nouveau PBehavior à partir du corps de la requête.

**URL** : `/api/v2/pbehavior`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
	"author": "root",
	"name": "Pbehavior test 1",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"filter": {
		"$or": [{
			"impact": {
				"$in": ["pbehavior_test_1"]
			}
		}, {
			"$and": [{
				"type": "component"
			}, {
				"name": "pbehavior_test_1"
			}]
		}]
	},
	"type_": "Hors plage horaire de surveillance",
	"reason": "Autre",
	"rrule": null,
	"comments": [],
	"exdate": []
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
	"author": "root",
	"name": "Pbehavior test 1",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"filter": {
		"$or": [{
			"impact": {
				"$in": ["pbehavior_test_1"]
			}
		}, {
			"$and": [{
				"type": "component"
			}, {
				"name": "pbehavior_test_1"
			}]
		}]
	},
	"type_": "Hors plage horaire de surveillance",
	"reason": "Autre",
	"rrule": null,
	"comments": [],
	"exdate": []
}' 'http://<Canopsis_URL>/api/v2/pbehavior'
```

#### Réponse en cas de réussite

**Condition** : le pbehavior est créé

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"e89d0a8f-8ccd-4357-83e9-ea3f8a53ebb2"
}
```

#### Cas particulier : Permettre la création d'un pbehavior avec un `_id` déjà existant

Lorsque les pbehaviors sont pilotés par l'API (via un ordonnanceur par exemple), il est parfois nécessaire de devoir 
"écraser" un pbehavior déjà existant mais qui serait expiré.  
Un paramètre complémentaire peut être passé à l'API dans ce cas : `replace_expired`.  


**URL** : `/api/v2/pbehavior?replace_expired=1`

**Exemple de corps de requête** :
```json
{
	"_id": "pbh1",
	"author": "root",
	"name": "Pbehavior test 1",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"filter": {
		"$or": [{
			"impact": {
				"$in": ["pbehavior_test_1"]
			}
		}, {
			"$and": [{
				"type": "component"
			}, {
				"name": "pbehavior_test_1"
			}]
		}]
	},
	"type_": "Hors plage horaire de surveillance",
	"reason": "Autre",
	"rrule": null,
	"comments": [],
	"exdate": []
}
```

Dans ce cas, si le pbehavior dont l'\_id vaut "pbh1" est expiré, alors

* ce pbehavior est renommé avec le format : `EXP{current_timestamp_in_milliseconds}-{_id}`
* un nouveau pbehavior avec les caractéristiques données en paramètre est créé

---

### Modification de PBehavior

Modifie un pbehavior à partir du corps de la requête.

Les champs `eids` et `comments` des pBehaviors ne sont pas modifiables avec cette route. Ils ne sont pas pris en compte s'ils sont présents dans le corps de la requête.

**URL** : `/api/v2/pbehavior/<pbehavior_id>`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
	"author": "root",
	"name": "Pbehavior test 1",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"filter": {
		"$or": [{
			"impact": {
				"$in": ["pbehavior_test_1"]
			}
		}, {
			"$and": [{
				"type": "component"
			}, {
				"name": "pbehavior_test_1"
			}]
		}]
	},
	"type_": "Hors plage horaire de surveillance",
	"reason": "Autre",
	"rrule": "FREQ=WEEKLY;BYDAY=FR,TH",
	"comments": [],
	"exdate": []
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut envoyer le JSON ci-dessus pour modifier le pbehavior dont l'`_id` vaut `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
	"author": "root",
	"name": "Pbehavior test 1",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"filter": {
		"$or": [{
			"impact": {
				"$in": ["pbehavior_test_1"]
			}
		}, {
			"$and": [{
				"type": "component"
			}, {
				"name": "pbehavior_test_1"
			}]
		}]
	},
	"type_": "Hors plage horaire de surveillance",
	"reason": "Autre",
	"rrule": "FREQ=WEEKLY;BYDAY=FR,TH",
	"comments": [],
	"exdate": []
}' 'http://<Canopsis_URL>/api/v2/pbehavior/e89d0a8f-8ccd-4357-83e9-ea3f8a53ebb2'
```

#### Réponse en cas de réussite

**Condition** : le pbehavior est modifié

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"name": "PBehavior de test",
	"author": "root",
	"reason": "Autre",
	"filter": "{\"$or\": [{\"impact\": {\"$in\": [\"pbehavior_test_1\"]}}, {\"$and\": [{\"type\": \"pbehavior_test_1\"}, {\"name\": \"pbehavior_test_1\"}]}]}",
	"type_": "Hors plage horaire de surveillance",
	"exdate": [],
	"tstart": 1567429156,
	"tstop": 1567601940
}
```

---

### Suppression de PBehavior

Supprime un PBehavior en fonction de son `id`.

**URL** : `/api/v2/pbehavior/<pbehavior_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer le PBehavior avec l'`id` `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/pbehavior/e89d0a8f-8ccd-4357-83e9-ea3f8a53ebb2'
```

#### Réponse en cas de réussite

**Condition** : La suppression du pbehavior a réussi.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"deletedCount": 1,
	"acknowledged": true
}
```

---

### Récupération des pbehaviors

Récupère un ou plusieurs pbehaviors appliqués sur une entité, via l'`eid` (pour `entity id`).

#### Récupération d'un pbehavior par eid

**URL** : `/api/v2/pbehavior_byeid/<entity_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer les pbehaviors s'appliquant à l'entité dont l'`_id` est `disk2/pbehavior_test_1` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/pbehavior_byeid/disk2/pbehavior_test_1'
```

##### Réponse en cas de réussite

**Condition** : Au moins un pbehavior appliqué à une entité correspondant à l'`id` est trouvé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[{
	"name": "Pbehavior test 1",
	"author": "root",
	"eids": ["pbehavior_test_1", "disk2/pbehavior_test_1"],
	"reason": "Autre",
	"filter": "{\"$or\": [{\"impact\": {\"$in\": [\"pbehavior_test_1\"]}}, {\"$and\": [{\"type\": \"component\"}, {\"name\": \"pbehavior_test_1\"}]}]}",
	"type_": "Hors plage horaire de surveillance",
	"rrule": "FREQ=WEEKLY;BYDAY=FR,TH",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"_id": "4c441d4e-9cc8-4f84-be73-9a4e97ba5e74",
	"isActive": true,
	"exdate": []
}]
```

---

#### Récupération de tous les pbehaviors en base de données

Récupère tous les pbehaviors stockés en base

**URL** : `/pbehavior/read`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Supporte la pagination** : Oui

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les pbehaviors :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/pbehavior/read'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"total": 1,
	"data": [{
		"count": 1,
		"total_count": 1,
		"data": [{
			"connector": "canopsis",
			"name": "Pbehavior test 1",
			"author": "root",
			"enabled": true,
			"reason": "Autre",
			"comments": [],
			"filter": "{\"$or\": [{\"impact\": {\"$in\": [\"pbehavior_test_1\"]}}, {\"$and\": [{\"type\": \"component\"}, {\"name\": \"pbehavior_test_1\"}]}]}",
			"type_": "Hors plage horaire de surveillance",
			"connector_name": "canopsis",
			"eids": [],
			"tstart": 1567439123,
			"timezone": "Europe/Paris",
			"tstop": 1569599100,
			"_id": "aaa9d5c3-b245-481f-b23a-844893cb3cfe",
			"rrule": "FREQ=WEEKLY;BYDAY=FR,TH",
			"exdate": []
		}]
	}],
	"success": true
}
```
