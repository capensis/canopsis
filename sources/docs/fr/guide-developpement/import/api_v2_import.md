# API d'import de contexte graphe

Cette API permet d'importer un référentiel externe (autrement appelé *contexte graphe* ou *context graph*) au format JSON dans Canopsis.

Un context graph est consitué d'un ensemble d'entités et des relations entre elles. Les entités vont enrichir les alarmes en leur rajoutant du contexte. Par exemple, si un équipement qui envoie une alarme dans Canopsis et qu'une entité représente cet équipement est présent, l'alarme sera entichie avec toutes les informations que contient l'entité. Ces informations peuvent être de toute sorte : adresse IP, OS installé, date de mise en service, criticité de l'équipement, responsable à contacter en cas de panne, etc.

L'API propose différentes routes pour téléverser un import ou suivre son évolution.

# Format d'import du context graph

Un référentiel externe est constitué de deux types de données : les entités et les liens entre ces entités.

Au moment de l'import, on va donc séparer ces deux types d'éléments dans un objet JSON. On va retrouver les entités dans le champ `cis` (venant de *Configuration Items*) et les relations entre les entités dans `links`.

Dans chaque entité présent dans la liste `cis` et dans chaque lien présent dans la liste `links`, le champ `action` va définir comment l'élément sera importé dans le contexte graphe.

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

# Entités

## Description des champs d'entités

La liste `cis` est une liste d'entités représentée sous forme d'objet JSON. Chaque entité possède ces différents champs :
* Le champ **`_id`** contient l'identifiant de l'entité, sous la forme d'une chaine de caractères, qui est concerné par l'action.
* Le champ **`type`** contient le type de l'entité sous forme d'une chaine de caractères. Ce champ ne peut prendre qu'une valeur parmi les 4 suivantes : `resource`, `component`, `connector`, `watcher`.
* Le champ **`infos`** contient les informations complémentaires. Ce sont des données totalement personnalisables que l'utilisateur peut modifier.
* Le champ **`measurements`** contient un ensemble de chaine de caractères correspondant aux métriques liées à l'entité.
* Le champ **`name`** contient le nom de l'entité sous forme de chaine de caractères.
* Le champ **`action`** contient le type de l'action à réaliser au moment de l'import.
* Le champ **`action_properties`** vient en complément du champ action en spécifiant des informations complémentaires pour la bonne réalisation de l'action.

## Description des actions sur les entités

Il existe 6 actions supportés sur les entités au moment de l'import d'un contexte graphe : `create`, `set`, `delete`, `update`, `disable` et `enable`.

### Create
`create` crée une nouvelle entité dans le contexte graphe. Si une entité existe déjà avec le même identifiant, une **mise à jour complète** de l'entité sera effectuée. Les nouvelles données vont écraser l'entité courante en base de données.

Pour fonctionner correctement, les champs `_id`, `type` et `action` sont obligatoires.

Les champs dans la section [Description des champs d'entités](#Description-des-champs-d'entités) vont être copiés tels quels dans la nouvelle entité, à l'exception d'`action` et d'`action_properties`. Si le champ n'existe pas, le champ dans la nouvelle entité sera initialisé avec une valeur par défaut :
* `name` sera initialisé avec la valeur du champ `_id`
* `depends` et `impact` seront initialisés avec une liste vide
* `measurements` sera initialisé avec une liste vide
* et le champ `infos` par un objet vide.

### Set
`set`, comme `create`, crée une nouvelle entité dans le contexte graphe. La différence est que l'action *Set* procède à une **mise à jour partielle** de l'entité. Seuls les champs fournis dans les données d'import impacteront l'entité.

Par exemple, si une entité existante possède `info1` et `info2` dans ses `infos` et que l'import contient `info1` avec une valeur différente et `info3`, alors `info1` sera mis à jour avec sa nouvelle valeur, `info2` ne changera pas et `info3` sera créé.

Un exemple pratique dans un paragraphe suivant permettra de mieux saisie la différente entre `set` et `create`.

### Delete
`delete` supprime une entité désignée par son identifiant (champ `_id`). Si l'entité n'existe pas dans le contexte graphe, une erreur est déclenchée, l'import s'arrêtera et les modifications de l'import en cours ne seront pas répercutées sur le référentiel de Canopsis.

Pour fonctionner correctement, les champs `_id` et `action` sont obligatoires.

Lors de la suppression de l'entité, les champs `impact` et `depends` des entités ayant une relation avec celle supprimée verront les références à l'entité détruite supprimées.

### Update
`update` met à jour l'entité désignée par son `_id`. Si l'entité désignée par `_id` n'existe pas initialement dans le contexte graphe, une erreur sera déclenchée, l'import s'arrêtera et les modifications ne seront pas répercutées dans le référentiel.

Les champs `_id` et `action` sont obligatoires.

Pour la mise à jour, tous les champs présents dans l'action hormis `action` et `action_properties` seront recopiés dans l'entité à mettre à jour.

### Enable
L'action `enable` va activer une entité, c'est-à-dire que le champ `enabled` de l'entité sera mis à `True`. Si l'entité n'existe pas dans le contexte graphe, une erreur sera déclenchée, l'import s'arrêtera et les modifications ne seront pas répercutées dans le contexte graphe.

Les champs `_id`, `action` et `action_properties` sont obligatoires. Le champ `action_properties` doit avoir un champ `enable` contenant soit un entier
représentant un timestamp soit une liste d'horodatage correspondant à l'activation de l'entité.

### Disable
`disable` réalise l'action contraire du `enable`, à savoir désactiver une entité.

Les champs `_id`, `action` et `action_properties` sont obligatoires. Le champ `action_properties` doit avoir un champ `disable` correspondant à la désactivation de l'entité.

# Liens
## Description des champs des liens
La liste `links` représente toutes les relations entres les entités. Ces liens sont stockés sous forme d'objet JSON avec ces différents champs :
* Le champ **`_id`** contient l'identifiant de l'action, il peut prendre n'importe quelle valeur.
* Le champ **`from`** contient l'identifiant de l'entité de départ du lien
* Le champ **`to`** contient l'identifiant de l'entité d'arrivée du lien
* Le champ **`action`** contient le type de l'action à réaliser sous forme d'une chaine de caractères.

Pour tous les liens, les champs `from`, `to` et `action` sont obligatoires.

Les liens décrits dans les actions sont des liens de type *impact-depends*, c'est-à-dire qu'ils représentent des liens dont l'entité de départ (champ `from` dans le lien) contient dans son champ `impact` une référence à l'entité d'arrivée du lien. Par conséquent, l'entité d'arrivée du lien (champ `to`) contiendra une référence à l'entité de départ dans son champ `depends`.

## Description des actions sur les liens

Il y a deux actions supportées sur les liens : `create` et `delete`.

### Create
`create` crée un lien entre les deux entités définis à l'aide des champs `to` et `from` selon les modalités définies ci-dessus. Si au moins une des deux entités n'existe pas dans le contexte graphe, une erreur est déclenchée, l'import s'arrêtera et les modifications de l'import en cours ne seront pas répercutées sur le référentiel de Canopsis.

### Delete
`delete` supprime un lien entres les deux entités définis à l'aide des champs `to` et `from` selon les modalités définis ci-dessus. Si au moins une des deux
entités n'existe pas dans le contexte graphe, une erreur est déclenchée, l'import s'arrêtera et les modifications de l'import en cours ne seront pas répercutées sur le référentiel de Canopsis.

# Fonctionnement de l'import
Afin d'exécuter un nouvel import, il faut tout d'abord téléverser le fichier
au format JSON sur la route *api/contextgraph/import* en utilisant le verbe
HTTP PUT et en ajoutant l'import dans le corps de la requête. Le serveur
retournera un objet JSON contenant l'identifiant de l'import.

```json
{
  "total": 1,
  "data": [
    {
      "import_id": "b95e227f-27a2-4636-9f9c-ad30109d075d"
    }
  ],
  "success": true
}
```

Une fois l'import téléversé, l'import sera stocké sous la forme d'un fichier
afin que la *task* responsable d'un import s'occupe de l'import.

Lors du traitement de l'import par la tâche, si le fichier JSON est mal formé
ou qu'une anomalie est survenue lors du traitement d'une action, l'import
est annulé et les modifications ne sont pas répercutées sur le contexte graphe.
Si aucune erreur n'est rencontré, le contexte graphe sera mise à jour.

Dans tous les cas, la progression et le résultat de l'import sont disponible
via la route */api/contextgraph/import/status/<import_id>* accessible en GET.
La route retourne un objet JSON contenant au moins les champs :

  * **_id** de l'import.
  * **status** qui contient le status de l'import. Il y en a actuellement
  quatre : *"pending"*, *"ongoing"*, *"failed"*, *"done"*
  * **creation** qui contient la date et l'heure du téléversement de l'import
  sur le serveur.

## Import en attente
Lorsqu'un import est téléversé alors qu'un autre est en cours d'exécution,
le nouvel import est mis en attente. L'interrogation de la route précédemment
citée avec l'identifiant du nouvel import retournera un objet JSON
dont le *status* sera à **pending**.
```json
{
  "status": "pending",
  "_id": "6ac1deb9-1049-41f3-9e85-48d694deaab3",
  "creation": "Mon Aug 28 17:41:27 2017"
}
```

## Import en cours
Lorsqu'un import est en cours de traitement, l'interrogation de la route
précédemment citée avec l'identifiant du nouvel import retournera un objet JSON
dont le *status* sera à **ongoing**.
```json
{
  "status": "ongoing",
  "_id": "6ac1deb9-1049-41f3-9e85-48d694deaab3",
  "creation": "Mon Aug 28 17:41:27 2017"
}
```

## Import terminé
Lorsqu'un import est traité complétement sans erreur, la route précédemment
citée retourne un objet JSON avec le *status* à **done**, le nombre
d'entité supprimées dans *stats.deleted*, le nombres d'entité mises à jour ou créées
*stats.updated* et le temps d'exécution de l'import dans le champ *exec_time*.
```json
{
  "status": "done",
  "exec_time": "00:03:58",
  "_id": "6ac1deb9-1049-41f3-9e85-48d694deaab3",
  "creation": "Mon Aug 28 17:41:27 2017",
  "stats": {
    "deleted": 0,
    "updated": 20000
  }
}
```

## Échec de l'import
En cas d'erreur, la route précédemment citée retourne un objet JSON contenant
un status **failed**, une description de l'erreur dans le champ **infos** et le
temps d'exécution total de l'import dans **exec_time**.
```json
{
  "status": "failed",
  "info": "ValueError(u'The ci of id connector_0 match an existing entity.',)",
  "_id": "6ac1deb9-1049-41f3-9e85-48d694deaab3",
  "creation": "Mon Aug 28 17:41:27 2017",
  "exec_time": "00:02:55"
}
```

### Relancer manuellement un import

Lorsqu'un job est bloqué en pending, il est possible de relancer manuellement la tâche d'importation en publiant un event forgé dans la queue `task_importctx` :
```json
{
  "jobid": "importctx_2c0c0b49-129a-41aa-bccd-cc23de478bbc",
  "jobs_uuid": "2c0c0b49-129a-41aa-bccd-cc23de478bbc"
}
```
(remplacer par l'id de votre job ; voir par exemple dans `/opt/canopsis/tmp` le nom du fichier json en attente d'importation ou dans la collection `default_importctx`)


# Exemple d'utilisation de l'import

 * Tout d'abord, il faut s'authentifier à auprès de l'API avec votre
*authentication key*

```
GET http://192.168.0.93:8082/?authkey=6b6ce450-5fd2-11e7-b5dd-0800279471b5
```

 * Ensuite, il faut téléverser l'import sur la route *api/contextgraph/import/*
```JSON
PUT http://127.0.0.1:8082/api/contextgraph/import/
{
    "cis": [
        {
            "type": "component",
            "infos": {
            },
            "_id": "component_0",
            "action": "create",
            "measurements": [],
            "action_properties": {
            },
            "name": "component_0",
            "impact": [],
            "depends": []
        },
        {
            "type": "resource",
            "infos": {
            },
            "_id": "resource_1",
            "action": "create",
            "measurements": [],
            "action_properties": {
                "disable": "494172"
            },
            "name": "resource_1",
            "impact": [],
            "depends": []
        }
    ],
    "links": [
        {
            "to": "component_0",
            "from": "resource_1",
            "infos": {
            },
            "_id": "component_0-to-resource_1",
            "action": "create",
            "action_properties": {
            }
        }
    ]
}
```
```json
{
  "total": 1,
  "data": [
    {
      "import_id": "01834fe6-181a-4312-8900-9ea24901bda0"
    }
  ],
  "success": true
}
```

 * Pour obtenir l'état de l'import :
```
GET http://192.168.0.93:8082/api/contextgraph/import/status/01834fe6-181a-4312-8900-9ea24901bda0
```
```json
{
  "status": "done",
  "exec_time": "00:00:01",
  "_id": "01834fe6-181a-4312-8900-9ea24901bda0",
  "creation": "Mon Aug 28 17:41:27 2017"
  "stats": {
    "deleted": 0,
    "updated": 20000
  }
}
```
