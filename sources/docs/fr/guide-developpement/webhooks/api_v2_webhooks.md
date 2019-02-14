# API Webhooks

L'API Webhooks permet de consulter, créer et supprimer des Webhooks.

### Creation de Webhook

Crée un nouveau Webhook à partir du corps de la requête.

**URL** : `/api/v2/webhook`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de corps de requête** :
```json
{
    "_id" : "declare_external_ticket",
    "hook" : {
        "triggers" : [
            "create"
        ],
        "event_patterns" : [
            {
                "connector" : "zabbix"
            }
        ],
        "alarm_patterns" : [
            {
                "ticket" : null
            }
        ],
        "entity_patterns" : [
            {
                "infos" : {
                    "output" : "MemoryDisk.*"
                }
            }
        ]
    },
    "request" : {
        "method" : "PUT",
        "url" : "{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}",
        "headers" : {
            "Content-type" : "application/json"
        },
        "payload" : "{{ $comp := .Alarm.Value.Component }}{{ $reso := .Alarm.Value.Resource }}{{ $val := .Alarm.Value.Status.Value }}{\"component\": \"{{$comp}}\",\"resource\": \"{{$reso}}\", \"parity\": {{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}},  \"value\": {{$val}} }"
    },
    "declare_ticket" : {
        "ticket_id" : "id",
        "ticket_creation_date" : "timestamp",
        "priority" : "priority"
    }
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le Json ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
    "_id" : "declare_external_ticket",
    "hook" : {
        "triggers" : [
            "create"
        ],
        "event_patterns" : [
            {
                "connector" : "zabbix"
            }
        ],
        "alarm_patterns" : [
            {
                "ticket" : null
            }
        ],
        "entity_patterns" : [
            {
                "infos" : {
                    "output" : "MemoryDisk.*"
                }
            }
        ]
    },
    "request" : {
        "method" : "PUT",
        "url" : "{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}",
        "headers" : {
            "Content-type" : "application/json"
        },
        "payload" : "{{ $comp := .Alarm.Value.Component }}{{ $reso := .Alarm.Value.Resource }}{{ $val := .Alarm.Value.Status.Value }}{\"component\": \"{{$comp}}\",\"resource\": \"{{$reso}}\", \"parity\": {{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}},  \"value\": {{$val}} }"
    },
    "declare_ticket" : {
        "ticket_id" : "id",
        "ticket_creation_date" : "timestamp",
        "priority" : "priority"
    }
}' 'http://<Canopsis_URL>/api/v2/webhook'
```

#### Réponse en cas de réussite

**Condition** : le Webhook est créé

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
"declare_external_ticket"
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

**Condition** : Si un Webhook similaire existe déjà en base.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "Error while creating an webhook"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "Error while creating an webhook"
}
```


### Suppression de Webhook

Supprime un Webhook en fonction de son `id`.

**URL** : `/api/v2/webhook/<webhook_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer le webhook avec l'`id` `declare_external_ticket` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/webhook/declare_external_ticket'
```

#### Réponse en cas de réussite

**Condition** : La suppresion du Webhook a réussi.

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
  "description": "Can not retrieve the webhook data from database, contact your administrator."
}
```


### Récupération des Webhooks

Récupère un ou plusieurs Webhook crée en base.

#### Récupération d'un Webhook par id

**URL** : `/api/v2/webhook/<webhook_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer le webhook avec l'`id` `declare_external_ticket` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/webhook/declare_external_ticket'
```

##### Réponse en cas de réussite

**Condition** : Un Webhook correspondant à l'`id` est trouvé.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "_id" : "declare_external_ticket",
    "hook" : {
        "triggers" : [
            "create"
        ],
        "event_patterns" : [
            {
                "connector" : "zabbix"
            }
        ],
        "alarm_patterns" : [
            {
                "ticket" : null
            }
        ],
        "entity_patterns" : [
            {
                "infos" : {
                    "output" : "MemoryDisk.*"
                }
            }
        ]
    },
    "request" : {
        "method" : "PUT",
        "url" : "{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}",
        "headers" : {
            "Content-type" : "application/json"
        },
        "payload" : "{{ $comp := .Alarm.Value.Component }}{{ $reso := .Alarm.Value.Resource }}{{ $val := .Alarm.Value.Status.Value }}{\"component\": \"{{$comp}}\",\"resource\": \"{{$reso}}\", \"parity\": {{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}},  \"value\": {{$val}} }"
    },
    "declare_ticket" : {
        "ticket_id" : "id",
        "ticket_creation_date" : "timestamp",
        "priority" : "priority"
    }
}
```

##### Réponse en cas d'erreur

**Condition** : Aucun Webhook correspondant à l'`id` n'est trouvé.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "No webhook found with ID declare_external_ticket"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "Can not retrieve the webhook data from database, contact your administrator."
}
```

#### Récupération de tous les Webhooks en base de données.

Récupèr tous les Webhooks stocké en base

**URL** : `/api/v2/webhook`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer tous les webhooks :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/webhook'
```

##### Réponse en cas de réussite

**Condition** : aucune.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
[
    {
      "_id": "declare_external_ticket",
      // ...
    },
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

##### Réponse en cas d'erreur

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "Can not retrieve the webhooks list from database, contact your administrator."
}
```
