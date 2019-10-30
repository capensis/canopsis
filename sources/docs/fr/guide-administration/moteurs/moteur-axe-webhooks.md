# Axe - Webhooks

!!! note
    Cette fonctionnalité n'est disponible que dans l'édition CAT de Canopsis.

Les webhooks sont une fonctionnalité du moteur `axe` permettant d'automatiser la gestion de la vie des tickets vers un service externe en fonction de l'état des évènements ou des alarmes.

Les webhooks sont définis dans la collection MongoDB `webhooks`, et peuvent être ajoutés et modifiés avec l'[API webhooks](../../guide-developpement/webhooks/api_v2_webhooks.md).

Des exemples pratiques d'utilisation des webhooks sont disponibles dans la partie [Exemples](#exemples).

## Activation du plugin webhooks

Les webhooks sont implémentés sous la forme d'un [plugin](../../guide-developpement/plugins/axe-post-processor.md) à ajouter dans le moteur `axe`. Ce plugin n'est disponible qu'avec une installation CAT de Canopsis.

### Activation avec Docker

Dans une installation Docker, l'image `canopsis/engine-axe-cat` remplace l'image par défaut `canopsis/engine-axe`. Le moteur `axe` doit ensuite être lancé au minimum avec l'option suivante pour que le plugin des webhooks soit chargé : `engine-axe -postProcessorsDirectory /plugins/axepostprocessor`

### Activation par paquets

Pour pouvoir utiliser les webhooks avec une installation par paquets, il faut :

*  compiler le plugin webhooks dans le répertoire contenant le plugin webhooks `CGO_ENABLED=1 go build -buildmode=plugin -o webhookPlugin.so main.go`
*  lancer le moteur `axe` avec l'option `-postProcessorsDirectory <dossier contenant webhookPlugin.so>`. Sauf configuration spécifique, `webhookPlugin.so` se trouve dans `/plugins/axepostprocessor`.

## Définition d'un webhook

Une règle est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel) : l'identifiant du webhook (généré automatiquement ou choisi par l'utilisateur).
 - `hook` (requis) : les conditions dans lesquelles le webhook doit être appelé, dont :
     - `alarm_patterns` (optionnel) : Liste de patterns permettant de filtrer les alarmes.
     - `entity_patterns` (optionnel) : Liste de patterns permettant de filtrer les entités.
     - `event_patterns` (optionnel) : Liste de patterns permettant de filtrer les évènements. Le format des patterns est le même que pour l'[event-filter](moteur-che-event_filter.md).
     - [`triggers`](../architecture-interne/triggers.md) (requis) : Liste de [triggers](../architecture-interne/triggers.md). Au moins un de ces [triggers](../architecture-interne/triggers.md) doit avoir eu lieu pour que le webhook soit appelé.
 - `disable_if_active_pbehavior` (optionnel, `false` par défaut) : `true` pour désactiver le webhook si un pbehavior est actif sur l'entité.
 - `request` (requis) : les informations nécessaires pour générer la requête vers le service externe, dont :
     - `auth` (optionnel) : les identifiants pour l'authentification HTTP
       - `username` (optionnel) : nom d'utilisateur employé pour l'authentification HTTP
       - `password` (optionnel) : mot de passé employé pour l'authentification HTTP
     - `headers` (optionnel) : les en-têtes de la requête
     - `method` (requis) : méthode HTTP
     - `payload` (requis) : le corps de la requête qui sera envoyée. Il s'agit d'une chaîne de texte qui est parsée pour être transformée en fichier JSON. Les caractères spéciaux doivent être échappés. Le payload peut être personnalisé grâce aux [Templates](#templates).
     - `url` (requis) : l'URL du service externe. L'URL est personnalisable grâce aux [Templates](#templates).
 - `declare_ticket` (optionnel) : les champs qui seront extraits de la réponse du service externe. Si `declare_ticket` est défini alors les données seront récupérées et un step `declareticket` est ajouté à l'alarme.
     - `ticket_id` est le nom du champ de la réponse contenant le numéro du ticket créé dans le service externe. La réponse du service est supposée être un objet JSON.
     - `empty_response` est un champ qui précise si la réponse du service externe est vide ou non. Si ce champ est présent et qu'il vaut `true`, alors le webhook va s'activer en ignorant les autres champs du `declare_ticket`.

Lors du lancement de moteur `axe`, plusieurs variables d'environnement sont utilisées (si elles existent) pour la configuration des webhooks :

- `SSL_CERT_FILE` indique un chemin vers un fichier de certificat SSL ;
- `SSL_CERT_DIR` désigne un répertoire qui contient un ou plusieurs certificats SSL qui seront ajoutés aux certificats de confiance ;
- `HTTPS_PROXY` et `HTTP_PROXY` seront utilisés si la connexion au service externe nécessite un proxy.

### Activation d'un webhook

Le champ `hook` représente les conditions d'activation d'un webhook. Il contient obligatoirement [`triggers`](../architecture-interne/triggers.md) qui est un tableau de triggers et éventuellement des `patterns` sur les alarmes, les entités et les évènements.

Pour plus d'informations sur les `triggers` disponibles, consulter la [`documentation sur les triggers`](../architecture-interne/triggers.md)

`entity_patterns` est un tableau pouvant contenir plusieurs patterns d'entités. Si plusieurs patterns sont ainsi définis, il suffit qu'un seul pattern d'entités corresponde à l'alarme en cours pour que la condition sur les `entity_patterns` soit validée. Il en va de même pour `alarm_patterns` (tableaux de patterns d'alarmes) et `event_patterns` (tableaux de patterns d'évènements).

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

Les champs `payload` et `url` sont personnalisables grâce aux templates. Les templates permettent de générer du texte en fonction de l'état de l'alarme, de l'évènement ou de l'entité.

Pour plus d'informations, vous pouvez consulter la [documentation sur les templates Golang](../architecture-interne/templates-golang.md).

### Données externes

Si `declare_ticket` est défini, les données récupérées du service externe sont stockées dans l'alarme et un step `declareticket` est ajouté à l'alarme.

#### Récupération des données externes

Dans `declare_ticket`, le champ `ticket_id` définit le champ où se trouve l'identifiant du ticket créé via le service externe. Par exemple, si l'API retourne l'id dans le champ `numberTicket`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "numberTicket"`.

Si l'API renvoie une réponse sous forme de JSON imbriqué, il faut prendre en compte le chemin d'accès. Par exemple, si l'API retourne une réponse de type ` {"result": {"number": "numéro d’incident"}}`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "result.number"`.

Les autres champs de `declare_ticket` sont stockés dans `Alarm.Value.Ticket.Data` de telle sorte que la clé dans `Data` corresponde à la valeur dans les données du service. Par exemple avec `"ticket_creation_date" : "timestamp"`, la valeur de `ticket["timestamp"]` sera mise dans `Alarm.Value.Ticket.Data["ticket_creation_date"]`.

#### Réponse vide du service externe

Dans le cas où le service externe renvoie une réponse (par exemple `200 OK`) mais sans contenu, il est possible d'ajouter le champ `empty_response` au `declare_ticket`. Si le champ `empty_response` est présent et vaut `true` alors tous les autres champs du `declare_ticket` sont ignorés. Un step `declareticket` est créé avec un numéro de ticket qui vaut `"N/A"`.

Si le champ `empty_response` n'est pas présent dans le `declare_ticket` ou qu'il vaut `false`, alors le webhook se comporte comme d'habitude.

### Exemples

```json
{
    "_id" : "declare_external_ticket",
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
    "declare_ticket" : {
        "ticket_id" : "id",
        "ticket_creation_date" : "timestamp",
        "priority" : "priority"
    }
}
```
