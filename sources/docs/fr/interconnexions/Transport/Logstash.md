# Logstash vers Canopsis

Transforme tout entrant Logstash en évènements Canopsis (ex : logs, requêtes HTTP…).

## Introduction

Cette documentation détaille la remontée des logs vers Canopsis via Logstash.

Connaître les base du fonctionnement de Logstash vous permettra d'avoir une meilleur compréhension de ce qui va suivre. Pour cela, nous vous invitons à vous documenter sur [l'input](https://www.elastic.co/guide/en/logstash/6.2/input-plugins.html), [les filtres](https://www.elastic.co/guide/en/logstash/6.2/filter-plugins.html) et [l'output](https://www.elastic.co/guide/en/logstash/6.2/output-plugins.html).

## Fonctionnement

La remontée des logs de logstash vers Canopsis se fait en passant par le bus AMQP.

Les logs devront donc être formatés d'une certaine manière afin que Canopsis puisse les comprendre et les intégrer.

### Informations AMQP

*  Vhost : `canopsis`
*  Routing key : `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`
*  Exchange : `canopsis.events`
*  Exchange Options : `type: topic`, `durable: true`, `auto_delete: false`
*  Content Type : `application/json`

!!! attention
    Faire très attention à la *routing key* !

### Structure d'un évènement

Actuellement, un évènement de type `log` doit être remonté dans Canopsis comme un évènement de type `check`.

Voici les champs attendus dans Canopsis afin que l'évènement puisse être reconnu :

```js
{
    // Les champs obligatoires

    'connector': 'logstash',
    'connector_name': 'logstash',
    'event_type': 'check',
    'source_type': 'resource',
    'component':   // Ex: Server5
    'resource':    // Ex: kernel
    'output':      // Message
    'state':       // Check state (0 - OK, 1 - WARNING, 2 - CRITICAL, 3 - UNKNOWN), default is 0
    'state_type':  // Check state type (0 - SOFT, 1 - HARD), default is 1
    'status':      // 0 == Ok | 1 == En cours | 2 == Furtif | 3 == Bagot | 4 == Annule

    // Les champs facultatifs

    'timestamp':   // UNIX timestamp for when the event was emitted (optional: set by the server to now)
}
```

!!! attention
    Les champs `state`, `state_type` et `status` doivent être de type `entier`.

!!! attention
    Le champ `timestamp`, s'il est utilisé, doit être au format timestamp Unix, sinon l'évènement risque de ne pas être interprété. S'il n'est pas renseigné, Canopsis utilisera la date et l'heure à laquelle il traite l'évènement.

## Exemple de configuration Logstash

#### Input

```
input {
  beats {
        port => "5044"
        tags => ["syslog"]
    }
}
```

#### Filter

```
filter {

    # ajout des champs essentiels pour un event Canopsis
        mutate {
          add_field => {
            "connector" => "logstash"
            "connector_name" => "logstash"
            "event_type" => "check"
            "state" => 1
            "state_type" => 1
            "status" => 1
          }
         }

    # Conversion des champs en integer
        mutate {
          convert => ["[state]", "integer"]
          convert => ["[state_type]", "integer"]
          convert => ["[status]", "integer"]
        }

   # S'assurer que l'on traite bien les event comportant le tag *syslog* défini dans l'input
        if "syslog" in [tags] {

            # Parse de la log et récupération des informations nécéssaires pour l'event et/ou la routing_key. Exemple 'component',
            # 'resource', 'output', 'timestamp'...etc.
            grok {
                match => {"message" => "%{SYSLOGTIMESTAMP:timestamp} %{SYSLOGHOST:component} %{NOTSPACE:resource}\[%{NUMBER:pid}\]\: %{GREEDYDATA:output}"}

            }

            #Ajout du champ 'source_type' avec la valeur 'resource'
            mutate {
                # ce champ est ici a ajouter manuellement
                add_field => {"source_type" => "resource"}
            }
        }

        Ex: exemple de traitement d'une valeur pour un besoin précis (Facultatif)
        mutate {
            update => { "resource" => "%{resource}-%{pid}" }
        }

        # construction de la clef de routage (routing key) necessaire a canopsis
        mutate {
            add_field => {"[@metadata][canopsis_rk_]" => "%{connector}.%{connector_name}.%{event_type}.%{source_type}.%{component}" }
        }

        if [source_type] == "resource" {
            mutate {
                add_field => {"[@metadata][canopsis_rk]" => "%{[@metadata][canopsis_rk_]}.%{resource}" }
            }
        } else {
            mutate {
                add_field => {"[@metadata][canopsis_rk]" => "%{[@metadata][canopsis_rk_]" }
            }
        }

        # nettoyage: suppression des champs et métadonnées maintenant inutiles
        mutate {
            remove_field => ["[@metadata][canopsis_rk_]", "message"]
        }


        # Ajouter un timestamp à l'event (convertir une date en timestamp) :

        ruby{
          code =>"event.set('timestamp', event.get('@timestamp').to_i)"
        }

}
```

#### Output

```
output {

  # Sortie debug (A supprimer)
  stdout { codec => rubydebug }

  rabbitmq {
        host => "XXX.XXX.XXX.XXX" # IP de la machine RabbitMQ (ou VIP)
        user => "my_user"
        password => "my_passwd"
        vhost => "canopsis"
        exchange => "canopsis.events"
        exchange_type => "topic"
        message_properties => { "content_type" => "application/json" }
        key => "%{[@metadata][canopsis_rk]}"
    }

}
```

!!! attention
    Faire bien attention à ce que la *routing key* soit au bon format (voir plus haut).

## Exemples

Un exemple d'évènement :

```
 message = {
  'connector': 'logstash',
  'connector_name': 'logstash',
  'event_type': 'check',
  'source_type': 'resource',
  'component': 'cft-server'
  'resource': 'cft-logging',
  'output': 'Server is going down',
  'state': 1,
  'state_type': 1,
  'status': 1
 }
```

Un exemple de `routing_key` (rk) :

```
"routing_key" => "logstash.logstash.check.resource.cft-server.cft-logging"
```

## Dépannage

#### L'évènement n'apparaît pas dans Canopsis

Lorsque l'on active l'option `stdout { codec => rubydebug }`, les logs traités sont affichés sur la sortie standard, ce qui peut faciliter le debug.

On retrouve la structure de l'évènement avec les champs ainsi que leurs valeurs, mais aussi la *routing key* au niveau des metadata.

1.  Vérifier que la *routing key* est correcte et au bon format (voir les exemples ci dessus).
2.  Vérifier que l'évènement est au bon format (champs, valeurs…).
3.  Vérifier que l'évènement remonte bien dans RabbitMQ dans la file `Engine_cleaner_events`
4.  Vérifier les logs des moteurs Canopsis. Vérifier qu'il n'y a pas d'erreurs dans `/opt/canopsis/var/log/engines/cleaner_events.log`
