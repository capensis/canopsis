# Webhooks

!!! note
    Cette fonctionnalité n'est disponible que dans l'édition CAT de Canopsis.

Les webhooks sont une fonctionnalité du moteur `axe` permettant d'automatiser la gestion de la vie des tickets vers un service externe en fonction de l'état des évènements ou des alarmes.

Les webhooks sont définis dans la collection MongoDB `webhooks`, et peuvent être ajoutés et modifiés avec l'[API webhooks](../../guide-developpement/webhooks/api_v2_webhooks.md).

Des exemples pratiques d'utilisation des webhooks sont disponibles dans la partie [Exemples](#exemples).

## Activation du plugin webhooks

Les webhooks sont implémentés sous la forme d'un plugin à ajouter dans le moteur `axe`. Ce plugin n'est disponible qu'avec une installation CAT de Canopsis.

### Activation avec Docker

Dans une installation Docker, l'image `canopsis/engine-axe-cat` remplace l'image par défaut `canopsis/engine-axe`. Le moteur `axe` doit ensuite être lancé au minimum avec l'option suivante pour que le plugin des webhooks soit chargé : `engine-axe -postProcessorsDirectory /plugins/axepostprocessor`

### Activation par paquets

Pour pouvoir utiliser les webhooks avec une installation par paquets, il faut :
*  compiler le plugin webhooks dans le répertoire contenant le plugin webhooks `CGO_ENABLED=1 go build -buildmode=plugin -o webhookPlugin.so main.go`
*  lancer le moteur `axe` avec l'option `-postProcessorsDirectory <dossier contenant webhookPlugin.so>`. Sauf configuration spécifique, `webhookPlugin.so` se trouve dans `/plugins/axepostprocessor`.

## Définition d'un webhook

Une règle est un document JSON contenant les paramètres suivants :

 - `_id` (optionnel): l'identifiant du webhook (généré automatiquement ou choisi par l'utilisateur).
 - `hook` (requis) : les conditions dans lesquelles le webhook doit être appelé, dont :
     - `alarm_patterns` (optionnel) : Liste de patterns permettant de filtrer les alarmes.
     - `entity_patterns` (optionnel) : Liste de patterns permettant de filtrer les entités.
     - `event_patterns` (optionnel) : Liste de patterns permettant de filtrer les évènements. Le format des patterns est le même que pour l'[event-filter](../event-filter/index.md).
     - `triggers` (requis) : Liste de triggers. Au moins un de ces triggers doit avoir eu lieu pour que le webhook soit appelé.
 - `request` (requis) : les informations nécessaires pour générer la requête vers le service externe, dont :
     - `auth` (optionnel) : les identifiants pour l'authentification HTTP
       - `username` (optionnel) : nom d'utilisateur employé pour l'authentification HTTP
       - `password` (optionnel) : mot de passé employé pour l'authentification HTTP
     - `headers` (optionnel) : les en-têtes de la requête
     - `method` (requis) : méthode HTTP
     - `payload` (requis) : le corps de la requête qui sera envoyé. Il s'agit d'une chaîne de texte qui est parsée pour être transformée en fichier json. Les caractères spéciaux doivent être échappés. Le payload peut être personnalisé grâce aux [Templates](#templates).
     - `url` (requis) : l'url du service externe. L'URL est personnalisable grâce aux [Templates](#templates).
 - `declare_ticket` (optionnel) : les champs qui seront extraits de la réponse du service externe. Si `declare_ticket` est défini alors les données seront récupérées et un step `declareticket` est ajouté à l'alarme.
     - `ticket_id` est le mom du champs de la réponse contenant le numéro du ticket créé dans le service externe. La réponse du service est supposée être un objet JSON.

Si la variable d'environnement `CPS_CERTIFICATES_DIRECTORY` est définie, et qu'il s'agit d'un chemin vers un dossier existant, les fichiers de ce dossier sont ajoutés aux certificats SSL de confiance pour les requêtes effectuées par le service de webhooks.

Lors du lancement de moteur `axe`, plusieurs variables d'environnement sont utilisées (si elles existent) pour la configuration des webhooks :
- `SSL_CERT_FILE` indique un chemin vers un fichier de certificat SSL;
- `SSL_CERT_DIR` désigne un répertoire qui contient un ou plusieurs certificats SSL qui seront ajoutés aux certifcats de confiance;
- `HTTPS_PROXY` et `HTTP_PROXY` seront utilisés si la connexion au service externe nécessite un proxy.

### Activation d'un webhook

Le champ `hook` représente les conditions d'activation d'un webhook. Il contient obligatoirement `triggers` qui est un tableau de triggers et éventuellement des `patterns` sur les alarmes, les entités et les évènements.

Les triggers possibles sont : `"stateinc"`, `"statedec"`, `"create"`, `"ack"`, `"ackremove"`, `"cancel"`, `"uncancel"`, `"declareticket"`, `"assocticket"`, `"snooze"`, `"unsnooze"`, `"resolve"`, `"done"`, et `"comment"`.

| Nom                      | Description                                              |
|:-------------------------|:---------------------------------------------------------|
| `"ack"`                  | Acquittement d'une alerte                                |
| `"ackremove"`            | Suppression de l'acquittement                            |
| `"assocticket"`          | Asociation d'un ticket à l'alarme                        |
| `"cancel"`               | Annulation de l'évènement                                |
| `"create"`               | Création de l'évènement                                  |
| `"comment"`              | Envoi d'un commentaire                                   |
| `"declareticket"`        | Déclaration d'un ticket à l'alarme                       |
| `"done"`                 | Fin de l'alarme                                          |
| `"resolve"`              | Résolution de l'alarme                                   |
| `"snooze"`               | Report de l'alarme                                       |
| `"statedec"`             | Diminution de la criticité de l'alarme                   |
| `"stateinc"`             | Augmentation de la criticité de l'alarme                 |
| `"uncancel"`             | Retablissement de l'alarme                               |
| `"unsnooze"`             | Fin du report de l'alarme                                |

`entity_patterns` est un tableau pouvant contenir plusieurs patterns d'entités. Si plusieurs patterns sont ainsi définis, il suffit qu'un seul pattern d'entités corresponde à l'alarme en cours pour que la condition sur les `entity_patterns` soit validée. Il en va de même pour `alarm_patterns` (tableaux de patterns d'alarmes) et `event_patterns` (tableaux de patterns d'évènements).

Si des triggers et des patterns sont définies dans le même hook, le webhook est activé s'il correspond à la liste des triggers et en même temps aux différentes listes de patterns.

Par exemple, ce webhook va être activé si le trigger reçu par le moteur correspond à `"stateinc"` ou `"statedec"` ET que l'évènement ait comme `connector` soit `zabbix`, soit `shinken` ET que dans l'entité, l'`output` corresponde à l'expression régulière `MemoryDisk.*`.

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

Les templates des champs `payload` et `url` peuvent se décomposer en deux parties : la déclaration de variables et le corps du texte lui-même. La déclaration de variables doit être positionnée avant le corps du message. Les variables se distinguent du corps du message par le fait qu'elles sont entourés de doubles accolades.

Les variables stockent des informations sur les alarmes, les événements et les entités. `{{ .Alarm }}` permet d'accéder aux propriétés d'une alarme, de même que `{{ .Event }}` pour un évènement et `{{ .Entity }}` pour une entité. Ces trois éléments contiennent plusieurs propriétés qu'on peut utiliser pour créer des chaînes dynamiques. Par exemple, `"Component : {{ .Alarm.Value.Component }}` va créer la chaîne de caractères `"Component : comp"` si le nom du component est `comp`.

Dans l'exemple suivant, `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}`, on déclare d'abord deux variables, `$comp` et `$res`. Ensuite, on utilise ces deux variables pour générer l'adress URL que va appeler le moteur axe, par exemple `http://mon-api.xyz/nom-du-component/nom-de-la-ressource`.

On peut également générer du texte en fonction de l'état de la variable. Dans le cas suivant `"{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}"`, on obtiendra `"http://127.0.0.1:5000/even"` si le statut de l'alarm vaut 0, 2 ou 4, `"http://127.0.0.1:5000/odd"` sinon.

De même façon que le `or` et le `eq`, il est possible de tester les conditions avec `and`, `not`, `ne` (not equal), `lt` (less than), `le` (less than or equal), `gt` (greater than) ou `ge` (greater than or equal).

La fonction `js`, qui renvoie une chaîne de caractères échappée, peut être également mentionnée. Si, par exemple, la valeur dans `{{ .Event.Output }}` contient des caractères spéciaux comme des guillemets ou des backslashs, `{{ .Event.Output | js }}` permet d'échapper ces caractères.

Pour plus d'informations, vous pouvez consulter la [documentaion officielle de Go sur les templates](https://golang.org/pkg/text/template).

### Données externes

Si `declare_ticket` est défini, les données récupérées du service externe sont stockées dans l'alarme et un step `declareticket` est ajouté à l'alarme.

Dans `declare_ticket`, le champ `ticket_id` définit le champ où se trouve l'identifiant du ticket créé via le service externe. Par exemple, si l'API retourne l'id dans le champ `numberTicket`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "numberTicket"`.

Si l'API renvoie une réponse sous forme de json imbriqué, il faut prendre en compte le chemin d'accès. Par exemple, si l'API retourne une réponse de type ` {"result": {"number": "numéro d’incident"}}`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "result.number"`.

Les autres champs de `declare_ticket` sont stockés dans `Alarm.Value.Ticket.Data` de telle sorte que la clé dans `Data` corresponde à la valeur dans les données du service. Par exemple avec `"ticket_creation_date" : "timestamp"`, la valeur de `ticket["timestamp"]` sera mise dans `Alarm.Value.Ticket.Data["ticket_creation_date"]`.

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