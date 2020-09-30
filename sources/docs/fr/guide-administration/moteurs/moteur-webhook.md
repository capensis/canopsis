# Moteur `engine-webhook` (Go, CAT)

!!! info
    Disponible à partir de Canopsis 3.34.0, uniquement en édition CAT.

Le moteur `engine-webhook` permet d'automatiser la gestion de la vie des tickets vers un service externe en fonction de l'état des évènements ou des alarmes.

Jusqu'en `3.33.0`, les webhooks étaient une fonctionnalité implémentée sous la forme d'un plugin dans le moteur `engine-axe` (en version CAT).

Depuis la `3.34.0`, ils sont devenus leur propre moteur (disponible uniquement en version CAT).

Depuis la `3.37.0`, la fonction de répétition est disponible.

Depuis la `3.39.0`, le webhook peut être activé ou désactivé avec l'attribut `enabled`.

Depuis la `3.41.0`, les valeurs des champs `declare_ticket` peuvent être définies sous forme d'expressions régulières.

Les webhooks peuvent être ajoutés et modifiés via l'[API webhooks](../../guide-developpement/api/api-v2-webhooks.md).

Des exemples pratiques d'utilisation des webhooks sont disponibles dans la partie [Exemples](#exemples).

## Utilisation

### Options du moteur

La commande `engine-webhook -help` liste toutes les options acceptées par le moteur.

## Fonctionnement

À l'arrivée dans sa file, le moteur va vérifier si l'événement correspond à un ou plusieurs de ces Webhooks.

Si oui, il va alors immédiatement appliquer les Webhooks correspondant (il n'y a pas de *beat*).

En cas d'échec, il existe un mécanisme de réémission du webhook.

Vous pouvez trouver des cas d'usage pour la [notification via un outil tiers dans le guide d'utilisation](../../guide-utilisation/cas-d-usage/notifications.md).

### Définition d'un webhook

Une règle est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel) : l'identifiant du webhook (généré automatiquement ou choisi par l'utilisateur).
 - `enabled` (requis) : le webhook est-il activé ou non (booléen).
 - `hook` (requis) : les conditions dans lesquelles le webhook doit être appelé, dont :
     - `alarm_patterns` (optionnel) : Liste de patterns permettant de filtrer les alarmes.
     - `entity_patterns` (optionnel) : Liste de patterns permettant de filtrer les entités.
     - `event_patterns` (optionnel) : Liste de patterns permettant de filtrer les évènements. Le format des patterns est le même que pour l'[event-filter](moteur-che-event_filter.md).
     - [`triggers`](../architecture-interne/triggers.md) (requis) : Liste de [triggers](../architecture-interne/triggers.md). Au moins un de ces [triggers](../architecture-interne/triggers.md) doit avoir eu lieu pour que le webhook soit appelé.
 - `disable_if_active_pbehavior` (optionnel, `false` par défaut) : `true` pour désactiver le Webhook si un comportement périodique est actif sur l'entité.
 - `request` (requis) : les informations nécessaires pour générer la requête vers le service externe, dont :
     - `auth` (optionnel) : les identifiants pour l'authentification HTTP
       - `username` (optionnel) : identifiant utilisateur employé pour l'authentification HTTP
       - `password` (optionnel) : mot de passé employé pour l'authentification HTTP
     - `headers` (optionnel) : les en-têtes de la requête
     - `method` (requis) : méthode HTTP
     - `payload` (requis) : le corps de la requête qui sera envoyée. Il s'agit d'une chaîne de texte qui est parsée pour être transformée en fichier JSON. Les caractères spéciaux doivent être échappés. Le payload peut être personnalisé grâce aux [Templates](#templates).
     - `url` (requis) : l'URL du service externe. L'URL est personnalisable grâce aux [Templates](#templates).
 - `retry` (optionnel) : politique à suivre en cas d'échec
     - `count` (optionnel) : nombre de répétition
     - `delay` (optionnel) : intervalle entre 2 essais
     - `unit` (optionnel) : unité de temps de l'intervalle (notation : "s" pour seconde, "m" pour minute, "h" pour heure)
 - `declare_ticket` (optionnel) : les champs qui seront extraits de la réponse du service externe. Si `declare_ticket` est défini alors les données seront récupérées et un step `declareticket` est ajouté à l'alarme. Le [trigger `declareticketwebhook`](../architecture-interne/triggers.md) est également alors déclenché.
     - `ticket_id` est le nom du champ de la réponse contenant le numéro du ticket créé dans le service externe. La réponse du service est supposée être un objet JSON.
     - `regexp` est un booléen qui détermine si les valeurs des champs `ticket_id` ou tout autre champ de l'option `declare_ticket` doivent être traitées comme des expressions régulières.
     - `empty_response` est un champ qui précise si la réponse du service externe est vide ou non. Si ce champ est présent et qu'il vaut `true`, alors le webhook va s'activer en ignorant les autres champs du `declare_ticket`.

Lors du lancement du moteur `engine-webhook`, plusieurs variables d'environnement sont utilisées (si elles existent) pour la configuration des webhooks :

- `SSL_CERT_FILE` indique un chemin vers un fichier de certificat SSL ;
- `SSL_CERT_DIR` désigne un répertoire qui contient un ou plusieurs certificats SSL qui seront ajoutés aux certificats de confiance ;
- `NO_PROXY`, `HTTPS_PROXY` et `HTTP_PROXY` seront utilisés si la connexion au service externe nécessite un proxy.

!!! attention
    Les [`triggers`](../architecture-interne/triggers.md) `declareticketwebhook`, `resolve` et `unsnooze` n'étant pas déclenchés par des [évènements](../../guide-developpement/struct-event.md), ils ne sont pas utilisables avec les `event_patterns`.

!!! attention
    L'`event_pattern` sait uniquement parser les champs customs (autres que ceux définis dans les [`structures d'événements`](../../guide-developpement/struct-event.md#structure-basique-dun-evenement)) qui sont de type `string`. Il ne sait pas parser les champs customs de type `number` ou `boolean`.

### Activation d'un webhook

Le champ `hook` représente les conditions d'activation d'un webhook. Il contient obligatoirement [`triggers`](../architecture-interne/triggers.md) qui est un tableau de triggers et éventuellement des `patterns` sur les alarmes, les entités et les évènements.

Pour plus d'informations sur les `triggers` disponibles, consulter la [`documentation sur les triggers`](../architecture-interne/triggers.md)

`entity_patterns` est un tableau pouvant contenir plusieurs patterns d'entités. Si plusieurs patterns sont ainsi définis, il suffit qu'un seul pattern d'entités corresponde à l'alarme en cours pour que la condition sur les `entity_patterns` soit validée. Il en va de même pour `alarm_patterns` (tableaux de patterns d'alarmes) et `event_patterns` (tableaux de patterns d'évènements).

!!! attention
    L'activation d'un webhook ne doit pas être dépendante de la criticité d'une alarme appliquée par le moteur `engine-action`. En effet, dans l’enchaînement des moteurs, `engine-action` se situe après `engine-webhook`. Le webhook sera donc activé **avant** que le moteur `engine-action` ait pu changer la criticité de l'alarme.

Si des triggers et des patterns sont définies dans le même hook, le webhook est activé s'il correspond à la liste des triggers et en même temps aux différentes listes de patterns.

Par exemple, ce webhook va être activé si le trigger reçu par le moteur correspond à `stateinc` ou `statedec` ET si l'évènement a comme `connector` soit `zabbix`, soit `shinken` ET si dans l'entité, l'`output` correspond à l'expression régulière `MemoryDisk.*`.

```json
{
    "hook" : {

        "triggers" : ["stateinc", "statedec"],

        "event_patterns" : [
            {"connector" : "zabbix"},
            {"connector" : "shinken"}
        ],

        "entity_patterns" : [
            {"infos" :
                {"output" :
                    {
                        "value": {"regex_match": "MemoryDisk.*"}
                    }
                }
            }
        ],
    }
}
```

### Templates

Les champs `payload` et `url` sont personnalisables grâce aux templates. Les templates permettent de générer du texte en fonction de la criticité de l'alarme, de l'évènement ou de l'entité.

Pour plus d'informations, vous pouvez consulter la [documentation sur les templates Golang](../architecture-interne/templates-golang.md).

### Tentatives en cas d'échec

Lorsque le service appelé par le webhook répond une erreur (Code erreur HTTP != 200 ou timeout du service), plusieurs nouvelles tentatives sont effectuées avec un délai.  
`count` représente le nombre de nouvelles tentatives.  
`delay` et `unit` représentent le délai avant une nouvelle tentative.   

`unit` est exprimé en "s" pour seconde, "m" pour minute, et "h" pour heure.

Ces paramètres sont positionnés dans la configuration de chaque webhook.  
Les paramètres par défaut sont précisés dans un fichier de configuration (option `-configPath` de la ligne de commande).

Exemple de fichier `webhook.conf` :

```ini
count=5
delay=1
unit="m"
```

### Données externes

Si `declare_ticket` est défini, les données récupérées du service externe sont stockées dans l'alarme et un step `declareticket` est ajouté à l'alarme.

#### Récupération des données externes

Dans `declare_ticket`, le champ `ticket_id` définit le champ où se trouve l'identifiant du ticket créé via le service externe. Par exemple, si l'API retourne l'id dans le champ `numberTicket`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "numberTicket"`.

Si l'API renvoie une réponse sous forme de JSON imbriqué, il faut prendre en compte le chemin d'accès. Par exemple, si l'API retourne une réponse de type ` {"result": {"number": "numéro d’incident"}}`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "result.number"`.

Les autres champs de `declare_ticket` sont stockés dans `Alarm.Value.Ticket.Data` de telle sorte que la clé dans `Data` corresponde à la valeur dans les données du service. Par exemple avec `"ticket_creation_date" : "timestamp"`, la valeur de `ticket["timestamp"]` sera mise dans `Alarm.Value.Ticket.Data["ticket_creation_date"]`.

A partir de la version 3.41.0 de Canopsis, les valeurs des champs `ticket_id` et autres champs de `declare_ticket` peuvent être définies sous forme d'expressions régulières.

Pour cela, il est nécessaire de positionner l'option `regexp` à `true` comme dans l'exemple suivant :

```json
"declare_ticket": {
  "ticket_id": "objects\\.UserRequest.*\\.fields\\.id",
  "regexp": true
}
```

#### Réponse vide du service externe

Dans le cas où le service externe renvoie une réponse (par exemple `200 OK`) mais sans contenu, il est possible d'ajouter le champ `empty_response` au `declare_ticket`. Si le champ `empty_response` est présent et vaut `true` alors tous les autres champs du `declare_ticket` sont ignorés. Un step `declareticket` est créé avec un numéro de ticket qui vaut `"N/A"`.

Si le champ `empty_response` n'est pas présent dans le `declare_ticket` ou qu'il vaut `false`, alors le webhook se comporte comme d'habitude.

### Exemples

```json
{
    "_id" : "declare_external_ticket",
    "enabled" : true,
    "hook" : {
        "triggers" : [
            "create"
        ],
        "event_patterns" : [
            {
                "connector" : "zabbix"
            }
        ],
        "entity_patterns" : [
            {
                "infos" : {
                    "output" : {
                        "value": {"regex_match": "MemoryDisk.*"}
                    }
                }
            }
        ]
    },
    "disable_if_active_pbehavior": true,
    "request" : {
        "method" : "PUT",
        "auth" : {
            "username" : "ABC",
            "password" : "a!(b)-c_"
        },
        "url" : "{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}",
        "headers" : {
            "Content-type" : "application/json"
        },
        "payload" : "{{ $comp := .Alarm.Value.Component }}{{ $reso := .Alarm.Value.Resource }}{{ $val := .Alarm.Value.Status.Value }}{\"component\": \"{{$comp}}\",\"resource\": \"{{$reso}}\", \"parity\": {{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}},  \"value\": {{$val}} }"
    },
    "retry" : {
        "count": 5,
        "delay" : 1,
        "unit" : "m"
    },
    "declare_ticket" : {
        "ticket_id" : "id",
        "ticket_creation_date" : "timestamp",
        "priority" : "priority",
        "regexp": false
    }
}
```

## Collection MongoDB associée

Les webhooks sont stockés dans la collection MongoDB `webhooks`.

```json
{
    "_id" : "declare_external_ticket",
    "enabled" : true,
    "disable_if_active_pbehavior" : false,
    "request" : {
        "url" : "{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}",
        "headers" : {
            "Content-type" : "application/json"
        },
        "method" : "PUT",
        "auth" : {
            "username" : "ABC",
            "password" : "a!(b)-c_"
        },
        "payload" : "{{ $comp := .Alarm.Value.Component }}{{ $reso := .Alarm.Value.Resource }}{{ $val := .Alarm.Value.Status.Value }}{\"component\": \"{{$comp}}\",\"resource\": \"{{$reso}}\", \"parity\": {{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}},  \"value\": {{$val}} }"
    },
    "retry" : {
        "count": 5,
        "delay" : 1,
        "unit" : "m"
    },
    "declare_ticket" : {
        "priority" : "priority",
        "ticket_id" : "id",
        "ticket_creation_date" : "timestamp",
        "regexp": false
    },
    "hook" : {
        "entity_patterns" : [
            {
                "infos" : {
                    "output" : {
                        "value" : {
                            "regex_match" : "MemoryDisk.*"
                        }
                    }
                }
            }
        ],
        "event_patterns" : [
            {
                "connector" : "zabbix"
            }
        ],
        "triggers" : [
            "create"
        ]
    }
}
```
