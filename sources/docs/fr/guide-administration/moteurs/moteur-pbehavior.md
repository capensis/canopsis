# Pbehavior

Les *pbehaviors* (pour *periodical behaviors*) sont des évènements de calendrier récurrents qui permettent de mettre en pause la surveillance d'une alarme pendant une période donnée (pour des maintenances ou des astreintes, par exemple).

Ils permettent de créer des « downtimes », à savoir indiquer qu'une entité est en pause.

Les pbehaviors sont définis dans la collection MongoDB `default_pbehavior`, et peuvent être ajoutés et modifiés avec l'[API PBehavior](../../guide-developpement/api/api-v2-pbehavior.md).

## Fonctionnement

Dans une stack en Go classique, la file du moteur `pbehavior` n'est pas alimentée.

Un PBehavior contient un filtre (`filter`) qui est appliqué sur une entité.

Chaque minute, le moteur calcule les PBehaviors et leur application sur les entités.

## Définition d'un PBehavior

Un pbehavior se caractérise par les informations suivantes.

|   Champ    |  Type  |                                             Description                                              |     |
| ---------- | ------ | ---------------------------------------------------------------------------------------------------- | --- |
|   `_id`    | string |                   Identifiant unique du comportement, généré par MongoDB lui-même.                   |     |
|   `eids`   | liste  |                  Liste d'identifiants d'entité qui correspond au filtre précédent.                   |     |
|   `name`   | string |                     Type de pbehavior. `downtime` est la seule valeur acceptée.                      |     |
|  `author`  | string |                            Auteur ou application ayant créé le pbehavior.                            |     |
| `enabled`  |  bool  |    Activer ou désactiver le pbehavior, pour qu’il puisse être ignoré, même sur une plage active.     |     |
| `comments` | liste  |                                 `null` ou une liste de commentaires.                                 |     |
|  `rrule`   | string |                                 Champ texte défini par la RFC 2445.                                  |     |
|  `tstart`  |  int   | Timestamp fournissant la date de départ du pbehavior, recalculée à partir de la `rrule` si présente. |     |
|  `tstop`   |  int   |  Timestamp fournissant la date de fin du pbehavior, recalculée à partir de la `rrule` si présente.   |     |
|  `type_`   | string |                         Optionnel. Type de pbehavior (pause, maintenance…).                          |     |
|  `reason`  | string |                       Optionnel. Raison pour laquelle ce pbehavior a été posé.                       |     |
| `timezone` | string |                       Fuseau horaire dans lequel le pbehavior doit s'exécuter.                       |     |
|  `exdate`  | array  |                     La liste des occurrences à ignorer sous forme de timestamps                      |     |


Un exemple d'évènement pbehavior brut :
```js
{
   "_id" : string,
   "name" : string,
   "filter" : string,
   "comments" : [ {
       "_id": string,
       "author": string,
       "ts": timestamp,
       "message": string
   } ],
   "tstart" : timestamp,
   "tstop" : timestamp,
   "rrule" : string,
   "enabled" : boolean,
   "eids" : [ ],
   "connector" : string,
   "connector_name" : string,
   "author" : string,
   "timezone" : string,
   "exdate" : [
      timestamp
   ]
}
```

### Filtrage d'entités (`filter`)

Le champ `filter` permet de filtrer les entités sur lesquelles le PBehavior est appliqué.

Il peut prendre en charge les conditions `or` et `and` mais nécessite de les échapper.

Exemple :

```json
{
	"author": "root",
	"name": "Pbehavior test 2",
	"tstart": 1567439123,
	"tstop": 1569599100,
	"filter": {
		"$or": [{
			"impact": {
				"$in": ["pbehavior_test_1"]
			}
		}, {
			"$and": [{
				"type": "component"
			}, {
				"name": "pbehavior_test_1"
			}]
		}]
	},
	"type_": "Hors plage horaire de surveillance",
	"reason": "Problème d'habilitation",
	"rrule": null,
	"comments": [],
	"exdate": []
}
```

C'est un filtre appliqué directement sur les champs des entités contenues dans la collection `default_entities` de MongoDB.

### Règles de récurrence (`rrule`)

C'est une règle de récurrence du PBehavior.

Dans le cas où la `rrule` est absente, `tstart` et `tstop` font office de plage d’activation du pbehavior, sans récurrence.

Dans le cas où la `rrule` est présente, `tstart` et `tstop` seront recalculés afin de refléter la récurrence.

### Dates d'exclusion (`exdate`)

Il est possible d'empêcher l'exécution d'une occurrence d'un pbehavior, à l'aide du champ `exdate`.

`exdate` est une liste de timestamps correspondant au début d'une occurence à empêcher.

### Fuseau horaire (`timezone`)

L'exécution de chaque pbehavior se fait dans un fuseau horaire particulier.

Lorsqu'un pbehavior ne contient pas de champ `timezone`, le fuseau utilisé sera celui défini dans le fichier `/opt/canopsis/etc/pbehavior/manager.conf` sous le champ `default_timezone`.

Si le fichier de configuration n'existe pas ou si le champ `default_timezone` n'existe pas, le fuseau `Europe/Paris` sera utilisé.

Si le fuseau horaire choisi comporte des heures d'hiver et d'été, celles-ci seront respectées tout au long de l'année. Ainsi, un pbehavior devant se déclencher à 16 heures s'exécutera à 16 heures en heure d'été et à 16 heures en heure d'hiver.

## Représentation dans MongoDB

Les pbehaviors sont stockés dans la collection MongoDB `default_pbehavior` (voir [API PBehavior](../../guide-developpement/api/api-v2-pbehavior.md) pour la création des pbehaviors).

Un exemple de pbehavior appliqué pour une plage de maintenance sans `rrule` avec la raison `Problème d'habilitation` et le type `Maintenance` aux alarmes dont le composant est `pbehavior_test_1`.

```json
{
    "_id" : "145331d4-d536-4c58-8e6d-229d5d8f3f10",
    "filter" : "{\"$or\": [{\"impact\": {\"$in\": [\"pbehavior_test_1\"]}}, {\"$and\": [{\"type\": \"component\"}, {\"name\": \"pbehavior_test_1\"}]}]}",
    "name" : "Pbehavior test 2",
    "author" : "root",
    "enabled" : true,
    "type_" : "Hors plage horaire de surveillance",
    "comments" : [],
    "connector" : "canopsis",
    "reason" : "Problème d'habilitation",
    "connector_name" : "canopsis",
    "eids" : [
        "pbehavior_test_1",
        "disk2/pbehavior_test_1"
    ],
    "tstart" : 1567439123,
    "tstop" : 1569599100,
    "timezone" : "Europe/Paris",
    "exdate" : [],
    "rrule" : null
}
```
