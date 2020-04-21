# Import de contexte graphe

Cette API permet d'importer un référentiel externe (autrement appelé *contexte graphe* ou *context graph*) au format JSON dans Canopsis.

Un context graph est constitué d'un ensemble d'entités et de leurs relations. Les entités vont enrichir les alarmes en leur rajoutant du contexte.

Par exemple, pour un équipement envoyant une alarme dans Canopsis, si une entité représentant cet équipement est présente, alors l'alarme sera enrichie avec toutes les informations que contient l'entité. Ces informations peuvent être de toute sorte : adresse IP, OS installé, date de mise en service, criticité de l'équipement, responsable à contacter en cas de panne, etc.

L'API propose différentes routes pour déclencher un import ou pour suivre son évolution.

## Format d'import du context graph

Un référentiel externe est constitué de deux types de données : les entités et les liens entre ces entités.

Au moment de l'import, on va donc séparer ces deux types d'éléments dans un objet JSON. On va retrouver les entités dans le champ `cis` (venant de *Configuration ItemS*) et les relations entre les entités dans `links`.

Dans chaque entité présente dans la liste `cis` et dans chaque lien présent dans la liste `links`, le champ `action` va définir comment l'élément sera importé dans le contexte graphe.

```json
{
    "cis": [
        {"_id": "string",
         "type": "string",
         "infos": {},
         "measurements": [{}],
         "name": "string",
         "action": "string",
         "action_properties": {"disable": "string"},
        }
    ],
    "links": [
        {"_id": "string",
         "to": "string",
         "infos": {},
         "action": "string",
         "action_properties": {}
        }
    ]
}
```

## Entités

### Description des champs d'entités

La liste `cis` est une liste d'entités représentée sous forme d'objet JSON. Chaque entité possède ces différents champs :

* **`_id`** (chaîne de caractères) : identifiant de l'entité concernée par l'action.
* **`type`** (chaîne de caractères) : type de l'entité. Valeurs possibles : `resource`, `component`, `connector`, ou `watcher`.
* **`infos`** : informations complémentaires. Ce sont des données totalement personnalisables que l'utilisateur peut modifier.
* **`measurements`** : ensemble de chaînes de caractères correspondant aux métriques liées à l'entité.
* **`name`** (chaîne de caractères) : nom de l'entité.
* **`action`** (chaîne de caractères) : type de l'action à réaliser au moment de l'import.
* **`action_properties`** : vient en complément du champ `action` en spécifiant des informations complémentaires pour la bonne réalisation de l'action.

### Description des actions sur les entités

Les actions suivantes peuvent être réalisées sur les entités, au moment de l'import d'un contexte graphe : `create`, `set`, `delete`, `update`, `disable` et `enable`.

#### Create

`create` crée une nouvelle entité dans le contexte graphe. Si une entité existe déjà avec le même identifiant, une **mise à jour complète** de l'entité sera effectuée. Les nouvelles données vont écraser l'entité déjà présente en base de données.

Pour fonctionner correctement, les champs `_id`, `type` et `action` sont obligatoires.

Les champs de la section [Description des champs d'entités](#Description-des-champs-d'entités) vont être copiés tels quels dans la nouvelle entité, à l'exception d'`action` et d'`action_properties`. Si le champ n'existe pas, le champ dans la nouvelle entité sera initialisé avec une valeur par défaut :

* `name` sera initialisé avec la valeur du champ `_id`
* `depends` et `impact` seront initialisés avec une liste vide
* `measurements` sera initialisé avec une liste vide
* et le champ `infos` par un objet vide.

#### Set

`set`, comme `create`, crée une nouvelle entité dans le contexte graphe. La différence est que l'action `set` procède à une **mise à jour partielle** de l'entité. Seuls les champs fournis dans les données d'import impacteront l'entité.

Par exemple, si une entité existante possède `info1` et `info2` dans ses `infos` et que l'import contient `info1` avec une valeur différente et `info3`, alors `info1` sera mis à jour avec sa nouvelle valeur, `info2` ne changera pas et `info3` sera créé.

Un exemple pratique dans le paragraphe suivant permettra de mieux saisir la différente entre `set` et `create`.

## Différence entre l'action create et l'action set

Imaginons une nouvelle entité `france` qui possède au départ deux informations (champ `infos` dans l'objet JSON) : `capitale` et `metropole1`.

```json
{
    "_id" : "france",
    "infos" : {
        "capitale" : {
            "name" : "capitale",
            "value" : "Poitiers",
            "description" : "Capitale de la France"
        },
        "metropole1" : {
            "name" : "metropole1",
            "value" : "Lyon",
            "description" : "Agglomération française"
        }
    },
    "action":""
}
```

On importe cette entité en base, mais on se rend compte trop tard que la `capitale` est incorrecte. On décide donc de mettre à jour l'entité avec `capitale` et d'en profiter pour ajouter le champ supplémentaire `metropole2`, qui vaut `Marseille`.

Le deuxième appel à l'API d'import aura cette forme.

```json
{
    "_id" : "france",
    "infos" : {
        "capitale" : {
            "name" : "capitale",
            "value" : "Paris",
            "description" : "Capitale corrigée de la France"
        },
        "metropole2" : {
            "name" : "metropole2",
            "value" : "Marseille",
            "description" : "Une autre agglomération française"
        }
    },
    "action":""
}
```

Avec le premier appel à l'API d'import (`{"capitale":"Poitiers", "metropole1":"Lyon"}`), les actions `create` et `set` sont identiques : elles vont créer en base de données l'entité `france`. La différence se fait lors du second appel à l'API (`{"capitale":"Paris", "metropole2":"Marseille"}`).

L'action `create` procède à une **mise à jour complète** de l'entité, c'est-à-dire les données précédentes vont être écrasées. En base de données, on se retrouve donc avec la `capitale` qui vaut `Paris`, `metropole1` qui a disparu et `metropole2` qui est présent.

L'action `set` est quant à elle une **mise à jour partielle**, en fournissant seulement les champs avec de nouvelles données, ce qui conserve le reste de l'entité. On retrouve donc dans l'entité la `capitale` qui vaut `Paris`, `metropole1` qui vaut toujours `Lyon` et qui n'a pas été impacté par le nouvel import et `metropole2` qui vaut `Marseille`.

```json
{
    "_id" : "france",
    "infos" : {
        "capitale" : {
            "name" : "capitale",
            "value" : "Paris",
            "description" : "Capitale corrigée de la France"
        },
        "metropole1" : {
            "name" : "metropole1",
            "value" : "Lyon",
            "description" : "Agglomération française"
        },
        "metropole2" : {
            "name" : "metropole2",
            "value" : "Marseille",
            "description" : "Une autre agglomération française"
        }
    }
}
```

#### Delete

`delete` supprime une entité désignée par son identifiant (champ `_id`). Si l'entité n'existe pas dans le contexte graphe, une erreur est déclenchée, l'import s'arrêtera et les modifications de l'import en cours ne seront pas répercutées sur le référentiel de Canopsis.

Pour fonctionner correctement, les champs `_id` et `action` sont obligatoires.

Lors de la suppression de l'entité, les champs `impact` et `depends` des entités ayant une relation avec celle supprimée verront les références à l'entité détruite supprimées.

#### Update

`update` met à jour l'entité désignée par son `_id`. Si l'entité désignée par `_id` n'existe pas initialement dans le contexte graphe, une erreur sera déclenchée, l'import s'arrêtera et les modifications ne seront pas répercutées dans le référentiel.

Les champs `_id` et `action` sont obligatoires.

Pour la mise à jour, tous les champs présents dans l'action, hormis `action` et `action_properties`, seront recopiés dans l'entité à mettre à jour.

#### Enable

L'action `enable` va activer une entité, c'est-à-dire que le champ `enabled` de l'entité sera mis à `True`. Si l'entité n'existe pas dans le contexte graphe, une erreur sera déclenchée, l'import s'arrêtera et les modifications ne seront pas répercutées dans le contexte graphe.

Les champs `_id`, `action` et `action_properties` sont obligatoires. Le champ `action_properties` doit avoir un champ `enable` contenant soit un entier représentant un timestamp, soit une liste d'horodatage correspondant à l'activation de l'entité.

#### Disable

`disable` réalise l'action contraire du `enable`, à savoir désactiver une entité.

Les champs `_id`, `action` et `action_properties` sont obligatoires. Le champ `action_properties` doit avoir un champ `disable` correspondant à la désactivation de l'entité.

## Liens

### Description des champs des liens

La liste `links` représente toutes les relations entre les entités. Ces liens sont stockés sous forme d'objet JSON avec ces différents champs :

* **`_id`** (chaîne de caractères) : identifiant de l'action, il peut prendre n'importe quelle valeur.
* **`from`** (chaîne de caractères) : identifiant de l'entité de départ du lien.
* **`to`** (chaîne de caractères) : identifiant de l'entité d'arrivée du lien.
* **`action`** (chaîne de caractères) : type de l'action à réaliser.

Pour tous les liens, les champs `from`, `to` et `action` sont obligatoires.

Les liens décrits dans les actions sont des liens de type *impact-depends*, c'est-à-dire qu'ils représentent des liens dont l'entité de départ (champ `from` dans le lien) contient dans son champ `impact` une référence à l'entité d'arrivée du lien. Par conséquent, l'entité d'arrivée du lien (champ `to`) contiendra une référence à l'entité de départ dans son champ `depends`.

### Description des actions sur les liens

Deux actions sont possibles sur les liens : `create` et `delete`.

#### Create

`create` crée un lien entre les deux entités définies à l'aide des champs `to` et `from` selon les modalités définies ci-dessus. Si au moins une des deux entités n'existe pas dans le contexte graphe, une erreur est déclenchée, l'import s'arrêtera et les modifications de l'import en cours ne seront pas répercutées sur le référentiel de Canopsis.

#### Delete

`delete` supprime un lien entre les deux entités définies à l'aide des champs `to` et `from` selon les modalités définies ci-dessus. Si au moins une des deux entités n'existe pas dans le contexte graphe, une erreur est déclenchée, l'import s'arrêtera et les modifications de l'import en cours ne seront pas répercutées sur le référentiel de Canopsis.

## Fonctionnement de l'import

Pour importer un référentiel, il faut envoyer les données au format JSON avec `PUT /api/contextgraph/import`. Le serveur retournera un objet JSON contenant l'identifiant de l'import (`import_id`). Sur le serveur, l'import sera stocké sous la forme d'un fichier, afin que la *task* responsable s'occupe de l'import.

Lors du traitement de l'import par la tâche, si le fichier est mal formé ou qu'une anomalie est survenue lors du traitement d'une action, l'import est annulé et les modifications ne sont pas répercutées sur le contexte graphe. Si aucune erreur n'est rencontrée, le contexte graphe sera mis à jour.

Dans tous les cas, la progression et le résultat de l'import sont disponibles via la route `GET /api/contextgraph/import/status/<import_id>`. La route retourne un objet JSON contenant au moins les champs `_id` (identifiant de l'import), `creation` (date à laquelle l'import a été envoyé à Canopsis) et `status`. Ce `status` est le statut de l'import et peut prendre quatre valeurs : `pending` (en attente), `ongoing` (en cours), `failed` (échec de l'import), `done` (succès de l'import).

Dans le cas où l'import est `failed` ou `done`, la route `GET /api/contextgraph/import/status/<import_id>` retourne des informations supplémentaires : temps d'exécution, nombre d'éléments impactés ou raison de l'échec.

Différents exemples de retour des routes `PUT /api/contextgraph/import` et `GET /api/contextgraph/import/status/<import_id>` sont disponibles dans le paragraphe suivant.

## Utilisation de l'API d'import

### Import d'un context graph

**URL** : `/api/contextgraph/import`

**Méthode** : `PUT`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :

```json
{
    "json":{
        "cis":[
            {
                "name":"capitals",
                "enabled":true,
                "action":"create",
                "infos":{
                    "info1":{
                        "name":"info1",
                        "value":"Paris",
                        "description":"Capitale de la France"
                    },
                    "info2":{
                        "name":"info2",
                        "value":"Londres",
                        "description":"Capitale de la Grande-Bretagne"
                    }
                },
                "_id":"capitals",
                "type":"component"
            }
        ],
        "links":[]
    }
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X PUT -u root:root -H "Content-Type: application/json" -d '{
    "json":{
        "cis":[
            {
                "name":"capitals",
                "enabled":true,
                "action":"create",
                "infos":{
                    "info1":{
                        "name":"info1",
                        "value":"Paris",
                        "description":"Capitale de la France"
                    },
                    "info2":{
                        "name":"info2",
                        "value":"Londres",
                        "description":"Capitale de la Grande-Bretagne"
                    }
                },
                "_id":"capitals",
                "type":"component"
            }
        ],
        "links":[]
    }
}' 'http://<Canopsis_URL>/api/contextgraph/import'
```

#### Réponse en cas de réussite

**Condition** : l'import a bien été réalisé dans Canopsis

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{
    "total":1,
    "data":[
        {
            "import_id":"c3090ed6-5b17-4c75-ad23-82238cffa62f"
        }
    ],
    "success":true
}
```

### Statut d'un import envoyé à Canopsis

**URL** : `/api/contextgraph/import/status/<import_id>`

**Méthode** : `GET`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut connaître le statut de la task `c3090ed6-5b17-4c75-ad23-82238cffa62f` : `curl -u root:root http://<Canopsis_URL>/api/contextgraph/import/status/c3090ed6-5b17-4c75-ad23-82238cffa62f`

#### Import en attente

Lorsqu'un import est déclenché alors qu'un autre est en cours d'exécution, le nouvel import est mis en attente. L'interrogation de la route précédemment citée avec l'identifiant du nouvel import retournera un objet JSON dont le `status` sera à **`pending`**.
```json
{
  "status": "pending",
  "_id": "c3090ed6-5b17-4c75-ad23-82238cffa62f",
  "creation": "Mon Aug 28 17:41:27 2017"
}
```

#### Import en cours

Lorsqu'un import est en cours de traitement, l'interrogation de la route précédemment citée avec l'identifiant du nouvel import retournera un objet JSON dont le `status` sera à **`ongoing`**.
```json
{
  "status": "ongoing",
  "_id": "c3090ed6-5b17-4c75-ad23-82238cffa62f",
  "creation": "Mon Aug 28 17:41:27 2017"
}
```

#### Import terminé

Lorsqu'un import est traité complètement sans erreur, la route précédemment citée retourne un objet JSON avec le `status` à **`done`**, le nombre d'entités supprimées dans `stats.deleted`, le nombre d'entités mises à jour ou créées `stats.updated` et le temps d'exécution de l'import dans le champ `exec_time`.
```json
{
  "status": "done",
  "exec_time": "00:03:58",
  "_id": "c3090ed6-5b17-4c75-ad23-82238cffa62f",
  "creation": "Mon Aug 28 17:41:27 2017",
  "stats": {
    "deleted": 0,
    "updated": 20000
  }
}
```

### Échec de l'import

En cas d'erreur, la route précédemment citée retourne un objet JSON contenant un status **`failed`**, une description de l'erreur dans le champ **`infos`** et le temps d'exécution total de l'import dans **`exec_time`**.
```json
{
  "status": "failed",
  "info": "ValueError(u'The ci of id connector_0 match an existing entity.',)",
  "_id": "c3090ed6-5b17-4c75-ad23-82238cffa62f",
  "creation": "Mon Aug 28 17:41:27 2017",
  "exec_time": "00:02:55"
}
```
