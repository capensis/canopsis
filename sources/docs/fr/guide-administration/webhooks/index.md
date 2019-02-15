# Webhooks

!!! note
    Cette fonctionnalité n'est disponible que dans l'édition CAT de Canopsis.

Les webhooks sont une fonctionnalité du moteur `axe` permettant d'automatiser la gestion de la vie des tickets vers un service externe en fonction de l'état des événements ou des alarmes.

Les webhooks sont définis dans la collection MongoDB `webhooks`, et
peuvent-être ajoutées et modifiées avec l'[API webhooks](../../guide-developpement/webhooks/api_v2_webhooks.md).

Des exemples pratiques d'utilisation des webhooks sont disponibles dans la partie [Exemples](#exemples).

## Définition d'un webhook

Une règle est un document JSON contenant les paramètres suivants :
 - `_id` (optionnel): l'identifiant du webhook (généré automatiquement ou choisi par l'utilisateur).
 - `hook` (requis) : les conditions dans lesquelles le webhook doit être appelé, dont :
     - `entity_patterns` (optionnel) : Liste de patterns permettant de filtrer les entités.
     - `event_patterns` (optionnel) : Liste de patterns permettant de filtrer les événements. Le format des patterns est le même que pour l'[event-filter](../event-filter/index.md).
     - `triggers` (requis) : Liste de triggers. Au moins un de ces triggers doit avoir eu lieu pour que le webhook soit appelé.
 - `request` (requis) : les informations nécessaires pour générer la requête vers le service externe, dont :
     - `auth` (optionnel) : les identifiants pour l'authentification HTTP
     - `headers` (optionnel) : les en-têtes de la requête
     - `method` (requis) : méthode HTTP
     - `payload` (requis) : le corps de la requête qui sera envoyé. Le payload peut être personnalisé grâce aux [Templates](#templates).
     - `url` (requis) : l'url du service externe. L'URL est personnalisable grâce aux [Templates](#templates).
 - `declare_ticket` (optionnel) : les champs qui seront extraits de la réponse du service externe. Si `declare_ticket` est défini alors les données seront récupérées et un step `declareticket` est ajouté à l'alarme.
     - `ticket_id` est le mom du champs de la réponse contenant le numéro du ticket créé dans le service externe. La réponse du service est supposé être un objet JSON.

Si la variable d'environnement `CPS_CERTIFICATES_DIRECTORY` est définie, et est un chemin vers un dossier existant, les fichiers de ce dossier sont ajoutés aux certificats SSL de confiance pour les requêtes effectuées par le service de webhooks.

### Activation d'un webhook

Le champ `hook` représente les conditions d'activation d'un webhook. Il contient obligatoirement `triggers` qui est un tableau de triggers et éventuellement des `patterns` sur les alarmes, les entités et les évenements.

Les triggers possibles sont : `"stateinc"`, `"statedec"`, `"create"`, `"ack"`, `"ackremove"`, `"cancel"`, `"uncancel"`, `"declareticket"`, `"assocticket"`, `"snooze"`, `"unsnooze"`, `"resolve"`, `"done"`, et `"comment"`.

`entity_patterns` est un tableau pouvant contenir plusieurs patterns d'entités. Si plusieurs patterns sont ainsi définies, il suffit qu'un seul pattern d'entités corresponde à l'alarme en cours pour que la condition sur les `entity_patterns` soit validée. Il en va de même pour `event_patterns` (tableaux de patterns d'événements).

Si des triggers et des patterns sont définies dans le même hook, le webhook est activé s'il correspond à la liste des triggers et en même temps aux différentes listes de patterns.

### Templates

Les champs `payload` et `url` sont personnalisables grâce aux templates. Les templates permettent de générer du texte en fonction de l'état de l'alarme, de l'évenement ou de l'entité.

`{{ .Alarm }}` permet d'accéder aux propriétés d'une alarme, de même que `{{ .Event }}` pour un événement et `{{ .Entity }}` pour une entité.

Ces trois éléments contiennent plusieurs qu'on peut utiliser pour créer des chaînes dynamiques. Par exemple, `"Component : {{ .Alarm.Value.Component }}` va créer la chaîne de caractères `"Component : comp"` si le nom du component est `comp`.

On peut également générer du texte en fonction de l'état de la variable. Dans le cas suivant `"{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}"`, on obtiendra `"http://127.0.0.1:5000/even"` si le statut de l'alarm vaut 0, 2 ou 4, `"http://127.0.0.1:5000/odd"` sinon.

De même façon que le `or` et le `eq`, il est possible de tester les conditions avec `and`, `not`, `ne` (not equal), `lt` (less than), `le` (less than or equal), `gt` (greater than) ou `ge` (greater than or equal).

Pour plus d'informations, vous pouvez consulter la [documentaion offficelle de Go sur les templates](https://golang.org/pkg/text/template).

### Données externes

Si `declare_ticket` est défini, les données récupérées du service externe sont stockées dans l'alarme et un step `declareticket` est ajouté à l'alarme.

Dans `declare_ticket`, le champ `ticket_id` définit le champ où se trouve l'identificant du ticket créé via le service externe. Par exemple, si l'API retourne l'id dans le champ `numberTicket`, il faudra ajouter dans `declare_ticket` la ligne `"ticket_id" : "numberTicket"`.

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

## Activation du plugin webhooks

Le moteur axe par défaut ne contient pas ce plugin gérant les webhooks. Pour pouvoir utiliser les webhooks, il faut :
- compiler le plugin webhooks dans le répertoire contenant le plugin webhooks `CGO_ENABLED=1 go build -buildmode=plugin -o webhookPlugin.so main.go`
- lancer le moteur axe avec l'option `-postProcessorsDirectory <dossier contenant webhookPlugin.so>`. Sauf configuration spécifique, `webhookPlugin.so` se trouve dans `/plugins/axepostprocessor`.

Le plugin webhook se trouve dans l'image `engine-axe-cat`.

!!! note
    Il est déconseillé d'utiliser l'option `-autoDeclareTickets` qui est **dépréciée** depuis la version 3.10.0 et l'ajout des webhooks.