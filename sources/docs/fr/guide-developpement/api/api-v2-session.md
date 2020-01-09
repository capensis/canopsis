# Session

L'API session permet deffectuer des stats sur les sessions des utilisateur de Canopsis.

Elle permet notament de reseptionner la `keepalive` envoyé par l'UIV3.

### Configuration 

Les paramètres de configuration (notamment la durée d'une session) se trouvent dans le fichier `canopsis/sources/canopsis/etc/session/session.conf`

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

Permet de mettre à jour les temps en cas de déplacement dans l'application.

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

### /api/v2/sessions

Permet de lister les sessions.

**URL** : `/api/v2/session-hide`

**Méthode**: `GET`

**Permissions requises** : Aucune


**Paramètres**

 - `usernames[]` (optionnel) : noms des utilisateurs dont les sessions doivent être renvoyées
 - `active` (optionnel) : true pour ne renvoyer que les sessions en cours
 - `started_after` (optionnel) : si ce paramètre est défini, ne renvoie que les sessions qui ont commencé après ce timestamp
 - `stopped_before` (optionnel) : si ce paramètre est défini, ne renvoie que les sessions qui ont expirées avant ce timestamp


!!! tip "Note"
    La partie "Interaction avec la BDD" sera prochainement deplacer dans la documentation du webserver.

### Document dans la BDD

Le travail sur les sessions s'effectu dans une collection dédier `Default_session`.

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
cette collection permet de garder un historique des sessions des utilisateur ainsi qu'un comptage des temps passer sur les pages.

