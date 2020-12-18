# Moteur `pbehavior` (Go, Core)

Les comportements périodiques (*pbehaviors*, pour *periodical behaviors*) sont des évènements de calendrier, récurrents ou non, qui permettent de modifier la surveillance d'une alarme pendant une période donnée (pour des maintenances ou des astreintes, par exemple).

Ils permettent de créer des « downtimes » et donc d'indiquer qu'une entité est en pause.

Les comportements sont définis dans la collection MongoDB `pbehavior`.

!!! Note
    Avec la v4 de Canopsis le fonctionnement des comportements périodiques à été complètement revu.
    Les informations qui figurent sur cette page ne sont donc valables que pour cette version.

## Utilisation

### Options du moteur

La commande `engine-pbehavior -help` liste toutes les options acceptées par le moteur.

## Fonctionnement

Ce moteur doit toujours être présent. Ce moteur est écrit en langage Go.

Un comportement périodique est caractérisé par un type et une raison (voir ci-dessous). Il contient également un filtre (`filter`) qui est appliqué sur le référentiel des entités.

Les comportements périodiques existants sont appliqués immédiatement sur les nouvelles alarmes. De la même façon, les comportements périodiques nouvellement créés sont appliqués immédiatement sur les alarmes existantes.

Ensuite, chaque minute, le moteur calcule les comportements périodiques et leur application sur les entités.

Un seul comportement peut être actif, à un moment donné, sur une entité.

## Administration de la planification

### Configuration des types de comportements périodiques

Rendez vous dans le menu Administration puis dans Administration de la planification.

![Menu administration de la planification](./img/menu-administration-planification.png)

Les types par défaut s'affichent à l'écran : `actif`, `inactif`, `maintenance` et `pause`. Ils ne peuvent être ni supprimés, ni modifiés. La priorité des types est gérée dans l'ordre croissant, c'est à dire, 0 est la priorité la plus faible et 3 est la plus forte et sera traitée avant les autres. Un seul type de comportement périodique peut être actif sur une entité à un moment donné.

![Types de comportements périodiques](./img/admin-planification-types-defaut.png)

### Création d'un type personnalisé

Cliquez sur le bouton `+` en bas à droite de la fenêtre pour ouvrir la fenêtre de création.

![Créer un type personnalisé](./img/admin-planification-creer-type.png)

Renseignez les différents champs, choisissez un type parmi la liste et affectez lui une priorité et une icône.

!!! Attention
    Deux types ne peuvent avoir la même priorité.

![Formulaire type personnalisé](./img/admin-planification-type-personnalise.png)

Cliquez sur le bouton Soumettre et votre type personnalisé apparaît dans la liste.

![Liste des types personnalisés](./img/admin-planification-liste-type-perso.png)

### Configuration des raisons

Cliquez sur l'onglet Raison. Par défaut, la liste des raisons est vide. Comme pour les types vous pouvez cliquer sur le bouton `+` pour créer une nouvelle raison. Chaque raison doit avoir un nom et une description.

Voici, par exemple, une liste de raisons personnalisées :

![Liste de raisons personnalisées](./img/admin-planification-liste-raisons.png)

## Configuration des dates d'exception

Il est également possible de configurer des dates d'exceptions dans l'onglet dédié. Pour cela cliquez de nouveau sur le bouton `+` pour créer une liste d'exceptions.

Vous pourrez alors renseigner un nom, une description et ajouter des dates dans la liste. A chaque date vous pourrez associer un des types existants.

![Créer une liste d'exceptions](./img/admin-planification-liste-exceptions.png)

## Définition d'un comportement périodique

Un comportement périodique se caractérise par les informations suivantes.

|   Champ    |  Type  |                                             Description                                              |     |
| ---------- | ------ | ---------------------------------------------------------------------------------------------------- | --- |
|   `_id`    | string |                   Identifiant unique du comportement, généré par MongoDB lui-même.                   |     |
|   `name`   | string |       Nom usuel donné au comportement périodique.                                                    |     |
|  `author`  | string |              Auteur ou application ayant créé le comportement périodique.                            |     |
| `enabled`  |  bool  |    Activer ou désactiver le pbehavior, pour qu’il puisse être ignoré, même sur une plage active.     |     |
| `comments` | liste  |                                 `null` ou une liste de commentaires.                                 |     |
|  `rrule`   | string | Règle de récurrence, champ texte [défini par la RFC 2445](https://www.kanzaki.com/docs/ical/recur.html).  |   |
|  `tstart`  |  int   | Timestamp fournissant la date de départ, recalculée à partir de la `rrule` si présente.              |     |
|  `tstop`   |  int   |  Timestamp fournissant la date de fin, recalculée à partir de la `rrule` si présente.                |     |
|  `type_`   | string |             Type de comportement périodique (pause, maintenance…).                                   |     |
|  `reason`  | string | Id de la raison pour laquelle ce comportement périodique a été posé. Les raisons de trouvent dans la collection `pbehavior_reason`.         |     |
|  `exdates` | liste  |  Liste des dates d'exception de la rrule.                                                            |     |
| `exceptions` | liste | Liste des id des listes d'exceptions attachées à ce comportement. Les listes d'exceptions se trouvent dans la collection `pbehavior_exception` |     |

Un exemple d'évènement `pbehavior` brut :
```js
{
  "_id" : "string",
  "name" : "string",
  "filter" : "string",
  "comments" : [ {
    "_id": "string",
    "author": "string",
    "ts": timestamp,
    "message": "string"
  } ],
  "tstart" : timestamp,
  "tstop" : timestamp,
  "rrule" : "string",
  "enabled" : boolean,
  "author" : "string",
  "reason" : "string",
  "exceptions" : [ {
    timestamp
  } ],
  "exdates" : [ {
    "begin" : timestamp,
    "end" : timestamp,
    "type" : "string"
  } ]
}   
```

#### Définition d'une raison

Un comportement périodique se caractérise par les informations suivantes.

|   Champ    |  Type  |                                             Description                                              |     |
| ---------- | ------ | ---------------------------------------------------------------------------------------------------- | --- |
|   `_id`    | string |                   Identifiant unique de la raison, généré par MongoDB lui-même.                      |     |
|   `name`   | string |               Nom usuel renseigné lors de la création de la raison.                                  |     |
| `description` | string |            Description de la raison.                                                              |     |

#### Définition d'une liste d'exceptions

Une liste d'exceptions se caractérise par les informations suivantes.

|   Champ    |  Type  |                                             Description                                              |     |
| ---------- | ------ | ---------------------------------------------------------------------------------------------------- | --- |
|   `_id`    | string |                   Identifiant unique de la liste, généré par MongoDB lui-même.                      |     |
|   `name`   | string |               Nom usuel renseigné lors de la création de la liste.                                  |     |
| `description` | string |            Description de la liste.                                                              |     |
| `exdates`  | int   |               Liste des dates d'exceptions et de leur type                                           |     |

### Filtrage d'entités (`filter`)

Le champ `filter` permet de filtrer les entités sur lesquelles le comportement périodique est appliqué.

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

C'est une règle de récurrence du comportement périodique.

Dans le cas où la `rrule` est absente, `tstart` et `tstop` font office de plage d’activation, sans récurrence.

Dans le cas où la `rrule` est présente, `tstart` et `tstop` seront recalculés afin de refléter la récurrence.
