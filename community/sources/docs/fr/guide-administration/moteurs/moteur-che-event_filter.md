# `engine-che` - Event-filter

L'event-filter est une fonctionnalité du moteur [`engine-che`](moteur-che.md) permettant de définir des règles manipulant les évènements.

Des exemples pratiques d'utilisation de l'event-filter sont disponibles dans la partie [Exemples](#exemples).

## Règles

Une règle est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel) : un identifiant unique (généré automatiquement s'il n'est pas défini par l'utilisateur).
 - `description` (optionnel) : une description de la règle, donnée par l'utilisateur.
 - `type` (requis) : le type de la règle (voir [Types de règles](#types-de-regles) pour plus de détails).
 - `patterns` (optionnel) : une liste de patterns permettant de sélectionner les évènements auxquels la règle doit être appliquée. Si le paramètre est absent, la règle est appliquée à tous les évènements (voir [Patterns](#patterns) pour plus de détails).
 - `priority` (optionnel, 0 par défaut) : la priorité de la règle. Les règles sont appliquées par ordre de priorité croissante.
 - `enabled` (optionnel, `true` par défaut) : `false` pour désactiver la règle.

#### Exemple

```json
{
    "type": "drop",
    "patterns": [
        {"resource": "invalid_resource"}
    ],
    "priority": 10
}
```

Le `type` de cette règle vaut `drop`, indiquant que c'est une règle qui supprime les évènements.

Le pattern de cette règle sélectionne les évènements dont la ressource vaut `invalid_resource` (voir [Patterns](#patterns) pour plus de détails).

Cette règle supprime donc les évènements dont la ressource vaut `invalid_resource`.

### Application des règles

Lors de la réception d'un évènement par le moteur `engine-che`, les règles sont parcourues par ordre de priorité croissante. Si l'évènement est reconnu par l'un des `patterns` d'une règle, celle-ci est appliquée à l'évènement. Le traitement effectué dépend du `type` de la règle (voir [Types de règles](#types-de-regles) pour plus de détails).

**Note :** Pour pouvoir être traité par Canopsis, un évènement doit respecter l'une des conditions suivantes :

 - son champ `source_type` vaut `component`, et son champ `component` est défini et ne vaut pas `""`

**ou**

 - son champ `source_type` vaut `resource`, et ses champs `component` et `resource` sont définis et ne valent pas `""`.

Les évènements ne respectant pas l'une de ces conditions en entrée du moteur `engine-che` ou en sortie de l'event-filter sont supprimés.

### Patterns

Les patterns d'une règle permettent de sélectionner les évènements auxquels elle doit être appliquée.

!!! note
    Afin de connaître l'intitulé exact des patterns que vous souhaitez utiliser, vous pouvez, en tant qu'administrateur, prendre une alarme représentative dans votre Bac à alarmes, et vous aider du bouton d'action « Liste des variables disponibles ».

    Cette vue vous présente le contenu de vos alarmes et de vos entités, sous la forme d'un arbre.

    Elle vous permet par exemple de savoir que le composant est directement accessible depuis le champ `component`, tandis que les hostgroups se situent dans `infos.hostgroups.value`.

#### Patterns simples

Un pattern peut être défini comme un objet JSON contenant les valeurs de certains champs d'un évènement.

Par exemple, une règle contenant le pattern suivant sera appliquée aux évènements dont le composant vaut `component_name` et dont la ressource vaut `resource_name`, ainsi qu'aux évènements dont le composant vaut `foobar` :

```json
"patterns": [
    {"component": "component_name", "resource": "resource_name"},
    {"component": "foobar"}
]
```

#### Patterns avancés

Pour plus d'expressivité, il est possible d'associer à un champ un objet contenant des couples `operateur: valeur`. Les opérateurs disponibles sont :

 - `>=`, `>`, `<`, `<=` : compare une valeur numérique à une autre valeur.
 - `regex_match` : filtre la valeur d'une clé selon une expression régulière. La syntaxe des expressions régulières est [définie dans un document dédié](../../guide-utilisation/formats-et-syntaxe/format-regex.md).

Par exemple, le pattern suivant sélectionne les évènements dont la criticité est comprise entre 1 et 3 (mineur, majeur ou critique) et dont l'output vérifie une expression régulière :

```json
"patterns": [
    {
        "state": {
            ">=": 1,
            "<=": 3
        },
        "output": {
            "regex_match": "Warning: CPU Load is critical \(.*\)"
        }
    }
]
```

Il est également possible de créer des patterns basés sur l'absence d'un champ dans l'événement. Par exemple, le pattern suivant sélectionne les événements qui ne contiennent pas le champ `customer`.

```json
"patterns": [
    {
        "customer": null
    }
]
```

Des opérateurs avancés permettent de tester les éléments d'un tableau :

- `has_every` : le tableau doit contenir **tous** les éléments listés dans l'opérateur.
- `has_one_of` : le tableau doit contenir **au moins un** des éléments listés dans l'opérateur.
- `has_not` : le tableau ne doit contenir **aucun** des éléments listés dans l'opérateur.

Par exemple, avec le pattern `has_every` suivant :

```json
"patterns": [
    {
        "hostgroups": {
            "has_every": [
                "windows",
                "windows-2012r2"
            ]
        }
    }
]
```

Un événement contenant le tableau suivant sera sélectionné car il contient à la fois `windows` et `windows-2012r2` :

```json
{
    "hostgroups": [
        "windows-2012r2",
        "windows",
        "vmware-windows-guests",
        "vmware-guests",
    ]
}
```

Tandis qu'un événement contenant le tableau suivant ne sera pas sélectionné car il ne contient qu'un des 2 éléments de `has_every` :
```json
{
    "hostgroups": [
        "windows-2012r2",
        "debian",
        "vmware-windows-guests",
        "vmware-guests",
    ]
}
```

### Types de règles

#### `drop`

Lorsqu'une règle de type `drop` est appliquée à un évènement, cet évènement est supprimé. Les règles suivantes ne sont pas appliquées à cet évènement, et il est ignoré par Canopsis.

Lorsqu'un évènement provoque le déclenchement d'une règle « drop », le moteur `engine-che` l'affiche sur sa sortie de log :

```
2019/05/02 12:45:19 event dropped by event filter: {"hostgroups":["HG_FOOBAR"],"event_type":"check","execution_time":2.139087200164795,"timestamp":1556793914,"component":"foobar","state_type":0,"source_type":"resource","resource":"PING","current_attempt":1,"connector":"foobar","long_output":"","state":2,"connector_name":"foobar","output":"foo","command_name":"foo","perf_data":"","max_attempts":2}
```

#### `break`

Lorsqu'une règle de type `break` est appliquée à un évènement, cet évènement sort de l'event-filter. Les règles suivantes ne sont pas appliquées, et l'évènement est traité par Canopsis.

#### `enrichment`

Les règles de types `enrichment` sont des règles d'enrichissement, qui permettent d'appliquer des actions modifiant les évènements.

Ces règles peuvent avoir les paramètres supplémentaires suivants :

 - `actions` (requis) : une liste d'actions à appliquer à l'évènement (voir [Actions](#actions) pour plus de détails).
 - `external_data` (optionnel) : des sources de données externes (voir [Données externes](#donnees-externes) pour plus de détails).
 - `on_success` (optionnel, `pass` par défaut) : le résultat de la règle en cas de succès (`pass`, `break` ou `drop`).
 - `on_failure` (optionnel, `pass` par défaut) : le résultat de la règle en cas d'échec (`pass`, `break` ou `drop`).

Lorsqu'une règle d'enrichissement est appliquée, les données externes sont récupérées, puis les actions sont exécutées, dans l'ordre dans lequel elles ont été définies. Si la récupération de données ou l'exécution d'une action échoue, l'application de la règle est interrompue, et son résultat sera la valeur de `on_failure`. Sinon, son résultat est la valeur de `on_success`.

Si le résultat de la règle est `drop`, l'évènement est supprimé. Les règles suivantes ne sont pas appliquées à cet évènement, et il est ignoré par Canopsis.

Si le résultat de la règle est `break`, l'évènement sort de l'event-filter. Les règles suivantes ne sont pas appliquées.

Si le résultat de la règle est `pass`, l'exécution de l'event-filter continue.

### Actions

Une action est un objet JSON contenant un champ `type` indiquant le type de l'action, et des paramètres. Les actions disponibles sont les suivantes:

  - `set_field` : Définir un champ d'un événement à une constante.
  - `set_field_from_template` : Définir un champ de type chaîne de caractères d'un événement en utilisant un modèle.
  - `set_entity_info` : Définir une info d'une entité à une constante.
  - `set_entity_info_from_template` : Définir une info de type chaîne de caractères d'une entité en utilisant un modèle.
  - `copy` : Copier une valeur d'un champ d'un événement à un autre.
  - `copy_to_entity_info` : Copier une valeur d'un champ d'un événement vers une info d'une entité.

**Note :** Les actions utilisent la représentation interne à Canopsis des évènements. Voir [Champs des évènements](#champs-des-evenements) pour la correspondance entre les noms des champs des évènements en JSON et dans la représentation de Canopsis.

#### Action `set_field`

L'action `set_field` permet de modifier un champ de l'évènement.

Les paramètres de l'action sont :

 - `name` (requis) : le nom du champ (tel que [défini dans l'annexe](#champs-des-evenements)).
 - `value` (requis) : la nouvelle valeur du champ.

Par exemple, l'action suivante passe la criticité d'un évènement en critique :

```json
{
    "type": "set_field",
    "name": "State",
    "value": 3
}
```

#### Action `set_field_from_template`

L'action `set_field_from_template` permet de modifier un champ de l'évènement à partir un modèle.

Les paramètres de l'action sont :

 - `name` (requis) : le nom du champ (tel que [défini dans l'annexe](#champs-des-evenements)).
 - `value` (requis) : le modèle utilisé pour déterminer la valeur du champ.

L'event-filter utilise le [moteur de modèles Go](https://golang.org/pkg/text/template/). Les champs de l'évènement peuvent être utilisés dans les modèles de la manière suivante : `{{.Event.NomDuChamp}}`. Il est également possible d'utiliser les expressions régulières des patterns pour utiliser des sous-groupes dans les modèles (voir [Expressions régulières](#expressions-regulieres) pour plus de détails), ou d'utiliser des [données externes](#donnees-externes).

Par exemple, l'action suivante modifie l'output d'un évènement pour y ajouter son auteur :

```json
{
    "type": "set_field_from_template",
    "name": "Output",
    "value": "{{.Event.Output}} (by {{.Event.Author}})"
}
```
#### Action `set_entity_info`

Cette action permet de définir une info d'entité de n'importe quel type.

Les paramètres de l'action sont :

  - `name` (requis) : le nom de l'information.
  - `description` (optionnel) : la description de l'information.
  - `value` (requis) : la nouvelle valeur du champ. 

Par exemple, l'action suivante définit le nom du client :

``` json
{
    "type": "set_entity_info",
    "name": "customer",
    "description": "Client",
    "value": "StarkTelecom"
}
```

#### Action `set_entity_info_from_template`

L'action `set_entity_info_from_template` permet de modifier une information de l'entité correspondant à l'évènement.

Les paramètres de l'action sont :

 - `name` (requis) : le nom de l'information.
 - `description` (optionnel) : la description de l'information.
 - `value` (requis) : le modèle utilisé pour déterminer la valeur de l'information.

L'event-filter utilise le [moteur de modèles go](https://golang.org/pkg/text/template/). Les champs de l'évènement peuvent être utilisés dans les modèles de la manière suivante : `{{.Event.NomDuChamp}}`. Il est également possible d'utiliser les expressions régulières des patterns pour utiliser des sous-groupes dans les modèles (voir [Expressions régulières](#expressions-regulieres) pour plus de détails), ou d'utiliser des [données externes](#donnees-externes).

Par exemple, l'action suivante modifie l'information `customer` d'une entité pour y ajouter la localisation :

```json
{
    "type": "set_entity_info_from_template",
    "name": "customer",
    "description": "Site client",
    "value": "StarkTelecom {{.Event.Location}}"
}
```

Cette action échoue si l'entité n'a pas été ajoutée à l'évènement au préalable. Pour utiliser cette action, il est donc nécessaire de définir une règle [ajoutant les entités aux évènements](#ajout-de-lentite-correspondant-a-un-evenement), avec une priorité inférieure à celles des règles contenant des actions de type `set_entity_info_from_template`.

#### Action `copy`

L'action `copy` permet de copier la valeur d'un champ dans un évènement.

Les paramètres de l'action sont :

 - `from` (requis) : le nom du champ dont la valeur doit être copiée. Il peut s'agir d'un champ de l'évènement (`Event.NomDuChamp`), d'un sous-groupe d'une expression régulière (voir [Expressions régulières](#expressions-regulieres)), ou d'une donnée externe (voir [Données externes](#donnees-externes)).
 - `to` (requis) : le nom du champ de l'évènement dans lequel la valeur doit être
   copiée.

Par exemple, l'action suivante va vérifier si l'entité à l'origine de l'évènement existe déjà dans le référentiel interne de Canopsis. Si c'est le cas elle sera copiée dans le champ `Entity` de l'évènement (reportez-vous à cet [exemple](#ajout-dinformations-a-lentite) pour l'utilisation du champ entity de l'évènement) :

```json
{
    "type": "copy",
    "from": "ExternalData.entity",
    "to": "Entity"
}
```

#### Action `copy_to_entity_info`

Cette action permet de copier une valeur de n'importe quel type dans les infos de l'entité.

Les paramètres de l'action sont :

 - `from` (requis) : le nom du champ dont la valeur doit être copiée. Il peut s'agir d'un champ de l'évènement (`Event.NomDuChamp`), d'un sous-groupe d'une expression régulière (voir [Expressions régulières](#expressions-regulieres)), ou d'une donnée externe (voir [Données externes](#donnees-externes)).
 - `name` (requis) : le nom de l'information.
 - `description` (optionnel) : la description de l'information.

Par exemple, l'action suivante va copier un champ environnement dans les infos de l'entité :

``` json
{
  "type" : "copy_to_entity_info",
  "from" : "Event.ExtraInfos.env",
  "name" : "env",
  "description" : "Environnement"
}
```

### Expressions régulières

Si l'un des patterns d'une règle contient une expression régulière (avec l'opérateur
`regex_match`) et des sous-groupes nommés, les valeurs de ces
sous-groupes peuvent être utilisés dans les templates des actions de type
`set_field_from_template` et `set_entity_info_from_template`, et comme champ
`from` des actions de type `copy`.

Pour nommer un groupe `([\.0-9]+)`, la syntaxe est la suivante :
```
(?P<nom_du_match>[\.0-9]+)
```

La valeur du sous-groupe sera alors disponible dans
`{{.RegexMatch.<nom du champ>.nom_du_match}}` pour les templates, et dans
`RegexMatch.<nom du champ>.nom_du_match` pour les actions de type `copy`.

Par exemple, si la liste des patterns vaut :

```json
"patterns": [
    {
        "State": {">=": 2},
        "Output": {"regex_match": "Warning: CPU Load is critical \((?P<load>.*)%\)"}
    }
]
```

et si l'output de l'évènement vaut `Warning: CPU Load is critical (97.5%)`,
alors il est possible d'utiliser l'expression `{{.RegexMatch.Output.load}}`,
qui vaudra `97.5`, dans un template.

### Données externes

Le champ `external_data` est un objet JSON contenant des couples `<nom de la
donnée>: <source de données>`.

Lors de l'application d'une règle d'enrichissement à un évènement, les sources
de données sont utilisées pour récupérer les données correspondant à cet
évènement. Ces données sont alors disponibles dans `{{.ExternalData.<nom de la
donnée>}}` pour les templates, et dans `ExternalData.<nom de la donnée>` pour
les actions de type `copy`.

Une source de données est un objet JSON contenant un champ `type` indiquant le
type de la source, et des paramètres. Les différents types de source de données
sont précisés ci-dessous.

Dans une installation Canopsis Pro, l'acquisition de données externes peut aussi
se faire avec un plugin `datasource`. Ceci nécessite [l'activation des plugins
datasource](moteur-che.md#activation-des-plugins-denrichissement-externe-datasource).

#### Entités

Une source de données de type `entity` renvoie l'entité correspondant à un
évènement. Elle ne prend pas de paramètre.

Voir [Ajout de l'entité correspondant à un
évènement](#ajout-de-lentite-correspondant-a-un-evenement) pour un exemple
d'utilisation de cette source de données.

#### Collection MongoDB externe

Les sources de données de type `mongo` permettent de récupérer des documents depuis
une collection de la base MongoDB de Canopsis. Ce type de source de données
n'est disponible que dans Canopsis Pro.

Les paramètres de cette source de données sont :

 - `collection` (requis) : le nom de la collection MongoDB à partir de laquelle
   les données seront importées. Il s'agit de votre propre collection, que vous
   devez créer et peupler vous-même avec les données voulues.
 - `select` (requis) : un objet JSON utilisé pour sélectionner un document de
   la collection. Cet objet contient des couples `<champ>: <template>`, où
   `<champ>` est le nom d'un champ de la collection, et `template` un template
   utilisé pour déterminer la valeur du champ.

L'event-filter utilise le [moteur de templates
Go](https://golang.org/pkg/text/template/). Les champs de l'évènement peuvent
être utilisés dans les templates de la manière suivante :
`{{.Event.NomDuChamp}}`. Il est également possible d'utiliser les expressions
régulières des patterns pour utiliser des sous-groupes dans les templates (voir
[Expressions régulières](#expressions-regulieres) pour plus de détails).

Par exemple, la source de données ci-dessous permet de récupérer un document
dont l'id est le nom du composant d'un évènement depuis une collection
`components` :

```json
{
    "type" : "mongo",
    "select" : {
        "_id" : "{{.Event.Component}}"
    },
    "collection" : "components"
}
```

Voir [Ajout d'informations à l'entité](#ajout-dinformations-a-lentite) pour un
exemple de règle utilisant cette source de données.

**Note :** Chaque utilisation d'une source de données de type `mongo` effectue
une requête MongoDB, ce qui risque d'affecter les performances du moteur `engine-che`.
Il est donc déconseillé de les utiliser dans des règles appliquées à tous les
évènements.

### Cas particulier des méta alarmes

Lorsqu'une méta alarme est créée ou mise à jour, un événement de type `metaalarm` ou `metaalarmupdated` est généré.

Cet événement met à disposition du moteur d'enrichissement plusieurs attributs :

| Représentation interne de Canopsis | Notes                                                       |
| ---------------------------------- | ----------------------------------------------------------- |
| .Event.ExtraInfos.Meta.Count       |  Nombre d'alarmes conséquences   |
| .Event.ExtraInfos.Meta.Children    |  Objet représentant la dernière alarme conséquence attachée |
| .Event.ExtraInfos.Meta.Rule        |  Les informations de la règle méta en elle-même            |

Voici un exemple qui permet d'ajouter un attribut texte sur l'entité de la méta alarme et dont le contenu vaut :

**`Count : Nombre d'alarmes conséquences; Children : Sévérité de la dernière alarme conséquence attachée; Rule : Nom de la règle ayant permis le regroupement`**

```
{
  "type": "set_entity_info_from_template",
  "name": "children",
  "description": "children",
  "value": "Count: {{ .Event.ExtraInfos.Meta.Count }}; Children: {{ .Event.ExtraInfos.Meta.Children.alarm.v.state.m }}; Rule: {{ .Event.ExtraInfos.Meta.Rule.Name }}"
}
```

## En cas de problème

Les logs du conteneur Docker `che` peuvent contenir des informations sur
l'exécution de l'event-filter.

### Délai d'application des règles

Les règles de l'event-filter sont rechargées toutes les minutes, il peut donc y
avoir un délai d'une minute entre la création d'une règle et son application.

### Règles invalides

Lors du chargement des règles, les règles dont le format est invalide sont
ignorées. Un message d'erreur est écrit dans les logs, contenant la règle et la
raison pour laquelle elle n'est pas valide.

### Mode debug

Pour tester le fonctionnement de l'event-filter, il est possible d'envoyer des
évènements en mode debug, permettant de tracer l'exécution de l'event-filter
sur cet évènement. Pour cela, il faut définir le champ `debug` de l'évènement à
`true`.

La trace de l'exécution de l'event-filter est alors enregistrée dans les logs.

## Exemples

### Conversion d'une adresse IP en nom de domaine

La règle suivante remplace le composant `192.168.0.1` par `example.com` dans
les évènements.

```javascript
{
    "description": "Conversion de 192.168.0.1 en example.com",
    "type": "enrichment",
    "patterns": [
        {"component": "192.168.0.1"}
    ],
    "actions": [
        {
            "type": "set_field",
            "name": "Component",
            "value": "example.com"
        }
    ],
    "priority": 100,
    "on_success": "pass",
    "on_failure": "pass"
}
```

### Traduction

Les règles suivantes permettent de traduire des messages de l'anglais vers le
français.

Puisque le champ `on_success` de ces deux règles vaut `break`, l'évènement sort
de l'event-filter dès que son output a été traduit.

```json
{
    "description": "Traduction du message de CPU critique",
    "type": "enrichment",
    "patterns": [
        {"output": {"regex_match": "Warning: CPU Load is critical \\((?P<load>.*)%\\)"}}
    ],
    "actions": [
        {
            "type": "set_field_from_template",
            "name": "Output",
            "value": "Attention, la charge CPU est critique ({{.RegexMatch.Output.load}}%)"
        }
    ],
    "priority": 100,
    "on_success": "break",
    "on_failure": "pass"
}
```

```json
{
    "description": "Traduction du message de disque presque plein",
    "type": "enrichment",
    "patterns": [
        {"output": {"regex_match": "Warning: The disk (?P<disk>.*) is almost full \\((?P<load>[\\.0-9]*)% used\\)"}}
    ],
    "actions": [
        {
            "type": "set_field_from_template",
            "name": "Output",
            "value": "Attention, le disque {{.RegexMatch.Output.disk}} est presque plein ({{.RegexMatch.Output.load}}% utilisés)"
        }
    ],
    "priority": 101,
    "on_success": "break",
    "on_failure": "pass"
}
```

Si l'on ne souhaite pas traduire l'output des évènements d'un composant, on
peut ajouter une règle de type `break` qui fait sortir ces évènements de
l'event-filter, avec une priorité inférieure à celle des règles de traduction
(pour qu'elle soit appliquée avant).

```json
{
    "type": "break",
    "patterns": [
        {"component": "component_name"}
    ],
    "priority": 50
}
```

### Passer des évènements en mode debug

La règle suivante permet de passer les évènements d'un composant en mode debug.

```json
{
    "type": "enrichment",
    "patterns": [
        {"component": "component_name"}
    ],
    "actions": [
        {
            "type": "set_field",
            "name": "Debug",
            "value": true
        }
    ],
    "priority": 0,
    "on_success": "pass",
    "on_failure": "pass"
}
```

### Ajout de l'entité correspondant à un évènement

La règle suivante permet d'ajouter l'entité correspondant à un évènement dans
son champ `Entity`.

```json
{
    "type": "enrichment",
    "patterns": [],
    "external_data": {
        "entity": {
            "type": "entity"
        }
    },
    "actions": [
        {
            "type": "copy",
            "from": "ExternalData.entity",
            "to": "Entity"
        }
    ],
    "on_success": "pass",
    "on_failure": "pass",
    "priority": 100
}
```

### Ajout d'informations à l'entité

Les règles suivantes permettent d'ajouter des informations à l'entité
correspondant à un évènement.

La première règle ajoute l'entité correspondant à l'évènement dans le champ
`Entity`, comme dans l'exemple précédent.

```json
{
    "type": "enrichment",
    "patterns": [],
    "external_data": {
        "entity": {
            "type": "entity"
        }
    },
    "actions": [
        {
            "type": "copy",
            "from": "ExternalData.entity",
            "to": "Entity"
        }
    ],
    "on_success": "pass",
    "on_failure": "pass",
    "priority": 100
}
```

La seconde règle utilise une source de données externes pour récupérer le
client et le responsable d'un composant dans une collection MongoDB
`components`, contenant des documents de la forme :

```json
{
    "_id": "nom_du_composant",
    "customer": "nom_du_client",
    "manager": "nom_du_responsable"
}
```

 - Cette règle n'est appliquée que si les informations `customer` et `manager`
   ne sont pas déjà définies, pour éviter de réaliser une requête MongoDB à
   chaque évènement.
 - Elle a une priorité supérieure à la règle ajoutant l'entité, pour que les
   actions de type `set_entity_info_from_template` n'échouent pas.

```json
{
    "type" : "enrichment",
    "patterns" : [{
        "current_entity": {
            "infos": {
                "customer": null,
                "manager": null
            }
        }
    }],
    "external_data" : {
        "component" : {
            "type" : "mongo",
            "select" : {
                "_id" : "{{.Event.Component}}"
            },
            "collection" : "components"
        }
    },
    "actions" : [
        {
            "type" : "set_entity_info_from_template",
            "name" : "customer",
            "description" : "Client",
            "value" : "{{.ExternalData.component.customer}}"
        },
        {
            "type" : "set_entity_info_from_template",
            "name" : "manager",
            "description" : "Responsable du service",
            "value" : "{{.ExternalData.component.manager}}"
        }
    ],
    "priority" : 200,
    "on_success" : "pass",
    "on_failure" : "pass"
}
```

La troisième règle utilise les informations de l'entité pour modifier l'output
des évènements.

```json
{
    "type" : "enrichment",
    "patterns" : [],
    "actions" : [
        {
            "type" : "set_field_from_template",
            "name" : "Output",
            "value" : "{{.Event.Output}} (client: {{.Event.Entity.Infos.customer.Value}})"
        }
    ],
    "priority" : 205,
    "on_success" : "pass",
    "on_failure" : "pass"
}
```

## Annexe

### Champs des évènements

| Nom du champ JSON | Nom de la représentation interne de Canopsis |
|-------------------|----------------------------------------------|
| `author`          | `Author` |
| `debug`           | `Debug` |
| `component`       | `Component` |
| `connector`       | `Connector` |
| `connector_name`  | `ConnectorName` |
| `current_entity`  | `Entity` |
| `event_type`      | `EventType` |
| `long_output`     | `LongOutput` |
| `output`          | `Output` |
| `perf_data`       | `PerfData` |
| `resource`        | `Resource` |
| `state`           | `State` |
| `state_type`      | `StateType` |
| `status`          | `Status` |
| `source_type`     | `SourceType` |
| `ticket`          | `Ticket` |
| `timestamp`       | `Timestamp` |


**Note :** le champ `Entity` n'est pas défini au début de l'exécution de l'event-filter. Pour y accéder, ou pour modifier les informations de l'entité, il faut utiliser une [règle ajoutant les entités aux évènements](#ajout-de-lentite-correspondant-a-un-evenement).
