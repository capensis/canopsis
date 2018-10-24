# API Get-alarms

Cette API permet de lister et filtrer les alarmes présentes dans Canopsis. 



## URL : 

`GET /alerts/get-alarms`

## paramètres


nom            | type   | Valeur par défaut | description 
---------------|--------|-------------------|-------------------------
tstart         | int    | null              | Timestamp UNIX de l'alarme la plus ancienne à inclure
tstop          | int    | null              | Timestamp UNIX de l'alarme la plus récente à include
opened         | bool   | true              | Si true, inclut les alarmes ouvertes. Si false, ne les inclut pas
resolved       | bool   | true              | Si true, inclut les alarmes résolues. Si false, ne les inclut pas
lookups        | list   | []                | inutilisé (déprécié)
filter         | dict   | {}                | filtre Mongo pour limiter les résultats de recherche
search         | string | ''                | Chaîne de caractère à rechercher. Peut être un mot ou une expression du langage de recherche
sort_key       | string | 'opened'          | attribut utilisé pour ordonner les résultats
sort_dir       | string | 'DESC'            | ordre de tri. Valeurs admises : `ASC` ou `DESC`
skip           | int    | 0                 | nombre de résultats à ignorer (utile pour la pagination)
limit          | int    | 50                | nombre de résultats retournés (utile pour la pagination)
with_step      | bool   | false             | Inclure les steps dans l'alarme
natural_search | bool   | false             | Est-ce une recherche naturelle ? 



## Valeurs de retour

La réponse est un flux json défini par 3 champs: 

nom     | type   | description 
--------|--------|-------------------------
total   | int    | nombre total d'éléments dans la réponse
data    | struct | liste d'objets de résultat (cf. section suivante)
success | bool   | `true` si la requête a réussi, false si une erreur est survenue.



l'objet data est une liste de dictionnaires, dont le contenu  varie en fonction de l'API utilisée. Pour `get-alarms`, il ne contient qu'un dict dont la clef `alarms` est une liste d'alarmes (cf. ci-dessous)

La structure d'une alarme est la suivante:

```
{
    "d": "bet_service1_pod_1",
    "entity": {
        "impact": [
            "watcher_microservice1",
            "Engine/Fluentd"
        ],
        "name": "Service 1 - Pod",
        "enable_history": [
            1511443520
        ],
        "measurements": [],
        "enabled": true,
        "depends": [
            "Engine/Fluentd/Fluentd"
        ],
        "infos": {
            "parent_service": {
                "value": "watcher_microservice1"
            }
        },
        "_id": "bet_service1_pod_1",
        "type": "component"
	},
	"t": 1511443520,                
    "v": {
    	    "status": {
    	        "a": "Engine/Fluentd.Fluentd",
    	        "_t": "statusinc",
    	        "m": "Panne",
    	        "t": 1511443520,
    	        "val": 1
    	    },
    	    "resolved": null,
    	    "resource": null,
    	    "tags": [
    	        "stats-opened"
    	    ],
    	    "ack": null,
    	    "extra": {},
    	    "component": "bet_service1_pod_1",
    	    "creation_date": 1511443520,
    	    "connector": "Engine/Fluentd",
    	    "canceled": null,
    	    "state": {
    	        "a": "Engine/Fluentd.Fluentd",
    	        "_t": "stateinc",
    	        "m": "Panne",
    	        "t": 1511443520,
    	        "val": 3
    	    },
    	    "connector_name": "Fluentd",
    	    "initial_output": "Panne",
    	    "last_update_date": 1511443520,
    	    "snooze": null,
    	    "ticket": null,
    	    "hard_limit": null
        },
        "infos": {
                        "parent_service": {
                            "value": "watcher_microservice1"
                        }
                    },
                    "_id": "c0f2a3f2-d051-11e7-b76a-00163e12248c"
                }

```


L'alarme est elle-même constituée de plusieurs sous-entités. 

La structure de l'alarme est la suivante, les entités qu'elle contient sont décrites ci-dessous: 


nom     | type       | description 
--------|------------|-------------------------
d 		| string     | identifiant de l'entité concernée par l'alarme
entity  | entity     | Entité concernée par l'alarme
t       | int        | Timestamp de l'alarme (repris depuis l'événement envoyé par la sonde de supervision)
v       | AlarmValue |  Contenu de l'alarme



La structure entity est la suivante : 

nom             | type         | description 
----------------|--------------|----------------------------------------------------------------------------------------------
impact          | list[string] | liste d'ID d'entités impactées par cette entité
name            | string       | Nom de l'entité
enable_history  | list[int]    | liste de timestamps auxquel l'entité a été activée
disable_history | list[int]    | liste de timestamps auxquel l'entité a été désactivée (peut être vide)
measurements    | list         | inutilisé
enabled         | bool         | true si l'entité est active, false sinon. Aucune alarme n'est créée sur une entité inactive
depends         | list[string] | Liste  d'ID d'entités impactant cette entité
infos           | object       | ensemble de pairs clef-valeur libre décrivant des informations métier liées à l'entité
_id             | string       | id unique de l'entité


La structure AlarmValue est la suivante :

nom              | type         | description 
-----------------|--------------|----------------------------------------------------------------------------------------------
status           | AlarmStep	| Statut de l'alarme (criticité)
resolved         | int 			| Timestamp Unix de la date de résolution de l'alarme (null si alarme en cours)
resource         | string       | identifiant de la resource
tags             | list[]string | Déprécié
ack              | AlarmStep    | Stockage d'un ACK s'il est posé
extra            | dict         | Déprécié
component        | string       | identifiant du composant
creation_date    | int          | Timestamp UNIX de la date de création de l'alarme dans Canopsis
connector        | string       | identifiant du connecteur
canceled         | AlarmStep    | 
state            | AlartStep    |
connector_name   | string       |
initial_output   | string       |
last_update_date | int          |
snooze           | AlarmStep    |
ticket           | AlarmStep    |
hard_limit       | bool         |
infos            | dict         |
_id              | string       |