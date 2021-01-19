# Limitations de Canopsis

Le document suivant décrit les limitations connues de Canopsis à ce jour.

Les limitations se différencient des « bugs », dans le sens où elles sont :

* soit voulues ;
* soit inhérentes à l'architecture ou aux choix faits dans la conception de Canopsis ;
* soit non destinées à être corrigées dans l'immédiat (en raison d'une incidence mineure, d'une plus faible priorité, ou d'un coût de résolution trop important).

!!! note
    Ce document est encore en cours de rédaction, et peut être complété lors des différentes mises à jour de l'outil.

## Limitations de l'interface web

### Compatibilité des anciens navigateurs

Aucune prise en charge n'est prévue pour les navigateurs web ne disposant pas d'une version raisonnablement récente des moteurs de rendu Webkit (Google Chrome, Safari…) ou Gecko (Mozilla Firefox).

L'interface de Canopsis n'est donc **pas** compatible avec les navigateurs suivants :

* Internet Explorer, toutes versions ;
* Microsoft Edge, avant la version 80 ;
* Mozilla Firefox, avant la version 78 ESR ;
* Google Chrome, avant la version 85.

## Limitations des évènements 

### Encodage des évènements

Les évènements envoyés à Canopsis ne peuvent être encodés qu'en UTF-8 (ou en ASCII).

Les autres encodages, tels qu'ISO-8859-1, CP1252 ou UTF-16, ne sont pas gérés : les évènements en question peuvent être ignorés ou causer des problèmes d'affichage.

Vous devez donc vous assurer que vos appels à l'API et que vos connecteurs ne génèrent que de l'UTF-8 en sortie.

### Limite de 256 caractères sur l'unicité d'une alarme

(Ticket [#2554](https://git.canopsis.net/canopsis/canopsis/-/issues/2554) sur Gitlab)

!!! note
    Cette limite s'appliquant aux évènements Canopsis, elle s'applique aussi aux API, à l'interface web et aux connecteurs que vous branchez à Canopsis.

L'unicité d'une alarme est établie par la concaténation des champs `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`.

Cette chaîne sert actuellement de clé de routage (ou *routing key*) pour acheminer les évènements vers Canopsis, dans le cadre de notre utilisation de RabbitMQ. Le protocole AMQP qui lui est associé impose néanmoins une longueur maximale de 256 caractères à cette chaîne (cf. Section 4.9 de la [spécification AMQP](https://www.rabbitmq.com/resources/specs/amqp0-9-1.pdf)). Canopsis ne pouvant pas contourner cette limite du protocole, une exception `ShortStringTooLong` sera générée lorsque cette limite est dépassée.

Vous devez donc veiller à ce que l'ensemble `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]` ne dépasse jamais 256 caractères, sans quoi les évènements, traitements et alarmes associés ne pourront être traités par Canopsis.
