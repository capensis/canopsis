# Moteur `engine-watcher` (Go, Core)

!!! note
    Cette page concerne les observateurs nouvelle génération, disponibles uniquement avec le moteur Go `engine-watcher`.

Les moteur `engine-watcher` permet de surveiller et de répercuter les états d'alarmes ouvertes sur des entités surveillées.

Les observateurs sont définis dans la collection MongoDB `default_entities`, et peuvent être ajoutés et modifiés avec l'[API watcherng](../../guide-developpement/api/api-v2-watcherng.md).

Des exemples pratiques d'utilisation des observateurs sont disponibles dans la partie [Exemples](#exemples).

## Utilisation

La file du moteur est placée juste après le moteur [`engine-axe`](moteur-axe.md).

### Options de l'engine-watcher

```
  -autoRecomputeWatchers
      Automatically recompute watchers each minute.
  -d	debug
  -printEventOnError
      Print event on processing error
  -publishQueue string
      Publish event to this queue. (default "Engine_action")
  -version
      version infos
```

## Fonctionnement

### Concept d'un observateur

Un observateur représente un groupe de surveillance.
C'est-à-dire que la criticité d'une entité de type watcher dépendra de la criticité des entités surveillées, et des alarmes ouvertes sur ces entités.

Le but d'un observateur est de donner une visibilité accrue et claire sur l'état d'un groupe d'entités, afin de détecter un changement de criticité positif ou négatif sur les alarmes liées aux entités du groupe surveillé.

### Définition d'un observateur

Un observateur est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel): l'identifiant de l'observateur (généré automatiquement ou choisi par l'utilisateur).
 - `name` (requis) : Le nom de l'observateur, qui sera utilisé dans la météo de services.
 - `entities` (requis) : La liste des patterns permettant de filtrer les entités surveillées. Le format des patterns est le même que pour l'[event-filter](moteur-che-event_filter.md).
 - `state` (requis) : Un document contenant :
    - `method` (requis) : Le nom de la méthode de calcul de la criticité de l'observateur, en fonction des alarmes ouvertes sur les entités. Actuellement, seule la méthode `worst` est implémentée.
    - Les différents paramètres des méthodes ci-dessus.
- `output_template` (requis) : Le template utilisé par l'observateur pour déterminer la sortie de l'alarme.

Le schéma en base est proche, puisqu'il s'agit de ces paramètres, ajoutés à ceux déjà présents pour une entitié.

#### Méthodes

Actuellement, seule la méthode `worst` est implémentée.

- `worst` : la criticité de l'observateur est la criticité de la pire alarme ouverte sur les entités surveillées.

#### Templates

L'`output_template` est un [template](https://golang.org/pkg/text/template/)
permettant d'afficher diverses informations dans l'output de l'alarme
correspondant à l'observateur.

Les informations disponibles sont :

 - `{{.Alarms}}` : le nombre d'alarmes en cours sur les entités observées par l'observateur.
 - `{{.State.Info}}` : le nombre d'entités observées n'ayant pas d'alarmes, ou une alarme en criticité `Info`.
 - `{{.State.Minor}}` : le nombre d'alarmes mineures sur les entités observées.
 - `{{.State.Major}}` : le nombre d'alarmes majeures sur les entités observées.
 - `{{.State.Critical}}` : le nombre d'alarmes critiques sur les entités observées.
 - `{{.Acknowledged}}` : le nombre d'alarmes acquittées sur les entités observées.
 - `{{.NotAcknowledged}}` : le nombre d'alarmes non-acquittées sur les entités observées.

Par exemple, l'output d'un observateur avec l'`output_template` suivant :

```
Crit : {{.State.Critical}} / Total : {{.Alarms}}
```

sera

```
Crit : 12 / Total : 60
```

s'il y a 60 alarmes en cours dont 12 critiques sur les entités observées.

#### Exemples

```json
{
    "_id": "h4z25rzg6rt-64rge354-5re4g",
    "name": "Client Capensis",
    "entities": [{
        "infos": {
            "customer": {
                "value": "capensis"
            }
        }
    }, {
        "_id": {"regex_match": ".+/comp"}
    }],
    "state": {
        "method": "worst",
    },
    "output_template": "Alarmes critiques : {{.State.Critical}}, Alarmes acquittées : {{.Acknowledged}}"
}
```
