# Meta alarm rule 

L'API Meta-alarm-rule permet de consulter, créer et supprimer des règles de corrélation.

Pour plus d'informations sur ce qu'est une règle de corrélation, consulter la [documentation du moteur `engine-correlation`](../../guide-administration/moteurs/moteur-correlation.md).

### Création d'une règle

Crée une nouvelle règle à partir du corps de la requête.

**URL** : `/api/v2/metaalarmrule`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
  "_id" : "metarule1",
  "name": "Test groupement de type complexe",
  "type": "complex",
  "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
  "config": {
    "time_interval": 10,
    "threshold_count": 3,
    "alarm_patterns": [
      {
        "v": {
          "resource": {
            "regex_match" : "meta_complex"
          },
          "state": {
            "val": 3
          }
        }
      }
    ]
  }
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
  "_id" : "metarule1",
  "name": "Test groupement de type complexe",
  "type": "complex",
  "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
  "config": {
    "time_interval": 10,
    "threshold_count": 3,
    "alarm_patterns": [
      {
        "v": {
          "resource": {
            "regex_match" : "meta_complex"
          },
          "state": {
            "val": 3
          }
        }
      }
    ]
  }
}' 'http://localhost:8082/api/v2/metaalarmrule'
```

#### Réponse en cas de réussite

**Condition** : la règle est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```
metarule1
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
	"description": "Trying to insert MetaAlarmRule with already existing _id"
}
```

---

**Condition** : Si le type de la règle est invalide (simple)

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
	"name": "", 
	"description": "rule type invalid value simple"
}
```

---

**Condition** : Si l'attribut `config.threshold_rate` est utilisé avec une règle de type `attribute` :

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
	"name": "", 
	"description": "invalid rule_type attribute with config time_interval"
}
```

---

### Suppression de Règle

Supprime une règle en fonction de son `id`.

**URL** : `/api/v2/metaalarmrule/<rule_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer la règle avec l'`id` `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X DELETE -u root:root 'http://localhost:8082/api/v2/metaalarmrule/6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd'
```

#### Réponse en cas de réussite

**Condition** : La suppression de la règle a réussi.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "acknowledged": true,
    "deletedCount": 1
}
```

#### Réponse en cas d'erreur

**Condition** : En cas d'absence de règle avec l'`_id` dans la base de données (`deletedCount` vaut 0).

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "acknowledged": true,
    "deletedCount": 0
}
```


### Récupération des règles

Récupère une ou plusieurs règles créées en base.

#### Récupération d'une règle par id

**URL** : `/api/v2/metaalarmrule/<rule_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer la règle avec l'`id` `6dacc239-59e8-4ba9-b1d0-e9c08ab8eacd` :

```sh
curl -X GET -u root:root 'http://localhost:8082/api/v2/metaalarmrule/8c7b1732-522f-4bcf-a7ac-d08bd9c085eb'
```

##### Réponse en cas de réussite

**Condition** : Une règle correspondant à l'`id` est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
  "patterns": null,
  "_id": "8c7b1732-522f-4bcf-a7ac-d08bd9c085eb",
  "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
  "config": {
    "time_interval": 10,
    "alarm_patterns": [
      {
        "v": {
          "state": {
            "val": 3
          },
          "resource": {
            "regex_match": "meta_complex"
          }
        }
      }
    ],
    "threshold_count": 3
  },
  "name": "Test groupement de type complexe",
  "type": "complex"
}

```

##### Réponse en cas d'erreur

**Condition** : Aucune règle correspondant à l'`id` n'est trouvée.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```
null
```

---

#### Récupération de toutes les règles en base de données

Récupère toutes les règles stockées en base

**URL** : `/api/v2/metaalarmrule`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer toutes les règles :

```sh
curl -X GET -u root:root 'http://localhost:8082/api/v2/metaalarmrule'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
  {
    "patterns": null,
    "_id": "73da1ad7-058e-46af-8442-7ea3f246eb68",
    "config": null,
    "name": "Relation-composant-ressource",
    "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
    "type": "relation"
  },
  {
    "patterns": null,
    "_id": "47757b30-4499-4cb6-afc5-f3d3f44300e8",
    "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
    "config": {
      "time_interval": 10,
      "alarm_patterns": [
        {
          "v": {
            "state": {
              "val": 3
            },
            "resource": {
              "regex_match": "meta_complex"
            }
          }
        }
      ],
      "threshold_count": 3
    },
    "name": "Test groupement de type complexe",
    "type": "complex"
  },
  {
    "patterns": null,
    "_id": "8c7b1732-522f-4bcf-a7ac-d08bd9c085eb",
    "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
    "config": {
      "time_interval": 10,
      "alarm_patterns": [
        {
          "v": {
            "state": {
              "val": 3
            },
            "resource": {
              "regex_match": "meta_complex"
            }
          }
        }
      ],
      "threshold_count": 3
    },
    "name": "Test groupement de type complexe",
    "type": "complex"
  },
  {
    "patterns": null,
    "_id": "mon_id_de_regle",
    "output_template" : "{{ .Children.Alarm.Value.State.Message }}",
    "config": {
      "time_interval": 10,
      "alarm_patterns": [
        {
          "v": {
            "state": {
              "val": 3
            },
            "resource": {
              "regex_match": "meta_complex"
            }
          }
        }
      ],
      "threshold_count": 3
    },
    "name": "Test groupement de type complexe",
    "type": "complex"
  }
]
```
