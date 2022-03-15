# Moteur `engine-correlation` (Pro)

Le moteur `engine-correlation` permet de créer des méta alarmes à partir de `règles de gestion`.

Des exemples pratiques d'utilisation de la corrélation sont disponibles dans la partie [types de groupements](#types-de-groupements).

## Utilisation

### Options du moteur

La commande `engine-correlation -help` liste toutes les options acceptées par le moteur.

### Multi-instanciation

Il est possible, à partir de Canopsis 3.43.0, de lancer plusieurs instances du moteur `engine-correlation`, afin d'améliorer sa performance de traitement et sa résilience.

En environnement Docker, il vous suffit par exemple de lancer Docker Compose avec `docker-compose up -d --scale correlation=2` pour que le moteur `engine-correlation` soit lancé avec 2 instances.

Cette fonctionnalité sera aussi disponible en installation par paquets lors d'une prochaine mise à jour.

## Fonctionnement

Tous les événements qui circulent dans Canopsis sont transmis à la file `Engine_correlation` et seront donc lus par le moteur `engine-correlation`.

Le moteur va alors vérifier si une règle de corrélation doit s'appliquer et si c'est le cas, il générera ou modifiera une méta alarme.

Une méta alarme permet de grouper dynamiquement des alarmes selon des critères définis dans une règle.

Ces groupements d'alarmes permettent de voir rapidement les corrélations entre les différentes alarmes présentes dans le [bac à alarmes](../../guide-utilisation/interface/widgets/bac-a-alarmes/index.md).

### Activation de la corrélation

Dans les bacs à alarmes, un interrupteur situé en haut du widget vous permet d'activer la corrélation.

![Activer la corrélation](./img/correlation_activer.png)

### Définition d'un groupement

Un groupement d'alarmes se caractérise par les informations suivantes.

| Champ | Type | Description |
|-------|------|-------------|
| `_id`   | string  | Identifiant unique du groupement, généré par MongoDB lui-même.  |
| `name`  | string  | Nom donné au groupement, il apparaîtra sur les méta alarmes dans le bac.  |
| `type`  | string  | Type de groupement : `relation`, `timebased`, `attribute`, `complex`, ou `valuegroup`.  |
| `time_interval`  | int  | Intervalle de temps en secondes.  |
| `threshold_count`  | int  | Le *seuil de déclenchement* exprime un nombre d'alarmes au delà duquel le groupement sera effectué.  |
| `threshold_rate`  | float  | Le *taux de déclenchement* exprime le pourcentage d'entités impactées au delà duquel le groupement aura lieu. |
| `value_paths`  | array  | Les *chemins de valeurs* désignent les adresses d'attributs à partir desquels le groupement va opérer. |
| `auto_resolve`  | bool  | Permet de résoudre la méta alarme lorsque toutes ses alarmes conséquences sont résolues. |
| `output_template`  | string  | L'`output_template` permet de formater le message informatif associé à une méta alarme. |

Dans un groupement de type `complex`, `threshold_count` et `threshold_rate` sont mutuellement exclusifs.

### Types de groupements

Il existe plusieurs types de groupements pour corréler vos alarmes.

#### Groupement par relation parent-enfant

Il permet de regrouper les alarmes qui ont un lien de parenté.  
Par exemple, si un composant a provoqué une alarme, toutes les alarmes des ressources ayant le même composant seront regroupées dans une même méta alarme.

Exemple :

```json
{
  "name": "Règle de groupement par relation",
  "type": "relation"
}
```

Cette règle s'applique à toutes les entités.

#### Groupement par intervalle de temps

Ce type de méta alarme regroupe toutes les alarmes survenues dans un intervalle de temps donné.  
Par exemple toutes les alarmes apparues au cours de la même minute (à partir de la création de la première alarme) seront regroupées sous une même méta alarme.

Exemple :
```json
{
  "name": "Règle de groupement par intervalle de temps",
  "type": "timebased",
  "output_template" : "Regroupé par intervalle de temps",
  "config": {
    "time_interval": 60
  }
}
```

#### Groupement par attribut

Ce type de groupement utilise les mêmes [patterns](moteur-che-event_filter.md#patterns) que les autres moteurs pour identifier un attribut dans l'évènement, dans l'entité ou dans l'alarme.  
Par exemple si on utilise un `event_pattern` qui vaut `component = srv001`, toutes les alarmes créées à partir d'un évènement dont le composant est égal à srv001 seront regroupées dans une méta alarme.

Exemple :

```json
{
  "name": "Règle de groupement par attribut",
  "type": "attribute",
  "output_template" : "Regroupé par attribut",
  "config": {
    "alarm_patterns": [
      {
        "v": {
          "component": "srv001"
        }
      }
    ]
  }
}
```

#### Groupement complexe avec seuil de déclenchement

C'est une combinaison de groupement par attribut et de groupement par intervalle de temps, il possède aussi une notion de seuil de déclenchement.  
Par exemple, on pourra l'utiliser pour regrouper toutes les alarmes créées pour une même entité durant un intervalle de temps donné, seulement si le nombre d'alarmes créées dépasse un certain seuil.

Exemple :

```json
{
  "name": "Règle de groupement complexe",
  "type": "complex",
  "output_template" : "Regroupement complexe",
  "config": {
    "time_interval": 60,
    "threshold_count": 3,
    "alarm_patterns": [
      {
        "v": {
          "resource": "check"
        }
      }
    ]
  }
}
```

Cette règle s'applique si 3 alarmes ou plus, dont la ressource vaut `check`, ont été créées durant un intervalle de temps de 60 secondes.

#### Groupement complexe avec taux de déclenchement

Comme le précédent il s'agit d'un groupement par attribut et par intervalle de temps. Mais celui-ci calcule le taux d'entités en erreur par rapport à un groupe donné, défini par un [pattern](moteur-che-event_filter.md#patterns).  
Par exemple, vous avez 5 serveurs situés sur le site de Wasquehal et vous souhaitez grouper les alarmes si au moins 80% de ces serveurs sont en erreur sur une période de 5 minutes. Au préalable, vous aurez pris soin de créer ces entités dans le référentiel interne de Canopsis et de les enrichir avec `infos.site.value = Wasquehal`.

Exemple :

```json
{
  "name": "Groupement dysfonctionnement global site Wasquehal",
  "output_template" : "Regroupé à partir d'un seuil de déclenchement",
  "type": "complex",
  "config": {
    "time_interval": 300,
    "threshold_rate": 0.80,
    "entity_patterns": [
      {
        "infos": {
            "site": {
                "value": "Wasquehal"
            }
        }
      }
    ]
  }
}
```

**Cas particulier nécessitant une autre base de calcul du ratio**

Lorsque Canopsis reçoit des événements et qu'il enrichit les entités correspondantes grâce à des règles, l'application d'un ratio peut paraître inadaptée.  
Par exemple, si les événements reçus concernent des `ping`, pour le calcul du ratio, nous avons besoin :

* Du nombre d'alarmes `ping`
* Du nombre total d'entités concernées par ces ping

Dans le cas où vous n'avez pas de référentiel pré rempli dans Canopsis, le nombre total d'entités concernées par ces ping n'est autre que le nombre d'alarmes `ping` déjà reçues par Canopsis jusque là.

Ce ratio ainsi calculé n'a pas de sens métier dans cet exemple.

Le type `complex` propose donc une solution pour utiliser une base de ratio différente. Cette base est représentée par l'attribut `total_entity_patterns`.

La règle ci-après permet de créer une méta alarme lorsque le ratio entre le nombre d'alarmes `ping` et le nombre de composants est supérieur à 80%

```json 
{
	"_id" : "complextest-total-entity-pattern",
	"auto_resolve" : false,
	"name" : "complextest-total-entity-pattern",
	"config" : {
		"time_interval" : 3600,
		"alarm_patterns" : [
			{
				"v" : {
					"resource" : "ping"
				}
			}
		],
		"threshold_rate" : 0.8,
		"total_entity_patterns" : [
			{
				"type" : "component",
			}
		]
	},
	"patterns" : null,
	"output_template" : "{{ .Children.Alarm.Value.State.Message }}",
	"type" : "complex"
}
```

#### Groupement groupe de valeurs

Ce type de règle possède les mêmes attributs que le type `complex` avec la notion de *chemins de valeur (value_paths)* en plus.  
Ces *chemins de valeur* sont utilisés pour grouper les valeurs de manière unique.

Exemple :

```json
{
    "name": "Groupe de valeurs",
    "type": "valuegroup",
    "output_template" : "Regroupé par un chemin de valeurs",
    "config": {
        "time_interval": 3600,
        "threshold_count": 5,
        "value_paths": [
            "entity.infos.site.value"
        ],
        "entity_patterns": [
            {
                "infos": {
                    "site": {
                        "value": {
                            "regex_match": ".+"
                        }
                    }
                }
            }
        ]
    }
}
```

Cette règle s'applique si 5 alarmes ou plus concernant le même chemin de valeur (`entity.infos.site.value`), ont été créées durant un intervalle de temps de 3600 secondes.
Plusieurs chemins de valeur peuvent être observés. Dans ce cas, une méta alarme est créée pour chaque combinaison de chemins de valeurs.

### Démonstration vidéo

La démonstration en vidéo suivante permet de constater l'effet obtenu avec la
règle de corrélation `relation` (relation parent-enfant) :

<iframe width="560" height="315" src="https://www.youtube-nocookie.com/embed/Y-OwUIqiu_s" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

D'autres vidéos sont disponibles pour les autres types de règles de corrélation.
Voir la [playlist « Tutoriels Canopsis » (YouTube)][youtube-pl].

[youtube-pl]: https://www.youtube.com/playlist?list=PLOoYSSN74hIZVLA1oGkZpora2scQvuFNJ

### Processus de création d'une méta alarme

En reprenant l'exemple d'un groupement par relation parent-enfant, la création d'une méta alarme se déroule de la façon suivante.

Une erreur se produit sur un composant `component_meta` qui envoie un évènement à Canopsis.

```json
{
  "event_type": "check",
  "connector": "connector_meta",
  "connector_name": "test_meta",
  "component": "component_meta",
  "source_type": "component",
  "author": "Canopsis",
  "state": 3,
  "output": "Exemple de création méta alarmes - component"
}
```

Une alarme sur le composant est créée et visible dans le bac à alarmes.

![Alarme en cours sur le composant](img/correlation_alarme_isolee.png)

Comme le composant est en erreur, ses ressources envoient probablement elles aussi des évènements qui vont conduire à la création d'alarmes.

```json
{
  "event_type": "check",
  "connector": "connector_meta",
  "connector_name": "test_meta",
  "component": "component_meta",
  "resource": "resource_test_X",
  "source_type": "resource",
  "author": "Canopsis",
  "state": 3,
  "output": "Exemple de création méta alarmes - ressource X"
}
```

Le moteur corrélation détecte la relation et l'alarme sur le composant devient alors une méta alarme. Les alarmes sur les ressources sont vues comme des conséquences du dysfonctionnement du composant et sont regroupées dans la méta alarme.

Dans l'interface cela se traduit de la façon suivante :

![Meta-alarme relation parent-enfant](img/correlation_alarmes_groupees.png)

Au passage de la souris sur l'icône de la colonne « Extra Details » une info-bulle apparaît et indique le nom de la règle qui a créé la méta alarme et le nombre d'alarmes liées.

Le détail des alarmes liées peut-être consulté en déroulant le détail de la méta alarme et en affichant l'onglet « Alarmes Liées ».

Dans le cas ci-dessus (groupement par relation parent-enfant), c'est l'alarme parent qui devient la méta alarme. Ses valeurs de `connector`, `connector_name` et `component` sont alors conservées. Pour les autres types de groupement, une méta alarme est générée par Canopsis. Celle-ci aura toujours les valeurs suivantes :

- connector : `engine`
- connector_name : `correlation`
- component : `metaalarm`

### Variables d'une méta alarme

La méta alarme a été enrichie avec de nouvelles variables.

![Variables de la méta alarme](img/correlation_variables_meta_alarme_01.png)

- `consequences` : contient le nombre d'alarmes liées ainsi que les données de ces alarmes dans un tableau `data`.

- `metaalarm` : la valeur 1 indique qu'il s'agit d'une méta alarme.

- `rule` : contient l'`id` et le `name` de la règle qui a déclenché la création de cette méta alarme.

![Variables de la méta alarme](img/correlation_variables_meta_alarme_02.png)

- `v.children` : embarque les `id` des entités liées.

- `v.parents` : est présent dans les variables mais ne sera renseigné que pour les alarmes liées à la méta alarme.

### Format du message informatif d'une méta alarme

Le message informatif, ou `output`, peut être formaté grâce à un template porté par la variable `output_template`.

Cet attribut peut-être composé de **texte statique** et/ou de **variables**.

Les variables accessibles sont : 

* **.Count** : Le nombre d'alarmes conséquences attachées à la méta alarme.
* **.Children** : L'ensemble des variables de la dernière alarme conséquence attachée à la méta alarme.
* **.Rule** : Les informations administratives de la méta alarme en elle-même

Quelques exemples :

* {{ .Count }} Alarmes conséquences; {{ .Children.Alarm.Value.State.Message }}; Règle : {{ .Rule.Name }}
* Un message informatif statique
* Corrélé par la règle {{ .Rule.Name }}

!!! Note
    Il est possible à partir de la version 3.47.0 de Canopsis d'utiliser les informations des méta alarmes dans des règles d'enrichissement. La documentation associée se trouve [ici](./moteur-che-event_filter.md#cas-particulier-des-meta-alarmes)

## Collection MongoDB associée

Les entités sont stockées dans la collection MongoDB `meta_alarm_rules`.

```json
{
	"_id" : "73da1ad7-058e-46af-8442-7ea3f246eb68",
	"patterns" : null,
	"config" : null,
	"name" : "Relation-composant-ressource",
	"type" : "relation"
}
```
