# API Statistiques

## Spécification

#### URL

Cette API définit deux routes :

 - `POST /api/v2/stats/<nom de la statistique>` : pour calculer la statistique
   `<nom de la statistique>`.
 - `POST /api/v2/stats`: pour calculer la valeur de plusieurs statistiques.

#### Paramètres

Ces routes acceptent un objet JSON contenant les paramètres suivants :

 - `stats` (uniquement pour la route `/api/v2/stats`, requis) : la liste des
   statistiques à calculer.
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
un).

Le nom de tag peut être utilisé pour filtrer selon :

 - L'identité de l'entité, avec les tags `entity_id` and `entity_type`.
 - Les informations de l'entité, avec les tags `entity_infos.<information_id>`.
   Seules les informations spécifiées dans la [configuration du moteur
   statsng](../../../guide_administrateur/statsng.md#entity-tags) peuvent être
   utilisées.
 - L'alarme, avec les tags `connector`, `connector_name`, `component`,
   `resource` and `alarm_state`.

Le filtre de tag peut être :

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
 - `above_rate`: la proportion d'alarmes dont le temps d'acquittement est
   supérieur au SLA (entre 0 et 1)
 - `below_rate`: la proportion d'alarmes dont le temps d'acquittement est
   inférieur au SLA (entre 0 et 1)

Le SLA doit être indiqué en secondes dans la requête dans un champs
`ack_time_sla`.

### Taux de résolution inférieur ou supérieur au SLA

La statistique `resolve_time_sla` est un objet JSON avec les champs suivants :

 - `above`: le nombre d'alarmes dont le temps de résolution est supérieur au
   SLA
 - `below`: le nombre d'alarmes dont le temps de résolution est inférieur au
   SLA
 - `above_rate`: la proportion d'alarmes dont le temps de résolution est
   supérieur au SLA (entre 0 et 1)
 - `below_rate`: la proportion d'alarmes dont le temps de résolution est
   inférieur au SLA (entre 0 et 1)

Le SLA doit être indiqué en secondes dans la requête dans un champs
`resolve_time_sla`.

### Temps passé dans chaque état

La statistique `time_in_state` est un objet JSON avec :

 - un champ par état (0-3), contenant le temps passé par le watcher dans cet
   état en secondes
 - un champ `total`, contenant le temps total

Les périodes pendant lesquels un pbehavior est actif sont exclues des valeurs
ci-dessus. Le temps total peut donc être inférieur à la durée de la période
`tstop - tstart`.

### Disponibilité

La statistique `availability` est un objet JSON avec les champs suivants :

 - `available_time` : le temps pendant lequel le watcher était dans un état
   disponible en secondes
 - `unavailable_time` : le temps pendant lequel le watcher était dans un état
   indisponible en secondes
 - `available_rate` : la proportion du temps pendant lequel le watcher était
   dans un état disponible (entre 0 et 1)
 - `unavailable_rate` : la proportion du temps pendant lequel le watcher était
   dans un état indisponible (entre 0 et 1)

Un état est considéré comme disponible s'il est inférieur ou égal à la valeur
donnée en paramètre dans le champs `available_state`.

Les périodes pendant lesquels un pbehavior est actif sont exclues des valeurs
ci-dessus. Le temps total `available_time + unavailable_time` peut donc être
inférieur à la durée de la période `tstop - tstart`.

**Remarque :** Si cette statistique est calculée pour plusieurs watchers à la
fois, la somme des temps de disponibilité (ou d'indisponibilité) sera renvoyée.


## Exemples

### Calcul du nombre d'alarmes créées par un composant

`/api/v2/stats/alarms_created`

Requête:

```javascript
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

Réponse:

```javascript
[
    {
        "tags": {},
        "alarms_created": 13
    }
]
```

### Calcul du nombre d'alarmes résolues par ressources d'un composant

`/api/v2/stats/alarms_resolved`

Requête:

```javascript
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

Réponse:

```javascript
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

### Calcul du taux d'acquittement inférieur ou supérieur au SLA

`/api/v2/stats/ack_time_sla`

Requête:

```javascript
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
    "ack_time_sla": 600
}
```

Réponse:

```javascript
[
    {
        "tags": {},
        "ack_time_sla": {
            "above": 3
            "below": 9,
            "above_rate": 0.25,
            "below_rate": 0.75,
        }
    }
]
```

### Calcul du nombre d'alarmes critiques créées par un composant

`/api/v2/stats/alarms_created`

Requête:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "connector": "connector",
            "connector_name": "connector_name",
            "component": "component",
            "alarm_state": 3,
        }
    ]
}
```

Réponse:

```javascript
[
    {
        "tags": {},
        "alarms_created": 13
    }
]
```

### Calcul du temps passé par un watcher dans chaque état

`/api/v2/stats/time_in_state`

Requête:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "entity_id": "watcher_0"
        }
    ]
}
```

Réponse:

```javascript
[
    {
        "tags": {},
        "time_in_state": {
			"total": 2454,
			"0": 1707,
			"1": 105,
			"2": 23,
			"3": 619
		}
    }
]
```

### Calcul du temps pendant lequel un watcher était disponible

`/api/v2/stats/availability`

Requête:

```javascript
{
    "tstart": 1528290000,
    "tstop": 1528293000,
    "filter": [
        {
            "entity_id": "watcher_0"
        }
    ],
    "available_state": 2
}
```

Réponse:

```javascript
[
    {
        "tags": {},
        "availability": {
			"available_time": 1835,
			"unavailable_time": 619,
			"available_rate": 0.747758761206194,
			"unavailable_rate": 0.25224123879380606
		}
    }
]
```

### Calcul de plusieurs statistiques

`/api/v2/stats`

Requête:

```javascript
{
    "stats": ["alarms_created", "alarms_resolved"],
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

Réponse:

```javascript
[
    {
        "tags": {},
        "alarms_created": 13,
        "alarms_created": 8,
    }
]
```
