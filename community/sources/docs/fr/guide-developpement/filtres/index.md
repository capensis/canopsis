# Description du langage des filtres

## **Structure des filtres**

Un filtre ou un Pattern est une structure qui contient un tableau à 2 dimensions de `FieldConditions`.

La première dimension contient plusieurs groupes de `FieldConditions`.   
Le filtre match si au moins un groupe match (opérateur **OU**).  
  
Par exemple:

```json
[
    [
        {
            "field": "component",
            "cond": {
                "type": "eq",
                "value": "test-component-1"
            }
        }
    ],
    [
        {
            "field": "component",
            "cond": {
                "type": "eq",
                "value": "test-component-2"
            }
        }
    ]
]
```

Il y a un match si l'attribut `component` vaut `test-component-1` **OU** `test-component-2`.

La seconde dimension contient un seul groupe de `FieldConditions`.  
Pour qu'il y ait match, toutes les conditions doivent être vérifiées (opérateur **ET**).  

Par exemple:

```json
[
    [
        {
            "field": "component",
            "cond": {
                "type": "eq",
                "value": "test-component-1"
            }
        },
        {
            "field": "name",
            "cond": {
                "type": "eq",
                "value": "test-resource-1"
            }
        }
    ]
]
```

Il y a un match si l'attribut `component` vaut `test-component-1` **ET** que l'attribut `name` vaut `test-resource-1`.


Les dimensions peuvent être combinées :

```json
[
    [
        {
            "field": "component",
            "cond": {
                "type": "eq",
                "value": "test-component-1"
            }
        },
        {
            "field": "name",
            "cond": {
                "type": "eq",
                "value": "test-resource-1"
            }
        }
    ],
    [
        {
            "field": "component",
            "cond": {
                "type": "eq",
                "value": "test-component-2"
            }
        }
    ]
]
```



Il y a match si l'attribut `component` vaut `test-component-1` **ET** l'attribut  `name` vaut `test-resource-1` **OU** si l'attribut `component` vaut `test-component-2`.

### **Field condition**

Une `Field condition` contient une règle pour un attribut :

```json
{
    "field": "component",
    "cond": {
        "type": "eq",
        "value": "test-component-2"
    }
}
```

Chaque `Field condition` contient les attributs `field` et `cond`.

#### Field

Un `field` dépend du pattern auquel il est attaché. Chaque `field` correspond à l'attribut éponyme dans l'entité Canopsis, il est donc fortement typé.  
Le type induit un ensemble de conditions liées à celui-ci.

Un pattern d'entité contient :

| Field      |     Type      |
|:----------:|:-------------:|
| _id   |  string |
| name   |    string   |
| category   | string |
| type   | string |
| connector   | string |
| component   | string |
| impact_level   | int |
| last_event_date | time |
| infos.info_name | various |
| component_infos.info_name | various | 

Un pattern d'alarme contient : 

| Field      |     Type      |
|:----------:|:-------------:|
| v.display_name   |  string |
| v.output   |    string   |
| v.long_output   | string |
| v.initial_output   | string |
| v.initial_long_output   | string |
| v.connector   | string |
| v.connector_name   | string |
| v.component | string |
| v.resource | string |
| v.last_comment.m | string | 
| v.ack.a | string | 
| v.state.val | int | 
| v.status.val | int | 
| v.ack | reference | 
| v.ticket | reference | 
| v.canceled | reference | 
| v.snooze | reference | 
| v.activation_date | reference | 
| v.creation_date | time | 
| v.last_event_date | time | 
| v.last_update_date | time | 
| v.ack.t | time | 
| v.resolved | time | 
| v.activation_date | time |
| v.duration | duration |
| tags | string_array |
| v.infos.info_name | various | 

Un pattern d'événement contient :

| Field      |     Type      |
|:----------:|:-------------:|
| connector   |  string |
| connector_name   |    string   |
| component   | string |
| resource   | string |
| output   | string |
| long_output   | string |
| event_type   | string |
| source_type | string |
| state | integer |
| extra.extra_name | various |

Un pattern de comportement périodique contient :

| Field      |     Type      |
|:----------:|:-------------:|
| pbehavior_info.id   |  string |
| pbehavior_info.type   |    string   |
| pbehavior_info.canonical_type   | string |
| pbehavior_info.reason   | string |

##### Type Various

Etant donné que les `infos` d'une entité peuvent être de type différent, il faut explicitement indiquer le type :

```json
[
    [
        {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
                "type": "eq",
                "value": "test-value"
            }
        }
    ]
]
```

Les différents types sont :

- string
- int
- bool
- string_array


#### Conditions

Les conditions sont des règles qui peuvent aboutir à un match ou non. 
L'ensemble des conditions est défini par le type de l'attribut.

!!! warning "Avertissement"
    La condition doit être en adéquation avec le type de l'attribut, sous peine de générer des erreurs


##### Conditions sur les chaines de caractères (String conditions)

- `eq` - si `field` **est égal** à la `value`

```json
[
    [
        {
            "field": "name",
            "cond": {
                "type": "eq",
                "value": "test-resource-1"
            }
        }
    ]
]
```

- `neq` - si `field` **n'est pas égal** à `value`

```json
[
    [
        {
            "field": "name",
            "cond": {
                "type": "neq",
                "value": "test-resource-1"
            }
        }
    ]
]
```

- `is_one_of` - si `field` **est l'un** des éléments de `value`

```json
[
    [
        {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-1",
                "test-resource-2",
                "test-resource-3"
              ]
            }
        }
    ]
]
```

- `is_not_one_of` - si `field` **n'est pas l'un** des éléments de `value`

```json
[
    [
        {
            "field": "name",
            "cond": {
              "type": "is_not_one_of",
              "value": [
                "test-resource-1",
                "test-resource-2",
                "test-resource-3"
              ]
            }
        }
    ]
]
```

- `contain`, `not_contain`, `begin_with`, `not_begin_with`, `end_with`, `not_end_with` - D'autres conditions dont les noms sont parlants.  

```json
[
    [
        {
            "field": "name",
            "cond": {
              "type": "begin_with",
              "value": "test-"
            }
        }
    ]
]
```

- `regexp` - si `field` match au sens regex **regexp**

```json
[
    [
        {
            "field": "name",
            "cond": {
                "type": "regexp",
                "value": "CMDB:(?P<SI_CMDB>.*?)($|,)"
            }
        }
    ]
]
```

- `exist` - si `field` **existe**, signifie chaine non vide.  

```json
[
    [
        {
            "field": "name",
            "cond": {
                "type": "exist",
                "value": true
            }
        }
    ]
]
```

##### Conditions sur les entiers (Int conditions)

- `eq` - si `field` **est égal** à `value`

```json
[
    [
        {
            "field": "impact_level",
            "cond": {
                "type": "eq",
                "value": 2
            }
        }
    ]
]
```

- `neq` - si `field` **n'est pas égal** à `value`

```json
[
    [
        {
            "field": "impact_level",
            "cond": {
                "type": "neq",
                "value": 2
            }
        }
    ]
]
```

- `gt` - si `field` **est plus grand** que `value`

```json
[
    [
        {
            "field": "impact_level",
            "cond": {
                "type": "gt",
                "value": 2
            }
        }
    ]
]
```

- `lt` - si `field` **est plus petit** que `value`

```json
[
    [
        {
            "field": "impact_level",
            "cond": {
                "type": "lt",
                "value": 2
            }
        }
    ]
]
```

##### Conditions booléennes (Bool conditions)

- `eq` - si `field` **est égal** à `value`

```json
[
    [
        {
            "field": "infos.bool_info",
            "field_type: "bool",
            "cond": {
                "type": "eq",
                "value": true
            }
        }
    ]
]
```

##### Conditions de référence (Reference conditions)

- `exist` - si `field` **existe**, si la structure Canopsis contient `field`. 

```json
[
    [
        {
            "field": "v.ack",
            "cond": {
                "type": "exist",
                "value": true
            }
        }
    ]
]
```

!!! note
    Si une alarme n'est pas acquittée, il n'y aura pas de structure `ack`dans l'alarme, le pattern ne match pas.


##### String_array conditions

- `has_every` - si `field` **contient tous** les élément de `value`

```json
[
    [
        {
            "field": "infos.array_info",
            "field_type: "string_array",
            "cond": {
                "type": "has_every",
                "value": [
                    "value-1",
                    "value-2",
                    "value-3"
                ]
            }
        }
    ]
]
```

!!! note
    ["value-1", "value-2", "value-3"] match ce pattern autant que ["value-1", "value-2", "value-3", "value-4"]. En revanche ["value-1", "value-2"] ne match pas.

- `has_one_of` - si `field` **contient au moins** un élément de `value`

```json
[
    [
        {
            "field": "infos.array_info",
            "field_type: "string_array",
            "cond": {
                "type": "has_one_of",
                "value": [
                    "value-1",
                    "value-2",
                    "value-3"
                ]
            }
        }
    ]
]
```

!!! note
    ["value-1"] match ce pattern autant que ["value-3", "value-4"]. En revanche ["value-1", "value-4"] ne match pas.

- `has_not` - si `field` **ne contient aucun** élément de `value`

```json
[
    [
        {
            "field": "infos.array_info",
            "field_type: "string_array",
            "cond": {
                "type": "has_not",
                "value": [
                    "value-1",
                    "value-2"
                ]
            }
        }
    ]
]
```

!!! note
    ["value-3"] match ce pattern. En revanche ["value-1", "value-3"] ne match pas.

- `is_empty` - si `field` est un tableau vide

```json
[
    [
        {
            "field": "infos.array_info",
            "field_type: "string_array",
            "cond": {
                "type": "is_empty",
                "value": true
            }
        }
    ]
]
```

##### Conditions temporelles (Time conditions)

- `relative_time` - si `field` appartient à l'intervalle de `maintenant` à une valeur dans le passé.

```json
[
    [
        {
            "field": "v.last_event_date",
            "cond": {
                "type": "relative_time",
                "value": {
                    "value": 1,
                    "unit": "m"
                }
            }
        }
    ]
]
```

Dans cet exemple, il y a match si `last_event_date` n'a pas plus d'une minute

- `absolute_time` - si `field` appartient à une intervalle fixe de temps

```json
[
    [
        {
            "field": "v.ack.t",
            "cond": {
                "type": "absolute_time",
                "value": {
                    "from": 1605263992,
                    "to": 1605264992
                }
            }
        }
    ]
]
```

##### Conditions de durée (Duration conditions)

- `gt` - si `field` **est plus grand que** la durée `value`

```json
[
    [
        {
            "field": "v.duration",
            "cond": {
                "type": "gt",
                "value": {
                    "value": 3,
                    "unit": "m"
                }
            }
        }
    ]
]
```

- `lt` - si `field` **est plus petit** que la durée `value`

```json
[
    [
        {
            "field": "v.duration",
            "cond": {
                "type": "lt",
                "value": {
                    "value": 3,
                    "unit": "m"
                }
            }
        }
    ]
]
```
