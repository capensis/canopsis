# Sessions

## Information

Le service web de Canopsis récolte des informations sur les sessions utilisateurs dans l'objectif de permettre la réalisation de statistiques.

Ces informations sont stockées dans la collection `default_session` de MongoDB.

## Document dans la BDD

Les données renvoyées par l'API sont stockées dans la collection `default_session`.

```javascript
{
    "_id": "...",                     // id unique de la session
    "id_beaker_session" : "cm9vdF8xNTc0OTU3MjMz", //Permet d'identifier la session
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
