# Informations dynamiques

L'API dynamic-infos permet de consulter, créer, modifier et supprimer des règles de gestion des informations dynamiques.

!!! note
    Les règles de gestion des informations dynamiques seront utilisées par le moteur Go `dynamic-infos`, qui n'est pas encore implémenté, pour ajouter des informations aux alarmes.

## Création d'une règle

Crée une nouvelle règle à partir du corps de la requête.

**URL** : `/api/v2/dynamic-infos`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
    "name": "Lien critique",
    "description": "Ajoute un lien externe aux alarmes critiques",
    "alarm_patterns": [
        {
            "v": {
                "connector": "zabbix",
                "state": {
                    "val": 3
                }
            }
        }
    ],
    "infos": [
        {"name": "type", "value": "url"},
        {"name": "url", "value": "http://help.local/zabbix-critical"}
    ]
}
```

L'`_id` de la règle peut optionnellement être défini dans le corps de la requête. En l'absence de champ `_id`, celui-ci est défini automatiquement.

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
    "name": "Lien critique",
    "description": "Ajoute un lien externe aux alarmes critiques",
    "alarm_patterns": [
        {
            "v": {
                "connector": "zabbix",
                "state": {
                    "val": 3
                }
            }
        }
    ],
    "infos": [
        {"name": "type", "value": "url"},
        {"name": "url", "value": "http://help.local/zabbix-critical"}
    ]
}' 'http://<Canopsis_URL>/api/v2/dynamic-infos'
```

### Réponse en cas de réussite

**Condition** : La règle a été créée avec succès. Dans ce cas, les champs suivants sont ajoutés à la règle :

 - `author` : l'utilisateur ayant créé la règle
 - `creation_date` : la date de création de la règle
 - `last_modified_date` : la date de la dernière modification de la règle

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
{
    "_id": "7b17aae1-8c75-440c-a39c-5b304b850171",
    "name": "Lien critique",
    "author": "billy",
    "creation_date": 1576260000,
    "last_modified_date": 1576260000
    // ...
}
```

### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas un document JSON valide.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "Invalid JSON"
}
```

---

**Condition** : Si le corps de la requête n'est pas une règle valide.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "invalid dynamic infos: ..."
}
```

---

**Condition** : Si le corps de la requête contient un id qui est déjà utilisé par une autre règle.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "failed to create dynamic infos: duplicate id 7b17aae1-8c75-440c-a39c-5b304b850171"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "failed to create dynamic infos: ..."
}
```

## Modification d'une règle

Modifie une règle à partir du corps de la requête.

Les champs `author`, `creation_date` et `last_modified_date` du corps de la
requête sont ignorés par cette route. `author` et `creation_date` ne sont pas
modifiés, `last_modified_date` est automatiquement mis à jour.

**URL** : `/api/v2/dynamic-infos/<rule_id>`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :

```json
{
    "_id": "7b17aae1-8c75-440c-a39c-5b304b850171",
    "name": "Lien critique",
    "description": "Ajoute un lien externe aux alarmes critiques",
    "alarm_patterns": [
        {
            "v": {
                "connector": "zabbix",
                "state": {
                    "val": 3
                }
            }
        }
    ],
    "infos": [
        {"name": "type", "value": "url"},
        {"name": "title", "value": "Aide pour les alarmes critiques"},
        {"name": "url", "value": "http://help.local/zabbix-critical"}
    ]
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
    "_id": "7b17aae1-8c75-440c-a39c-5b304b850171",
    "name": "Lien critique",
    "description": "Ajoute un lien externe aux alarmes critiques",
    "alarm_patterns": [
        {
            "v": {
                "connector": "zabbix",
                "state": {
                    "val": 3
                }
            }
        }
    ],
    "infos": [
        {"name": "type", "value": "url"},
        {"name": "title", "value": "Aide pour les alarmes critiques"},
        {"name": "url", "value": "http://help.local/zabbix-critical"}
    ]
}' 'http://<Canopsis_URL>/api/v2/dynamic-infos/7b17aae1-8c75-440c-a39c-5b304b850171'
```

### Réponse en cas de réussite

**Condition** : La règle a été modifiée avec succès.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
{
    "_id": "7b17aae1-8c75-440c-a39c-5b304b850171",
    "name": "Lien critique",
    // ...
}
```

### Réponse en cas d'erreur

**Condition** : S'il n'y a pas de règle correspondant à l'id.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "no dynamic infos rule with id 7b17aae1-8c75-440c-a39c-5b304b850171"
}
```

## Suppression d'une règle

Supprime une règle.

**URL** : `/api/v2/dynamic-infos/<rule_id>`

**Méthode** : `DELETE`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut supprimer la règle avec l'id `7b17aae1-8c75-440c-a39c-5b304b850171` :

```sh
curl -X DELETE -u root:root 'http://<Canopsis_URL>/api/v2/dynamic-infos/7b17aae1-8c75-440c-a39c-5b304b850171'
```

### Réponse en cas de réussite

**Condition** : La règle a été supprimée avec succès.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{}
```

### Réponse en cas d'erreur

**Condition** : S'il n'y a pas de règle correspondant à l'id.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "no dynamic infos rule with id 7b17aae1-8c75-440c-a39c-5b304b850171"
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "name": "",
  "description": "failed to delete dynamic infos"
}
```

## Récupération d'une règle

**URL** : `/api/v2/dynamic-infos/<rule_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer la règle avec l'id `7b17aae1-8c75-440c-a39c-5b304b850171` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/dynamic-infos/7b17aae1-8c75-440c-a39c-5b304b850171'
```

### Réponse en cas de réussite

**Condition** : La règle existe.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
{
    "_id": "7b17aae1-8c75-440c-a39c-5b304b850171",
    "name": "Lien critique",
    // ...
}
```

### Réponse en cas d'erreur

**Condition** : S'il n'y a pas de règle correspondant à l'id.

**Code** : `404 NOT FOUND`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "no dynamic infos rule with id 7b17aae1-8c75-440c-a39c-5b304b850171"
}
```

## Récupération d'une liste de règles

**URL** : `/api/v2/dynamic-infos`

**Méthode** : `GET`

**Paramètres** :

 - `search` : la valeur à rechercher.
 - `search_fields` : les noms des champs sur lesquels la recherche est effectuée, séparés par des virgules. Les valeurs acceptés sont `_id`, `name`, `description`, `infos.name` et `infos.value`. Par défault la recherche est effectuée sur tous ces champs.
 - `limit` : le nombre de règles à renvoyer.
 - `offset` : le nombre de règles à passer.

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer toutes les règles :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/dynamic-infos'
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` pour récupérer les règles ajoutant une information avec le nom `url` :

```sh
curl -X GET -u root:root 'http://<Canopsis_URL>/api/v2/dynamic-infos?search=url&search_fields=infos.name'
```

### Réponse en cas de réussite

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
{
    "count": 17,  // Le nombre total de règles correspondant à la recherche
    "rules": [
        {
            "_id": "7b17aae1-8c75-440c-a39c-5b304b850171",
            // ...
        },
        // ...
    ]
}
```

