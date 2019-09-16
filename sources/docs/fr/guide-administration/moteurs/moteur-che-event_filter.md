# Che - Event-filter

!!! note
    Cette page concerne l'event-filter nouvelle génération, disponible uniquement
    avec le moteur GO `che`.

L'event-filter est une fonctionnalité du moteur [`che`](moteur-che.md) permettant de définir des règles manipulant les évènements.

Les règles sont définies dans la collection MongoDB `eventfilter`, et peuvent être ajoutées et modifiées avec l'[API event-filter](../../guide-developpement/event-filter/api_v2_event-filter.md).

Des exemples pratiques d'utilisation de l'event-filter sont disponibles dans la partie [Exemples](#exemples).

## Activation du plugin d'enrichissement depuis une ressource externe

!!! note
    Cette fonctionnalité n'est disponible que dans l'édition CAT de Canopsis.

L['event-filter](moteur-che-event_filter.md) peut utiliser des sources de données externes pour enrichir les évènements. Ces sources externes (à l'exception de `entity`) sont des [plugins](../../guide-developpement/plugins/event-filter-data-source.md) disponibles dans Canopsis CAT.

Les plugins doivent-être placés dans un dossier accessible par le moteur `che`.

### Activation avec Docker

Les plugins doivent-être ajoutés dans un volume dans l'image docker, et leur emplacement doit-être précisé dans la commande. Par exemple, avec `docker-compose` :

```yaml
  che:
    image: canopsis/engine-che-cat:${CANOPSIS_IMAGE_TAG}
    env_file:
      - compose.env
    restart: unless-stopped
    command: /engine-che -dataSourceDirectory /data-source-plugins
    volumes:
      - "./plugins:/data-source-plugins"
```

Dans une installation Docker, l'image `canopsis/engine-che-cat` remplace l'image par défaut `canopsis/engine-che`. Le moteur `che` doit ensuite être lancé au minimum avec l'option suivante pour que le plugin d'enrichissement externe soit chargé : `engine-che -dataSourceDirectory /data-source-plugins`

Les plugins doivent-être placés dans un dossier accessible par le moteur `che`.

L'exécutable `engine-che` accepte une option `-dataSourceDirectory` permettant de préciser le dossier contenant les plugins. Par défaut, ce dossier est celui contenant `engine-che`.

### Activation par paquets

Pour pouvoir utiliser l'enrichissement depuis une ressource externe, il faut :

*  lancer le moteur `che` avec l'option `-dataSourceDirectory <dossier contenant les plugins>`. Par défaut, ce dossier est celui contenant `engine-che`.


## Règles

Une règle est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel) : un identifiant unique (généré automatiquement s'il n'est pas défini par l'utilisateur).
 - `description` (optionnel) : une description de la règle de la règle, donnée par l'utilisateur.
 - `type` (requis) : le type de la règle (voir [Types de règles](#types-de-règles) pour plus de détails).
 - `pattern` (optionnel) : un pattern permettant de sélectionner les évènements auxquels la règle doit être appliquée. Si le pattern n'est pas précisé, la
   règle est appliquée à tous les évènements (voir [Patterns](#patterns) pour plus de détails).
 - `priority` (optionnel, 0 par défaut) : la priorité de la règle. Les règles sont appliquées par ordre de priorité croissante.
 - `enabled` (optionnel, `true` par défaut) : `false` pour désactiver la règle.

#### Exemple

```json
{
    "type": "drop",
    "pattern": {
        "resource": "invalid_resource"
    },
    "priority": 10
}
```

Le `type` de cette règle vaut `drop`, indiquant que c'est une règle qui supprime les évènements.

Le `pattern` de cette règle sélectionne les évènements dont la ressource vaut`invalid_resource` (voir [Patterns](#patterns) pour plus de détails).

Cette règle supprime donc les évènements dont la ressource vaut `invalid_resource`.

### Application des règles

Lors de la réception d'un évènement par le moteur `che`, les règles sont parcourues par ordre de priorité croissante. Si l'évènement est reconnu par le `pattern` d'une règle, celle-ci est appliquée à l'évènement. Le traitement effectué dépend du `type` de la règle (voir [Types de règles](#types-de-règles) pour plus de détails).

**Note :** Pour pouvoir être traité par Canopsis, un évènement doit respecter la condition suivante :

 - son champ `source_type` vaut `component`, et son champ `component` est défini et ne vaut pas `""`; *ou*
 - son champ `source_type` vaut `resource`, et ses champs `component` et `resource` sont définis et ne valent pas `""`.

Les évènements ne respectant pas cette condition en entrée du moteur `che` ou en sortie de l'event-filter sont supprimés.


### Patterns

Le pattern d'une règle permet de sélectionner les évènements auxquels elle doit être appliquée.

#### Patterns simples

Un pattern peut être défini comme un objet JSON contenant les valeurs de certains champs d'un évènement.

Par exemple, une règle contenant le pattern suivant sera appliquée aux évènements dont le composant vaut `component_name` et dont la ressource vaut `resource_name` :

```json
"pattern": {
    "component": "component_name",
    "resource": "resource_name"
}
```

#### Patterns avancés

Pour plus d'expressivité, il est possible d'associer à un champ un objet contenant des couples `operateur: valeur`. Les opérateurs disponibles sont :

 - `>=`, `>`, `<`, `<=` : compare une valeur numérique à une autre valeur.
 - `regex_match` : filtre la valeur d'une clé selon une expression régulière. La syntaxe des expressions régulières est celle de [go](https://golang.org/pkg/regexp/syntax/), et est similaire à celle acceptée par Perl et Python.

Par exemple, le pattern suivant sélectionne les évènements dont l'état est compris entre 1 et 3 (mineur, majeur ou critique) et dont l'output vérifie une expression régulière :

```json
"pattern": {
    "state": {">=": 1, "<=": 3},
    "output": {"regex_match": "Warning: CPU Load is critical \\(.*\\)"}
}
```

### Types de règles

#### Drop

Lorsqu'une règle de type `drop` est appliquée à un évènement, cet évènement est supprimé. Les règles suivantes ne sont pas appliquées à cet évènement, et il est ignoré par Canopsis.

Lorsqu'un évènement provoque le déclenchement d'une règle « drop », le moteur `che` l'affiche sur sa sortie de log :

```
2019/05/02 12:45:19 event dropped by event filter: {"hostgroups":["HG_FOOBAR"],"event_type":"check","execution_time":2.139087200164795,"timestamp":1556793914,"component":"foobar","state_type":0,"source_type":"resource","resource":"PING","current_attempt":1,"connector":"foobar","long_output":"","state":2,"connector_name":"foobar","output":"foo","command_name":"foo","perf_data":"","max_attempts":2}
```

#### Break

Lorsqu'une règle de type `break` est appliquée à un évènement, cet évènement sort de l'event-filter. Les règles suivantes ne sont pas appliquées, et l'évènement est traité par Canopsis.

#### Enrichment

Les règles de types `enrichment` sont des règles d'enrichissement, qui permettent d'appliquer des actions modifiant les évènements.

## Règles d'enrichissement

Les règles de types `enrichment` sont des règles d'enrichissement, qui permettent d'appliquer des actions modifiant les évènements.

Ces règles peuvent avoir les paramètres suivants (en plus de `type`, `pattern`, `priority` et `enabled`) :

 - `actions` (requis) : une liste d'actions à appliquer à l'évènement (voir [Actions](#actions) pour plus de détails).
 - `external_data` (optionnel) : des sources de données externes (voir [Données externes](#données-externes) pour plus de détails).
 - `on_success` (optionnel, `pass` par défaut) : le résultat de la règle en cas de succès (`pass`, `break` ou `drop`).
 - `on_failure` (optionnel, `pass` par défaut) : le résultat de la règle en cas d'échec (`pass`, `break` ou `drop`).

Lorsqu'une règle d'enrichissement est appliquée, les données externes sont récupérées, puis les règles sont appliquées, dans l'ordre dans lequel elles ont été définies. Si la récupération de données ou l'exécution d'une action échoue, l'application de la règle est interrompue, et son résultat sera la valeur de `on_failure`. Sinon, son résultat est la valeur de `on_success`.

Si le résultat de la règle est `drop`, l'évènement est supprimé. Les règles suivantes ne sont pas appliquées à cet évènement, et il est ignoré par Canopsis.

Si le résultat de la règle est `break`, l'évènement sort de l'event-filter. Les règles suivantes ne sont pas appliquées.

Si le résultat de la règle est `pass`, l'exécution de l'event-filter continue.

### Actions

Une action est un objet JSON contenant un champ `type` indiquant le type de l'action, et des paramètres. Les actions disponibles sont précisées ci-dessous.

**Note :** Les actions utilisent la représentation interne à Canopsis des évènements. Voir [Champs des évènements](#champs-des-évènements) pour la correspondance entre les noms des champs des évènements en JSON et dans la représentation de Canopsis.

#### `set_field`

L'action `set_field` permet de modifier un champ de l'évènement.

Les paramètres de l'action sont :

 - `name` (requis) : le nom du champ.
 - `value` (requis) : la nouvelle valeur du champ.

Par exemple, l'action suivante remplace l'état d'un évènement par un état critique :

```json
{
    "type": "set_field",
    "name": "State",
    "value": 3
}
```

#### `set_field_from_template`

L'action `set_field_from_template` permet de modifier un champ de l'évènement
avec un template.

Les paramètres de l'action sont :

 - `name` (requis) : le nom du champ.
 - `value` (requis) : le template utilisé pour déterminer la valeur du champ.

L'event-filter utilise le [moteur de templates
go](https://golang.org/pkg/text/template/). Les champs de l'évènement peuvent
être utilisés dans les templates de la manière suivante :
`{{.Event.NomDuChamp}}`. Il est également possible d'utiliser les expressions
régulières des patterns pour utiliser des sous-groupes dans les templates (voir
[Expressions régulières](#expressions-régulières) pour plus de détails), ou
d'utiliser des [données externes](#données-externes).

Par exemple, l'action suivante modifie l'output d'un évènement pour y ajouter
son auteur :

```json
{
    "type": "set_field_from_template",
    "name": "Output",
    "value": "{{.Event.Output}} (by {{.Event.Author}})"
}
```


#### `set_entity_info_from_template`

L'action `set_entity_info_from_template` permet de modifier une information de
l'entité correspondant à l'évènement.

Les paramètres de l'action sont :

 - `name` (requis) : le nom de l'information.
 - `description` (optionnel) : la description de l'information
 - `value` (requis) : le template utilisé pour déterminer la valeur de
   l'information.

L'event-filter utilise le [moteur de templates
go](https://golang.org/pkg/text/template/). Les champs de l'évènement peuvent
être utilisés dans les templates de la manière suivante :
`{{.Event.NomDuChamp}}`. Il est également possible d'utiliser les expressions
régulières des patterns pour utiliser des sous-groupes dans les templates (voir
[Expressions régulières](#expressions-régulières) pour plus de détails), ou
d'utiliser des [données externes](#données-externes).


Par exemple, l'action suivante modifie l'information `customer` d'une entité :

```json
{
    "type": "set_entity_info_from_template",
    "name": "customer",
    "description": "Client",
    "value": "jack"
}
```

Cette action échoue si l'entité n'a pas été ajoutée à l'évènement au préalable.
Pour utiliser cette action, il est donc nécessaire de définir une règle
[ajoutant les entités aux
évènements](#ajout-de-lentité-correspondant-à-un-évènement), avec une priorité
inférieure à celles des règles contenant des actions de type
`set_entity_info_from_template`.

#### `copy`

L'action `copy` permet de copier la valeur d'un champ dans un évènement.

Les paramètres de l'action sont :

 - `from` : le nom du champ dont la valeur doit être copiée. Il peut s'agir
   d'un champ de l'évènement (`Event.NomDuChamp`), d'un sous-groupe d'une
   expression régulière (voir [Expressions
   régulières](#expressions-régulières)), ou d'une donnée externe (voir
   [Données externes](#données-externes)).
 - `to` : le nom du champ de l'évènement dans lequel la valeur doit être
   copiée.

Par exemple, l'action suivante copie une entité dans le champ `Entity` d'un
évènement :

```json
{
    "type": "copy",
    "from": "ExternalData.entity",
    "to": "Entity"
}
```

### Expressions régulières

Si le pattern d'une règle contient une expression régulière (avec l'opérateur
`regex_match`) contenant des sous-groupes nommés, les valeurs de ces
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

Par exemple, si le pattern vaut :

```json
"pattern": {
    "State": {">=": 2},
    "Output": {"regex_match": "Warning: CPU Load is critical \\((?P<load>.*)%\\)"}
}
```

et si l'output de l'évènement vaut `Warning: CPU Load is critical (97.5%)`,
alors il est possible d'utiliser l'expression `{{.RegexMatch.Output.load}}`,
qui vaudra `97.5`, dans un template.

### Données externes

Le champ `external_data` est un objet JSON contenant des couples `<nom de la
données>: <source de données>`.

Lors de l'application d'une règle d'enrichissement à un évènement, les sources
de données sont utilisées pour récupérer les données correspondant à cet
évènement. Ces données sont alors disponibles dans `{{.ExternalData.<nom de la
données>}}` pour les templates, et dans `ExternalData.<nom de la données>` pour
les actions de type `copy`.

Une source de données est un objet JSON contenant un champ `type` indiquant le
type de la source, et des paramètres. Les différents types de source de données
sont précisés ci-dessous.

#### Entités

Une source de données de type `entity` renvoie l'entité correspondant à un
évènement. Elle ne prend pas de paramètres.

Voir [Ajout de l'entité correspondant à un
évènement](#ajout-de-lentité-correspondant-à-un-évènement) pour un exemple
d'utilisation de cette source de données.

#### Collection MongoDB

Les sources de données de type `mongo` permettent de récupérer des documents dans
une collection de la base MongoDB de Canopsis. Ce type de source de données
n'est disponible que dans Canopsis CAT.

Les paramètres de cette source de données sont :

 - `collection` (requis) : le nom de la collection MongoDB.
 - `select` (requis) : un objet JSON utilisé pour sélectionner un document de
   la collection. Cet objet contient des couples `<champ>: <template>`, où
   `<champ>` est le nom d'un champ de la collection, et `template` un template
   utilisé pour déterminer la valeur du champ.

L'event-filter utilise le [moteur de templates
go](https://golang.org/pkg/text/template/). Les champs de l'évènement peuvent
être utilisés dans les templates de la manière suivante :
`{{.Event.NomDuChamp}}`. Il est également possible d'utiliser les expressions
régulières des patterns pour utiliser des sous-groupes dans les templates (voir
[Expressions régulières](#expressions-régulières) pour plus de détails).

Par exemple, la source de données ci-dessous permet de récupérer un document
dont l'id est le nom du composant d'un évènement dans une collection
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

Voir [Ajout d'informations à l'entité](#ajout-dinformations-à-lentité) pour un
exemple de règle utilisant cette source de données.

**Note :** Chaque utilisation d'une source de données de type `mongo` effectue
une requête MongoDB, ce qui risque d'affecter les performances du moteur `che`.
Il est donc déconseillé de les utiliser dans des règles appliquées à tous les
évènements.


## En cas de problème

Les logs du container docker `che` peuvent contenir des informations sur
l'exécution de l'event-filter.

### Délai d'application des règles

Les règles de l'event-filter sont rechargées toutes les minutes, il peut donc y
avoir un délai d'une minute entre la création d'une règle et son application.

### Règles invalides

Lors du chargement des règles, les règles dont le format est invalide sont
ignorées. Un message d'erreur est écrit dans les logs, contenant la règle et la
raison pour laquelle elle n'est pas valide.

### Mode débug

Pour tester le fonctionnement de l'event-filter, il est possible d'envoyer des
évènements en mode débug, permettant de tracer l'exécution de l'event-filter
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
    "pattern": {
        "component": "192.168.0.1"
    },
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
    "pattern": {
        "output": {"regex_match": "Warning: CPU Load is critical \\((?P<load>.*)%\\)"}
    },
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
    "pattern": {
        "output": {"regex_match": "Warning: The disk (?P<disk>.*) is almost full \\((?P<load>[\\.0-9]*)% used\\)"}
    },
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
    "pattern": {
        "component": "component_name"
    },
    "priority": 50
}
```

### Passer des évènements en mode debug

La règle suivante permet de passer les évènements d'un composant en mode débug.

```json
{
    "type": "enrichment",
    "pattern": {
        "component": "component_name"
    },
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
    "pattern": {},
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
    "pattern": {},
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
   ne sont pas déjà définies, pour éviter d'effectuer une requête MongoDB à
   chaque évènement.
 - Elle a une priorité supérieure à la règle ajoutant l'entité, pour que les
   actions de type `set_entity_info_from_template` n'échouent pas.

```json
{
    "type" : "enrichment",
    "pattern" : {
        "current_entity": {
            "infos": {
                "customer": null,
                "manager": null
            }
        }
    },
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
    "pattern" : {},
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

| Évènement JSON | Représentation interne de Canopsis |                                                                                                                         Notes                                                                                                                         |
| -------------- | ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|   connector    |             Connector              |                                                                                                                                                                                                                                                       |
| connector_name |           ConnectorName            |                                                                                                                                                                                                                                                       |
|   event_type   |             EventType              |                                                                                                                                                                                                                                                       |
|   component    |             Component              |                                                                                                                                                                                                                                                       |
|    resource    |              Resource              |                                                                                                                                                                                                                                                       |
|   perf_data    |              PerfData              |                                                                                                                                                                                                                                                       |
|     status     |               Status               |                                                                                                                                                                                                                                                       |
|   timestamp    |             Timestamp              |                                                                                                                                                                                                                                                       |
|   state_type   |             StateType              |                                                                                                                                                                                                                                                       |
|  source_type   |             SourceType             |                                                                                                                                                                                                                                                       |
|  long_output   |             LongOutput             |                                                                                                                                                                                                                                                       |
|     state      |               State                |                                                                                                                                                                                                                                                       |
|     output     |               Output               |                                                                                                                                                                                                                                                       |
|     author     |               Author               |                                                                                                                                                                                                                                                       |
|     ticket     |               Ticket               |                                                                                                                                                                                                                                                       |
|     debug      |               Debug                |                                                                                                                                                                                                                                                       |
| current_entity |               Entity               | Ce champ n'est pas défini au début de l'exécution de l'event-filter. Pour y accéder, ou pour modifier les informations de l'entité, il faut utiliser une [règle ajoutant les entités aux évènements](#ajout-de-lentité-correspondant-à-un-évènement). |
