# API Statistiques

Cette API permet de calculer diverses statistiques sur les entités.


## Spécification

#### URL

  `POST /api/v2/stats/[nom de la statistique/]`


#### Paramètres

Ces routes acceptent un objet JSON contenant les paramètres suivants :

 - `tstart` (optionnel) : début de la période pour laquelle les statistiques
   doivent être calculées (timestamp).
 - `tstop` (optionnel) : fin de la période pour laquelle les statistiques
   doivent être calculées (timestamp).
 - `group_by` (optionnel) : la liste de tags par lesquels les résultats doivent
   être regroupés.
 - `limit` (optionnel) : si `group_by` est défini, le nombre maximal de groupes
   à retourner.
 - `offset` (optionnel) : si `group_by` est défini, l'indice du premier groupe
   à retourner.
 - `filter` (optionnel) : une liste de *groupes d'entités*. Une entité est
   prise en compte dans le calcul des statistiques si elle fait partie d'un des
   groupes d'entités.

Un *groupe d'entités* est un objet JSON contenant des couples
`"<nom de tag>" : <filtre de tag>`. Une entité fait partie d'un groupe
d'entités si chacun de ses tags valide le filtre correspondant (s'il y en a
un). Un filtre de tag peut être :

 - une chaîne de caractères, auquel cas la valeur du tag doit être égale à
   cette chaîne;
 - une liste de chaînes, auquel cas la valeur du tag doit faire partie de cette
   liste;
 - un objet de la forme `{"matches": "<regex>"}`, où `<regex>` est une
   [expression régulière](https://golang.org/pkg/regexp/syntax/), auquel cas la
   valeur du tag doit être reconnue par cette expression régulière.

```javascript
[ // Calcule les statistiques pour les entités appartenant à au moins un des groupes suivants.
    { // Ce groupe contient les entités dont les tags vérifient les conditions suivantes
        "<tag1>": "valeur",                    // la valeur de tag1 est "valeur" ET
        "<tag2>": ["valeur1", "valeur2", ...], // la valeur de tag2 est dans [...] ET
        "<tag3>": {"matches": "valeur\d+"}     // la valeur de tag3 est reconnue par la regex
    },
    // ...
]
```


#### Réponse

En cas de succès, la réponse est un objet JSON avec les champs suivants :

 - `total` : le nombre total de groupes (vaut 1 si `group_by` n'est pas
   défini).
 - `data` : un tableau contenant les résultats de la requête. Chaque élément du
   tableau est un objet contenant :
    - un champ `tags` : un tableau contenant les valeurs des tags défini en
      `group_by` (un tableau vide si `group_by` n'est pas défini)
    - des champs contenant les valeurs des statistiques calculées.


## Statistiques

### Nombre d'alarmes créées

La statistique `alarms_created` vaut le nombre d'alarmes créées pendant une
période.

### Nombre d'alarmes résolues

La statistique `alarms_resolved` vaut le nombre d'alarmes résolues pendant une
période.

### Nombre d'alarmes annulées

La statistique `alarms_canceled` vaut le nombre d'alarmes annulées pendant une
période.


## Exemples

### Calcul du nombre d'alarmes créées par un composant

`/api/v2/stats/alarms_created/`

Requête :
```json
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component"
        }
    ]
}
```

Réponse :
```json
{
    "total": 1,
    "data": [
        {
            "tags": [],
            "alarms_created": 13
        }
    ]
}
```

### Calcul du nombre d'alarmes résolues par ressources d'un composant

`/api/v2/stats/alarms_resolved/`

Requête :
```json
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component"
        }
    ],
    "group_by": ["resource"]
}
```

Réponse :
```json
{
    "total": 3,
    "data": [
        {
            "tags": ["resource1"],
            "alarms_resolved": 4
        },
        {
            "tags": ["resource2"],
            "alarms_resolved": 3
        },
        {
            "tags": ["resource3"],
            "alarms_resolved": 1
        },
    ]
}
```
