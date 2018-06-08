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

En cas de succès, la réponse est un tableau JSON contenant des objets avec
chacun :

 - un champ `tags` : un objet contenant les valeurs des tags défini en
   `group_by` (un objet vide si `group_by` n'est pas défini)
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

### Temps d'acquittement moyen

La statistique `mean_ack_time` vaut le temps d'acquittement moyen.

### Temps de résolution moyen

La statistique `mean_resolve_time` vaut le temps de résolution moyen.

### Taux d'acquittement inférieur ou supérieur au SLA

La statistique `ack_time_sla` est un objet JSON avec les champs suivants :

 - `above`: le nombre d'alarmes dont le temps d'acquittement est supérieur au
   SLA
 - `below`: le nombre d'alarmes dont le temps d'acquittement est inférieur au
   SLA
 - `above_rate`: le pourcentage d'alarmes dont le temps d'acquittement est
   supérieur au SLA
 - `below_rate`: le pourcentage d'alarmes dont le temps d'acquittement est
   inférieur au SLA

### Taux de résolution inférieur ou supérieur au SLA

La statistique `resolve_time_sla` est un objet JSON avec les champs suivants :

 - `above`: le nombre d'alarmes dont le temps de résolution est supérieur au
   SLA
 - `below`: le nombre d'alarmes dont le temps de résolution est inférieur au
   SLA
 - `above_rate`: le pourcentage d'alarmes dont le temps de résolution est
   supérieur au SLA
 - `below_rate`: le pourcentage d'alarmes dont le temps de résolution est
   inférieur au SLA


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
[
    {
        "tags": {},
        "alarms_created": 13
    }
]
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
[
    {
        "tags": {"resource": "resource1"},
        "alarms_resolved": 4
    },
    {
        "tags": {"resource": "resource2"},
        "alarms_resolved": 3
    },
    {
        "tags": {"resource": "resource3"},
        "alarms_resolved": 1
    }
]
```
