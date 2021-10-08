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

!!! attention
    Un événement de type `check` entraîne la création ou mise à jour d'une alarme par le [moteur `engine-axe`](../../guide-administration/moteurs/moteur-axe.md#evenements-de-type-check). Il est important de s'assurer qu'elle puisse être également résolue par Logstash. À défaut, il est nécessaire que vous l'annuliez vous-même une fois l'alarme traitée.

## Exemple de configuration Logstash

#### Input

```bash
input {
    beats {
        port => "5044"
            tags => ["syslog"]
    }
}
```

#### Filter

```bash
filter {
    # ajout des champs essentiels pour un event Canopsis
    mutate {
        add_field => {
            "connector" => "logstash"
                "connector_name" => "logstash"
                "event_type" => "check"
                "state" => 1
        }
    }

    # Conversion des champs en integer
    mutate {
        convert => ["[state]", "integer"]
    }

    # S'assurer que l'on traite bien les event comportant le tag *syslog* défini dans l'input
    if "syslog" in [tags] {

        # Parse du log et récupération des informations nécessaires pour l'event et/ou la routing_key. Exemple 'component',
        # 'resource', 'output', 'timestamp', etc.
        grok {
            match => {"message" => "%{SYSLOGTIMESTAMP:timestamp} %{SYSLOGHOST:component} %{NOTSPACE:resource}\[%{NUMBER:pid}\]\: %{GREEDYDATA:output}"}
        }

        # Ajout du champ 'source_type' avec la valeur 'resource'
        mutate {
            # ce champ est ici a ajouter manuellement
            add_field => {"source_type" => "resource"}
        }
    }

    # Ex : exemple de traitement d'une valeur pour un besoin précis (Facultatif)
    mutate {
        update => { "resource" => "%{resource}-%{pid}" }
    }

    # construction de la clef de routage (routing key) necessaire à Canopsis
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

    # nettoyage : suppression des champs et métadonnées maintenant inutiles
    mutate {
        remove_field => ["[@metadata][canopsis_rk_]", "message"]
    }


    # Ajouter un timestamp à l'event (convertir une date en timestamp) :
    ruby {
        code => "event.set('timestamp', event.get('@timestamp').to_i)"
    }

}
```

#### Output

```bash
output {

    # Sortie debug (à supprimer)
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
 }
```

Un exemple de `routing_key` (rk) :

```
"routing_key" => "logstash.logstash.check.resource.cft-server.cft-logging"
```

## Dépannage

#### L'évènement n'apparaît pas dans Canopsis

Lorsque l'on active l'option `stdout { codec => rubydebug }`, les logs traités sont affichés sur la sortie standard, ce qui peut faciliter le debug.

On retrouve la structure de l'évènement avec les champs ainsi que leurs valeurs, mais aussi la *routing key* dans les méta-données.

1.  Vérifier que la *routing key* est correcte et au bon format (voir les exemples ci-dessus).
2.  Vérifier que l'évènement est au bon format (champs, valeurs…).
3.  Vérifier que l'évènement remonte bien dans RabbitMQ.
4.  Vérifier les logs des moteurs Canopsis.
