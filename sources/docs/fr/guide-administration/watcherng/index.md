# Watcher

Les watchers nouvelle génération sont une fonctionnalité du moteur `axe` permettant de surveiller et de répercuter les états d'alarmes ouvertes sur des entités surveillées.

Les watchers sont définis dans la collection MongoDB `default_entities`, et
peuvent être ajoutés et modifiés avec l'[API watcherng](../../guide-developpement/watcherng/api_v2_watcherng.md).

Des exemples pratiques d'utilisation des watchers sont disponibles dans la partie [Exemples](#exemples).

## Concept d'un watcher

Un watcher, ancienne comme nouvelle génération, représente un groupe de surveillance.  
C'est à dire que l'état d'une entité de type watcher dépendra de l'état des entités surveillées, et des alarmes ouvertes sur ces entités.  

Le but de d'un watcher est de donner une visibilité accrue et claire sur l'état d'un groupe d'entités, afin de détecter un changement d'état positif ou négatif sur les alarmes liées aux entités du groupe surveillé.

Son fonctionnement interne est le suivant :  
**Initialisation**  
- A sa création, le watcher récupère les entités surveillées, et met à jour les champs `impact` des entités avec son `_id`, ainsi que son champ `depends` avec les `_id` des entités surveillées.  
- Lorsqu'une entité est créée, les watchers se mettent automatiquement à jour comme à leur création, ainsi que l'entité nouvellement créée.

**Calcul de l'état**
A chaque top de l'engine axe, chaque watcher recalcule son état. Pour calculer son état, chaque watcher :  
- Récupère ses entités surveillées.
- Récupère les alarmes ouvertes sur ces entités.
- Garde un compte de chaque état d'alarme, respectivement `Info`, `Minor`, `Major` et `Critical`.
- Applique l'algorithme présent dans le champ `state.method` avec les éventuels paramètres afin de calculer l'état.
- Calcule sa sortie en fonction de l'`output_template`.
- Envoie un évènement de type `check` à `axe` contenant l'état calculé `state` et la sortie `output` afin que le moteur puisse prendre en compte ces changements.

## Définition d'un watcher

Un watcher est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel): l'identifiant du watcher (généré automatiquement ou choisi par l'utilisateur).
 - `name` (requis) : Le nom du watcher, qui sera utilisé dans la météo de services.
 - `entities` (requis) : La liste des patterns permettant de filtrer les entités surveillées. Le format des patterns est le même que pour l'[event-filter](../event-filter/index.md).
 - `state` (requis) : Un document contenant :
    - `method` (requis) : Le nom de la méthode de calcul de l'état du watcher en fonction des alarmes ouvertes sur les entités. Actuellement, seule la méthode `worst` est implémentée.
    - Les différents paramètres des méthodes ci-dessus.
- `output_template` (requis) : Le template utilisé par le watcher pour déterminer la sortie de l'alarme.

Le schéma en base est proche, puisqu'il s'agit de ces paramètres, ajoutés à ceux déjà présents pour une entitié.

### Méthodes

Actuellement, seule la méthode `worst` est implémentée.
- `worst` : L'état du watcher est l'état de la pire alarme ouverte sur les entités surveillées.

### Templates

L'`output_template` est une chaîne de caractère contenant, entre autres, un accès à un compte des alarmes selon leur état. Les comptes sont accessibles dans la variable `{{.State}}`, et les différents comptes d'états sont `.Info`, `.Minor`, `.Major`, et `.Critical`.  
Un exemple de sortie avec un compte des alarmes dans l'état `Minor` serait donc : `Alarmes mineures : {{.State.Minor}}`.  

### Exemples

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
    "output_template": "Alarmes critiques : {{.State.Critical}}"
}
```