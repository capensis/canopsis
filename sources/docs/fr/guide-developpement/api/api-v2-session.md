# Session

L'API session stocke des données sur les sessions utilisateurs de Canopsis.

Les données disponibles sont :

- Nombre de sessions crées par utilisateur
- Temps de navigation passée sur chaque page par utilisateur
- Liste des sessions créées par utilisateur

Pour cela, elle se base notamment sur le signal keepalive envoyé par l'UIv3.


### Configuration 

Les paramètres de configuration (notamment la durée d'une session) se trouvent dans le fichier `canopsis/sources/canopsis/etc/session/session.conf`

### Création d'une session

Création d'une Session.

**URL** : `/api/v2/sessionstart`

**Méthode**: `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

#### Réponse en cas de réussite

**Condition** : la session est créée

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
	"description":"Session Start"
}
```



### Navigation dans l'application

Permet de stocker les temps de session en fonction des différentes pages consultées dans l'application.

**URL** : `/api/v2/session-hide`

**Méthode**: `POST`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Condition** : Session Active

**Exemple de corps de requête** :
```json
{
    "path":"[/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8,view-tab_edd5855b-54f1-4c51-9550-d88c2da60768]"
}
```
#### Réponse en cas de réussite

**Condition** : session-hide bien enregistrer

**Code** : `200 OK`

### Liste des sessions

Permet de lister les sessions.

**URL** : `/api/v2/sessions`

**Méthode**: `GET`

**Permissions requises** : Aucune


**Paramètres**

 - `usernames[]` (optionnel) : noms des utilisateurs dont les sessions doivent être renvoyées
 - `active` (optionnel) : true pour ne renvoyer que les sessions en cours
 - `started_after` (optionnel) : si ce paramètre est défini, ne renvoie que les sessions qui ont commencé après ce timestamp
 - `stopped_before` (optionnel) : si ce paramètre est défini, ne renvoie que les sessions qui ont expiré avant ce timestamp


#### Réponse en cas de réussite

**Condition** : sessions trouver 

**Code** : `200 OK`

```json
[
    {
        "_id": "...",
        "id_beaker_session" : "cm9vdF8xNTc0OTU3MjMz",
        "username": "jacques",
        "start": 1574150400,
        "last_ping": 1574182800,
        "last_visible_ping": 1574182800,
        "last_visible_path": [
            "/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8",
            "view-tab_edd5855b-54f1-4c51-9550-d88c2da60768",
        ],
        "visible_duration": 3750,

        "tab_duration": {
            "/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8": {
                "view-tab_edd5855b-54f1-4c51-9550-d88c2da60768": 3000
            },
            "/view/other-view": {
                "other-tab": 50
            },
            "/exploitation/pbehaviors": 700
        }
    },
    // ...
]
```
### Signal de session active

Envoi d'un keepalive. 

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

**Condition** : Keepalive bien enregistré

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


!!! tip "Note"
    La partie "Interaction avec la BDD" sera prochainement deplacée dans la documentation du webserver.

### Document dans la BDD

Les données renvoyées par l'API sont stockées dans la collection `Default_session`.

```javascript
{
    "_id": "...",                     // id unique de la session
    "id_beaker_session" : "cm9vdF8xNTc0OTU3MjMz", //Permet d'itentifier la session
    "username": "jacques",            // nom de l'utilisateur
    "start": 1574150400,              // date du début de la session
    "last_ping": 1574182800,          // date du dernier ping
    "last_visible_ping": 1574182800,  // date du dernier ping au premier plan
    "last_visible_path": [
        "/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8",     // id de la dernière vue visible au premier plan
        "view-tab_edd5855b-54f1-4c51-9550-d88c2da60768",  // id du dernier onglet visible au premier plan
    ],
    "visible_duration": 3750,  // durée passée au premier plan, en secondes

    // objet contenant les durées passées dans chaque onglet
    // (session["tab_duration"][view_id][tab_id] vaut la durée passée dans
    // l'onglet tab_id de la vue view_id)
    "tab_duration": {
        "/view/da7ac9b9-db1c-4435-a1f2-edb4d6be4db8": {
            "view-tab_edd5855b-54f1-4c51-9550-d88c2da60768": 3000
        },
        "/view/other-view": {
            "other-tab": 50
        },
        "/exploitation/pbehaviors": 700
    }
}
```
Cette collection stocke l'historique des sessions utilisateurs ainsi que le temps passé sur les différentes pages. 
Elle peut servir de base pour l'élaboration de statistiques utilisateurs.
