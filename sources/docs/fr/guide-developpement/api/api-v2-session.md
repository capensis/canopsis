# Session

L'API session permet deffectuer des stats sur les sessions des utilisateur de Canopsis.

### Création d'une session

Création d'une Session.

**URL** : `/api/v2/sessionstart`

**Méthode**: `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

#### Réponse en cas de réussite

**Condition** : la session est créé

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"description":"Session Start"
}
```

### keepalive

Envoie d'un keepalive. 

**URL** : `/api/v2/keepalive`

**Méthode**: `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Condition** Session Active

**Exemple de corps de requête** :
```json
{
    "visible":true
    "path":"[/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8,view-tab_edd5855b-54f1-4c51-9550-d88c2da60768]"
}
```
#### Réponse en cas de réussite

**Condition** : Keepalive bien enregistrer

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"description":"Session keepalive",
    "time":XXXXXXX,
    "visible":true,
    "paths":"[/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8,view-tab_edd5855b-54f1-4c51-9550-d88c2da60768]"
}
```

### /api/v2/session-hide

Permet de mettre a jour les temps en cas de déplacement dans l'application.

**URL** : `/api/v2/session-hide`

**Méthode**: `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Condition** Session Active

**Exemple de corps de requête** :
```json
{
    "path":"[/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8,view-tab_edd5855b-54f1-4c51-9550-d88c2da60768]"
}
```
#### Réponse en cas de réussite

**Condition** : session-hide bien enregistrer

**Code** : `200 OK`

### /api/v2/sessions

Permet de lister les sessions.

**URL** : `/api/v2/session-hide`

**Méthode**: `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune


**Paramètres**

 - `usernames[]` (optionnel) : noms des utilisateurs dont les sessions doivent être renvoyées
 - `active` (optionnel) : true pour ne renvoyer que les sessions en cours
 - `started_after` (optionnel) : si ce paramètre est défini, ne renvoie que les sessions qui ont commencé après ce timestamp
 - `stopped_before` (optionnel) : si ce paramètre est défini, ne renvoie que les sessions qui ont expirées avant ce timestamp
 - `limit`, `offset` (optionnel) : pour gérer la pagination

