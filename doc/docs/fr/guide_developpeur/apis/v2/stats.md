# API Statistiques

## Routes

### Calcul d'une statistique

#### URL

`POST /api/v2/stats/<nom de la statistique>`

#### Paramètres

Cette route accepte en requête un objet JSON contenant les paramètres
suivants :

 - `tstart` (optionnel) : le début de la dernière période pour laquelle la
   statistique doit être calculées (timestamp).
 - `tstop` (optionnel) : le fin de la dernière période pour laquelle la
   statistique doit être calculées (timestamp).
 - `periods` (optionnel) : le nombre de périodes pour lesquelles la statistique
   doit être calculées, en secondes.
 - `group_by` (optionnel) : la liste de tags par lesquels les résultats doivent
   être regroupés. Les tags disponibles sont les mêmes que ceux utilisés dans
   les *groupes d'alarmes*, et sont définis ci-dessous.
 - `filter` (optionnel) : une liste de *groupes d'alarmes* (voir ci-dessous).
   Une alarme est prise en compte dans le calcul des statistiques si elle fait
   partie d'un des groupes d'alarmes.
 - `parameters` (optionnel) : un objet contenant les paramètres spécifiques à
   la statistique calculée. Ces paramètres sont précisés dans la documentation
   de chacune des statistiques.

Un *groupe d'alarmes* est un objet JSON contenant des couples
`"<nom de tag>" : <filtre de tag>`. Une alarme fait partie d'un groupe
d'alarmes si chacun de ses tags valide le filtre correspondant (s'il y en a
un).

Le nom de tag peut être utilisé pour filtrer selon :

 - L'identité de l'entité qui a créé l'alarme, avec les tags `entity_id` et
   `entity_type`.
 - Les informations de l'entité, avec les tags `entity_infos.<information_id>`.
   Seules les informations spécifiées dans la [configuration du moteur
   statsng](../../../guide_administrateur/statsng.md#entity-tags) peuvent être
   utilisées.
 - L'alarme, avec les tags `connector`, `connector_name`, `component`,
   `resource` et `alarm_state`.

Le filtre de tag peut être :

 - une chaîne de caractères, auquel cas la valeur du tag doit être égale à
   cette chaîne;
 - une liste de chaînes, auquel cas la valeur du tag doit faire partie de cette
   liste;
 - un objet de la forme `{"matches": "<regex>"}`, où `<regex>` est une
   [expression régulière](https://golang.org/pkg/regexp/syntax/), auquel cas la
   valeur du tag doit être reconnue par cette expression régulière.

```javascript
[ // Calcule les statistiques pour les alarmes appartenant à au moins un des groupes suivants.
    { // Ce groupe contient les alarmes dont les tags vérifient les conditions suivantes
        "<tag1>": "valeur",                    // la valeur de tag1 est "valeur" ET
        "<tag2>": ["valeur1", "valeur2", ...], // la valeur de tag2 est dans [...] ET
        "<tag3>": {"matches": "valeur\d+"}     // la valeur de tag3 est reconnue par la regex
    },
    // ...
]
```

#### Réponse

En cas de succès, la réponse est un tableau JSON contenant les groupes obtenus
en groupant par les tags précisés dans `group_by`.

Chaque groupe est un objet contenant les champs suivants :

 - `tags` : un objet contenant les valeurs des tags définis en `group_by`.
 - `periods` : un tableau contenant les statistiques pour chaque période. Les
   périodes sont triées chronologiquement.

Chaque période est un objet contenant les champs suivants :

 - `tstart` : la date de début de la période (timestamp).
 - `tstop` : la date de fin de la période (timestamp).
 - `<nom de la statistique>` : la valeur de la statistique. Le type de cette
   valeur dépend de la statistique.

#### Exemple

La requête suivante renvoie le taux de résolution conforme à un SLA par jour
pour chaque ressource du composant `c`, pendant la semaine du 23 au 29 juillet
2018.

```javascript
POST /api/v2/stats/resolve_time_sla
{
    "filter": [{
        "component": "c",
        "entity_type": "resource"
    }],
    "group_by": ["resource"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "periods": 7,
    "parameters": {
        "sla": 3600
    }
}
```

Le calcul prend en compte les alarmes dont le composant vaut `c` et dont le
type d'entité est `resource`, c'est-à-dire les alarmes créées par les
ressources du composant `c`.

La statistique est calculée pour sept périodes consécutives et de même durée.
La dernière période commence à `tstart` (le 29 juillet 2018 à 00:00) et se
termine à `tstop` (le 30 juillet 2018 à 00:00). Cela correspond donc à sept
périodes d'une durée d'une journée, du 23 au 29 juillet.

La valeur du SLA est un paramètre spécifique à la statistique
`resolve_time_sla`, et est donc précisée dans le champ `parameters`.

Le document JSON ci-dessous est un exemple de réponse à la requête précédente.
La valeur de la statistique est un dictionnaire contenant plusieurs valeurs.
Voir la documentation de `resolve_time_sla` pour plus de détails.

```javascript
[ // Tableau de groupes
    {
        "tags": { // Tags du groupe
            "resource": "resource1"
        },
        "periods": [ // Tableau des périodes
            {
                "tstart": 1532296800, // 23 juillet à 00:00
                "tstop": 1532383200, // 24 juillet à 00:00
                "resolve_time_sla": { // Valeur de la statistique
                    "above": 10,
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9
                }
            },
            {
                "tstart": 1532383200, // 24 juillet à 00:00
                "tstop": 1532469600, // 25 juillet à 00:00
                "resolve_time_sla": {
                    "above": 0,
                    "below": 47,
                    "above_rate": 0,
                    "below_rate": 1
                }
            },
            // ...
            {
                "tstart": 1532815200, // 29 juillet à 00:00
                "tstop": 1532901600, // 30 juillet à 00:00
                "resolve_time_sla": {
                    "above": 4,
                    "below": 28,
                    "above_rate": 0.125,
                    "below_rate": 0.875
                }
            }
        ]
    },
    {
        "tags": {
            "resource": "resource2"
        },
        "periods": [
            // ...
        ]
    },
    // ...
]
```


### Calcul de plusieurs statistiques

#### URL

`POST /api/v2/stats`

#### Paramètres

Cette route accepte en requête un objet JSON contenant les mêmes paramètres que
la route précédente, avec deux exceptions :

 - Un champ supplémentaire `stats` (requis) contenant la liste des statistiques
   à calculer.
 - Le champ `parameters` contient un objet associant à chaque statistique ses
   paramètres.

#### Réponse

La réponse a le même format que pour la route précédente. Chaque période
contient les valeurs de plusieurs statistiques.

#### Exemple

La requête suivante renvoie le taux de résolution conforme à un SLA et le
nombre d'alarmes créées par jour pour chaque ressource du composant `c`, pendant
la semaine du 23 au 29 juillet 2018.

```javascript
POST /api/v2/stats
{
    "stats": ["resolve_time_sla", "alarms_created"],
    "filter": [{
        "component": "c",
        "entity_type": "resource"
    }],
    "group_by": ["resource"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "periods": 7,
    "parameters": {
        "resolve_time_sla": {
            "sla": 3600
        }
    }
}
```

Le corps de cette requête est identique à l'exemple précédent, à deux
exceptions près :

 - La liste des statistiques à calculer a été ajoutée dans le champ `stats`.
 - Le paramètre `sla` qui était dans le champ `parameters` a été déplacé dans
   `parameters.resolve_time_sla`. La statistique `alarms_created` ne prend pas
   de paramètres. Si elle en prenait, ils devraient être définis dans
   `parameters.alarms_created`.

Le document JSON ci-dessous est un exemple de réponse à la requête précédente.

```javascript
[ // Tableau de groupes
    {
        "tags": { // Tags du groupe
            "resource": "resource1"
        },
        "periods": [ // Tableau des périodes
            {
                "tstart": 1532296800, // 23 juillet à 00:00
                "tstop": 1532383200, // 24 juillet à 00:00
                // Valeurs des deux statistiques
                "alarms_created": 100,
                "resolve_time_sla": {
                    "above": 10
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9
                }
            },
            // ...
        ]
    },
    // ...
]
```


## Statistiques

### Nombre d'alarmes créées

La statistique `alarms_created` renvoie le nombre d'alarmes créées. Les alarmes
créées alors qu'un pbehavior est actif ne sont pas prises en compte. Elle ne
prends pas de paramètres.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/alarms_created
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_created": 100
            }
        ]
    },
    // ...
]
```

### Nombre d'alarmes résolues

La statistique `alarms_resolved` renvoie le nombre d'alarmes résolues. Les
alarmes *créées* alors qu'un pbehavior est actif ne sont pas prises en compte.
Elle ne prend pas de paramètres.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/alarms_resolved
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_resolved": 86
            }
        ]
    },
    // ...
]
```

### Nombre d'alarmes annulées

La statistique `alarms_canceled` renvoie le nombre d'alarmes annulées. Les
alarmes *créées* alors qu'un pbehavior est actif ne sont pas prises en compte.
Elle ne prend pas de paramètres.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/alarms_canceled
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarms_canceled": 7
            }
        ]
    },
    // ...
]
```

### Temps d'acquittement moyen

La statistique `mean_ack_time` renvoie le temps d'acquittement moyen en
secondes. Les alarmes *créées* alors qu'un pbehavior est actif ne sont pas
prises en compte. Elle ne prend pas de paramètres.

Le temps d'acquittement est la différence entre la date du *premier*
acquittement et la date de création de l'alarme.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/mean_ack_time
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "mean_ack_time": 426
            }
        ]
    },
    // ...
]
```

### Temps de résolution moyen

La statistique `mean_resolve_time` renvoie le temps de résolution moyen en
secondes. Les alarmes *créées* alors qu'un pbehavior est actif ne sont pas
prises en compte. Elle ne prend pas de paramètres.

Le temps de résolution est la différence entre la date de résolution et la date
de création de l'alarme.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/mean_resolve_time
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "mean_resolve_time": 2687
            }
        ]
    },
    // ...
]
```

### Taux d'acquittement inférieur ou supérieur au SLA

La statistique `ack_time_sla` renvoie les nombres et taux d'acquittement
inférieur ou supérieur à un SLA. Les alarmes *créées* alors qu'un pbehavior est
actif ne sont pas prises en compte. Elle prend un paramètre `sla` dont la
valeur est le SLA en secondes, et renvoie un objet JSON contenant les champs
suivants :

 - `above` : le nombre d'alarmes dont le temps d'acquittement est supérieur au
   SLA.
 - `below` : le nombre d'alarmes dont le temps d'acquittement est inférieur au
   SLA.
 - `above_rate` : la proportion d'alarmes dont le temps d'acquittement est
   supérieur au SLA (entre 0 et 1).
 - `below_rate` : la proportion d'alarmes dont le temps d'acquittement est
   inférieur au SLA (entre 0 et 1).

Le temps d'acquittement est la différence entre la date du *premier*
acquittement et la date de création de l'alarme.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/ack_time_sla
{
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "sla": 900
    }
}
```

Réponse :

```javascript
[
    {
        "tags": {},
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "ack_time_sla": {
                    "above": 10
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9,
                }
            }
        ]
    }
]
```

### Taux de résolution inférieur ou supérieur au SLA

La statistique `resolve_time_sla` renvoie les nombres et taux de résolution
inférieur ou supérieur à un SLA. Les alarmes *créées* alors qu'un pbehavior est
actif ne sont pas prises en compte. Elle prend un paramètre `sla` dont la
valeur est le SLA en secondes, et renvoie un objet JSON contenant les champs
suivants :

 - `above` : le nombre d'alarmes dont le temps de résolution est supérieur au
   SLA.
 - `below` : le nombre d'alarmes dont le temps de résolution est inférieur au
   SLA.
 - `above_rate` : la proportion d'alarmes dont le temps de résolution est
   supérieur au SLA (entre 0 et 1).
 - `below_rate` : la proportion d'alarmes dont le temps de résolution est
   inférieur au SLA (entre 0 et 1).

Le temps de résolution est la différence entre la date de résolution et la date
de création de l'alarme.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/resolve_time_sla
{
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "sla": 3600
    }
}
```

Réponse :

```javascript
[
    {
        "tags": {},
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "ack_time_sla": {
                    "above": 10
                    "below": 90,
                    "above_rate": 0.1,
                    "below_rate": 0.9,
                }
            }
        ]
    }
]
```

### Temps passé dans chaque état

La statistique `time_in_state` renvoie un objet JSON avec :

 - un champ par état (entre 0 et 3), contenant le temps passé par l'entité dans
   cet état en secondes.
 - un champ `total`, contenant le temps total.

Les périodes pendant lesquels un pbehavior est actif sont exclues des valeurs
ci-dessus. Le temps total peut donc être inférieur à la durée de la période
`tstop - tstart`.

Cette statistique ne peut être calculée que pour des groupes contenant une
seule entité. Il est donc nécessaire de s'assurer que chaque groupe n'en
contient qu'une, par exemple en ajoutant `entity_id` au paramètre `group_by`.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/time_in_state
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "periods": 2
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532728800,
                "tstop": 1532815200,
                "time_in_state": {
                    0: 48159,
                    1: 34051,
                    2: 2203,
                    3: 1387,
                    "total": 85800
                }
            },
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "time_in_state": {
                    0: 52563,
                    1: 28465,
                    2: 4245,
                    3: 527,
                    "total": 85800
                }
            }
        ]
    },
    // ...
]
```

### Disponibilité

La statistique `availability` renvoie les temps et taux de disponibilité et
d'indisponibilité. Elle prend un paramètre `available_state` dont la valeur est
l'état jusqu'auquel une entité est considérée comme étant disponible. Elle
renvoie un objet JSON contenant les champs suivants :

 - `available` : le temps pendant lequel l'entité était dans un état
   disponible (inférieur ou égal à `available_state`) en secondes.
 - `unavailable` : le temps pendant lequel l'entité était dans un état
   indisponible (strictement supérieur à `available_state`) en secondes.
 - `available_rate` : la proportion du temps pendant lequel l'entité était dans
   un état disponible (inférieur ou égal à `available_state`). Cette valeur est
   entre 0 et 1.
 - `unavailable_rate` : la proportion du temps pendant lequel l'entité était
   dans un état indisponible (strictement supérieur à `available_state`). Cette
   valeur est entre 0 et 1.

Les périodes pendant lesquels un pbehavior est actif sont exclues des valeurs
ci-dessus. Le temps total `available + unavailable` peut donc être
inférieur à la durée de la période `tstop - tstart`.

Cette statistique ne peut être calculée que pour des groupes contenant une
seule entité. Il est donc nécessaire de s'assurer que chaque groupe n'en
contient qu'une, par exemple en ajoutant `entity_id` au paramètre `group_by`.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/availability
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "available_state": 1
    }
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "time_in_state": {
                    "available": 81028,
                    "unavailable": 4772,
                    "available_rate": 0.9443822843822843,
                    "unavailable_rate": 0.05561771561771562
                }
            }
        ]
    },
    // ...
]
```

### Maintenance

La statistique `maintenance` renvoie le temps pendant lequel un pbehavior était
ou n'était pas actif sur une entité. Elle ne prend pas de paramètres et renvoie
un objet JSON contenant les champs suivants :

 - `maintenance` : le temps pendant lequel l'entité avait un pbehavior actif,
   en secondes.
 - `no_maintenance` : le temps pendant lequel l'entité n'avait pas de pbehavior
   actif, en secondes.

Cette statistique ne peut être calculée que pour des groupes contenant une
seule entité. Il est donc nécessaire de s'assurer que chaque groupe n'en
contient qu'une, par exemple en ajoutant `entity_id` au paramètre `group_by`.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/maintenance
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "maintenance": {
                    "maintenance": 600,
                    "no_maintenance": 85800
                }
            }
        ]
    },
    // ...
]
```

### Indice de fiabilité

La statistique `mtbf` (Mean Time Between Failures) renvoie l'indice de
fiabilité, c'est-à-dire le temps hors maintenance divisé par le nombre
d'indisponibilités.

Cette statistique ne peut être calculée que pour des groupes contenant une
seule entité. Il est donc nécessaire de s'assurer que chaque groupe n'en
contient qu'une, par exemple en ajoutant `entity_id` au paramètre `group_by`.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/mtbf
{
    "filter": [{
        "component": "c"
    }],
    "group_by": ["entity_id"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "entity_id": "ressource1/c"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "mtbf": 31.931522143654632
            }
        ]
    },
    // ...
]
```

### Liste d'alarmes

La statistique `alarm_list` renvoie une liste d'alarmes. Elle ne prend pas de
paramètres, et renvoie un tableau d'objets JSON contenant les tags de l'alarme
(`entity_id`, `entity_type`, `entity_infos.<information_id>`, `connector`,
`connector_name`, `component`, `resource` et `alarm_state`), et les champs
suivants :

 - `time` : la date de création de l'alarme.
 - `pbehavior` : `true` s'il y avait un pbehavior actifs quand l'alarme a été
   créé, `false` sinon.
 - `resolve_time` : le temps de résolution de l'alarme.

Seules les alarmes résolues sont prises en compte.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/alarm_list
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarm_list": [
                    {
                        "time": 1532815202,
                        "entity_id": "resource1/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource1",
                        "alarm_state": "3",
                        "pbehavior": "False",
                        "value": 157
                    },
                    {
                        "time": 1532815325,
                        "entity_id": "resource2/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource2",
                        "alarm_state": "1",
                        "pbehavior": "False",
                        "value": 849
                    },
                    // ...
                ]
            }
        ]
    },
    // ...
]
```

### Groupes d'entités impactées par le plus d'alarmes

La statistique `most_alarms_impacting` renvoie une liste contenant les groupes
d'entités impactés par le plus d'alarmes. Elle prend les paramètres suivants :

 - `group_by` (requis) : une liste de tags utilisés pour regrouper les entités.
 - `filter` (optionnel) : un filtre d'entités. Le format de ce paramètre est le
   même que celui du champ `filter` principal.
 - `limit` (optionnel) : le nombre maximal de groupes à renvoyer.

Les paramètres `group_by` est `filter` doivent être définis dans le champ
`parameters` (ou `parameters.most_alarms_impacting` pour la route
`/api/v2/stats`), et sont distincts des champs `group_by` et `filter`
principaux. Par exemple, pour obtenir les ressources impactées par le plus
d'alarmes regroupées par composant, le champ `parameters.group_by` doit
contenir `resource` (pour calculer le nombre d'alarme par ressource), le champ
`parameters.filter` doit contenir `"entity_type": "resource"` (pour ne calculer
le nombre d'alarmes que pour les ressources), et le champ `group_by` doit
contenir `component` (pour regrouper les résultats par composant). Voir
l'exemple ci-dessous pour la requête complète.

La requête renvoie une liste d'objets triés par nombre d'alarmes décroissant,
avec les champs suivants :

 - `tags` : les tags du groupe d'entités.
 - `value` : le nombre d'alarmes impactant ce groupe d'entités.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/most_alarms_impacting
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "group_by": ["resource"],
        "filter": [{
            "entity_type": "resource"
        }],
        "limit": 2
    }
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "most_alarms_impacting": [
                    {
                        "tags": {
                            "resource": "resource3"
                        },
                        "value": 451
                    },
                    {
                        "tags": {
                            "resource": "resource1"
                        },
                        "value": 210
                    }
                ]
            }
        ]
    },
    // ...
]
```

### Entités avec le pire indice de fiabilité

La statistique `worst_mtbf` renvoie une liste de groupes d'entités ayant le
pire indice de fiabilité. La requête prend les paramètres suivants :

 - `group_by` (requis) : une liste de tags utilisés pour regrouper les entités.
 - `filter` (optionnel) : un filtre d'entités. Le format de ce paramètre est le
   même que celui du champ `filter` principal.
 - `limit` (optionnel) : le nombre maximal de groupes à renvoyer.

Les paramètres `group_by` est `filter` doivent être définis dans le champ
`parameters` (ou `parameters.wordt_mtbf` pour la route `/api/v2/stats`), et
sont distincts des champs `group_by` et `filter` principaux. Par exemple, pour
obtenir les ressources avec le pire MTBF regroupées par composant, le champ
`parameters.group_by` doit contenir `resource` (pour calculer le nombre
d'alarme par ressource), le champ `parameters.filter` doit contenir
`"entity_type": "resource"` (pour ne calculer le nombre d'alarmes que pour les
ressources), et le champ `group_by` doit contenir `component` (pour regrouper
les résultats par composant). Voir l'exemple ci-dessous pour la requête
complète.

La requête renvoie une liste d'objets triés par indice de fiabilité croissant,
avec les champs suivants :

 - `tags` : les tags du groupe d'entités.
 - `value` : l'indice de fiabilité.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/worst_mtbf
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600,
    "parameters": {
        "group_by": ["resource"],
        "filter": [{
            "entity_type": "resource"
        }],
        "limit": 2
    }
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "worst_mtbf": [
                    {
                        "tags": {
                            "resource": "resource3"
                        },
                        "value": 45
                    },
                    {
                        "tags": {
                            "resource": "resource1"
                        },
                        "value": 157
                    }
                ]
            }
        ]
    },
    // ...
]
```

### Alarmes les plus longues

La statistique `longest_alarms` renvoie une liste des alarmes qui ont pris le
plus de temps à être résolues. La requête prend les paramètres suivants :

 - `limit` (optionnel) : le nombre maximal d'alarmes à renvoyer.

La requête renvoie un tableau d'objets JSON contenant les tags de l'entité qui
a créé l'alarme (`entity_id`, `entity_type`, `entity_infos.<information_id>`,
`connector`, `connector_name`, `component`, `resource` et `alarm_state`), et
les champs suivants :

 - `time` : la date de création de l'alarme.
 - `pbehavior` : `true` s'il y avait un pbehavior actifs quand l'alarme a été
   créé, `false` sinon.
 - `resolve_time` : le temps de résolution de l'alarme.

#### Exemple

Requête :

```javascript
POST /api/v2/stats/longest_alarms
{
    "group_by": ["component"],
    "tstart": 1532815200,
    "tstop": 1532901600
}
```

Réponse :

```javascript
[
    {
        "tags": {
            "component": "component1"
        },
        "periods": [
            {
                "tstart": 1532815200,
                "tstop": 1532901600,
                "alarm_list": [
                    {
                        "time": 1532895472,
                        "entity_id": "resource2/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource2",
                        "alarm_state": "1",
                        "pbehavior": "False",
                        "value": 4892
                    },
                    {
                        "time": 1532854763,
                        "entity_id": "resource1/component1",
                        "entity_type": "resource",
                        "connector": "connector",
                        "connector_name": "connector_name",
                        "resource": "resource1",
                        "alarm_state": "3",
                        "pbehavior": "False",
                        "value": 3542
                    },
                   // ...
                ]
            }
        ]
    },
    // ...
]
```
