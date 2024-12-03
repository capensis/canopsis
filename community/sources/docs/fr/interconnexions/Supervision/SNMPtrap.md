# SNMP trap vers Canopsis

Réceptionne les traps SNMP, les traduit grâce à un jeu de règles et les convertit en évènements.

!!! info
    Ce connecteur n'est disponible que dans l'édition Pro de Canopsis.

## Introduction

Ce guide décrit la réception et la traduction de traps SNMP pour Canopsis.

Bon nombre d'équipements émettent des traps. Canopsis est en mesure de :

1. les réceptionner ;
2. les traduire grâce à un jeu de règles ;
3. les convertir en évènements.

Ce document aborde la première étape de cette mise en œuvre.

Pour comprendre le principe complet de traitement des traps, et 
notamment savoir ce qu'il adviendra des événements de type `trap` 
générés par le connecteur, consultez également la documentation des
[règles SNMP](../../guide-utilisation/menu-exploitation/regles-snmp.md).

## Émission des traps SNMP

L'émission des traps SNMP n'est pas traitée dans ce guide dans la mesure où cela concerne les équipements en eux-mêmes.

Il faut configurer sur les différents émetteurs l'adresse du récepteur de traps
ainsi que son port : il s'agit de l'adresse du connecteur `snmp2canopsis` et du
port UDP 162 (port par défaut d'un récepteur de traps SNMP).

## Connecteur snmp2canopsis

Le connecteur `snmp2canopsis` porte 3 missions :

1. Réceptionner les traps
2. Parser les traps et les transformer en JSON
3. Publier les messages JSON obtenus dans un exchange AMQP dédié sur Canopsis

### Configuration

#### Fichier de configuration

Le fichier de configuration du connecteur est `/etc/snmp2canopsis.conf`.

```ini
[snmp]
ip = 0.0.0.0
port = 162

[amqp]
host = rabbitmq
port = 5672
user = cpsrabbit
password = canopsis
vhost = canopsis
exchange = canopsis.snmp
```

#### Écoute SNMP

La section `[snmp]` contient la configuration pour l'IP et le port d'écoute des traps SNMP.

Pour permettre l'écoute quelle que soit l'adresse de la machine hôte, mettre la valeur `0.0.0.0`.

Par défaut le port d'écoute est le 162.

#### Connexion RabbitMQ

La section `[amqp]` contient la configuration pour la connexion au bus RabbitMQ.

Il faut donc vérifier que l'URL et les identifiants qui y figurent sont les bons.

La section `host` est à remplir avec l'IP ou le nom DNS du RabbitMQ.

### Déroulement

#### Décodage des traps

Une fois réceptionnés, les traps sont décodés puis transformés en JSON.

Exemple de JSON en sortie du connecteur :

```json
{
  "component": "127.0.0.1",
  "connector": "snmp",
  "connector_name": "snmp2canopsis",
  "event_type": "trap",
  "snmp_timeticks": "2350066",
  "snmp_trap_oid": "1.3.6.1.6.3.1.1.5.3",
  "snmp_vars": {
    "1.3.6.1.2.1.2.2.1.1": "1",
    "1.3.6.1.2.1.2.2.1.7": "2",
    "1.3.6.1.2.1.2.2.1.8": "2"
  },
  "snmp_version": "1",
  "source_type": "component",
  "state": 3,
  "state_type": 1,
  "timestamp": 1440075343.725282
}
```

Le connecteur ne possédant aucune MIB, `snmp_trap_oid` et le tableau `snmp_vars`
contiennent les OID des éléments sans aucune traduction.

On remarque que le connecteur produit une structure JSON qui rappelle le format
des évènements Canopsis. Notons le type d'évènement (`event_type`) : `trap`.

#### Publication du JSON

Les traps transformés en JSON sont publiés dans le bus AMQP de Canopsis, dans
un exchange dédié : `canopsis.snmp`.

#### Suite du traitement

Ces messages (évènements de type `trap`) seront traduits ultérieurement en
évènements de type `check`, par le moteur `snmp` de Canopsis.

### Exécution

Nous utilisons généralement un conteneur `Docker` pour cette étape.

```console
$ docker run -v snmp2canopsisdata:/connector-snmp2canopsis/etc \
    docker.canopsis.net/docker/pro/canopsis-pro-connector-snmp:v22.10.1

[2023-09-21 09:18:06.700607] INFO: snmp2canopsis: Read configuration from /connector-snmp2canopsis/etc/snmp2canopsis.conf
[2023-09-21 09:18:06.701409] DEBUG: amqp: Thread started
[2023-09-21 09:18:06.702131] INFO: amqp: Connecting to cpsrabbit@rabbitmq, on canopsis
[2023-09-21 09:18:06.701857] INFO: snmp: Start SNMP listener on 0.0.0.0:162
[2023-09-21 09:18:06.707382] DEBUG: amqp: Read the snmp queue
```

La configuration associée est la suivante :

```ini
[snmp]
ip = 0.0.0.0
port = 162

[amqp]
host = rabbitmq
port = 5672
user = cpsrabbit
password = canopsis
vhost = canopsis
exchange = canopsis.snmp
```

### Test de fonctionnement

Afin de valider le fonctionnement, nous pouvons générer un trap SNMP depuis une
machine.

Nous aurons besoin de la commande `snmptrap`.

Pour installer cette commande sous Debian :

```console
# apt install snmp
```

Pour installer cette commande sous EL (RHEL/Rocky Linux/AlmaLinux/…) :

```console
# yum install net-snmp-utils
```

Pour le test, nous allons nous appuyer sur la MIB Nagios
[NAGIOS-NOTIFY-MIB][notify_mib] et sa dépendance [nagios-root][root_mib].

[notify_mib]: https://github.com/monitoring-plugins/nagios-mib/raw/master/MIB/NAGIOS-NOTIFY-MIB
[root_mib]: https://github.com/nagios-plugins/nagios-mib/raw/master/src-mib/nagios-root.mib

Les deux fichiers récupérés doivent être placés dans le répertoire des MIB
SNMP : `/usr/share/snmp/mibs`.

Puisqu'il s'agit de traps SNMP, il faut s'intéresser au type `NOTIFICATION TYPE` présent dans les MIB.

Pour notre trap de test, nous allons utiliser l'objet `nSvcEvent` :

```
 nSvcEvent  NOTIFICATION-TYPE
   OBJECTS { nHostname, nHostStateID, nSvcDesc, nSvcStateID, nSvcAttempt,
             nSvcDurationSec, nSvcGroupName, nSvcLastCheck, nSvcLastChange,
             nSvcOutput }
   STATUS  current
   DESCRIPTION
     "The SNMP trap that is generated as a result of an event with the service
     in Nagios."
   ::= { nagiosNotify 7 }
```

Voici la ligne de commande utilisée avec `snmptrap` pour générer le trap :

```sh
snmptrap -v 2c -c public ${IP_RECEPTEUR} '' NAGIOS-NOTIFY-MIB::nSvcEvent \
  nHostname s "Equipement Impacte" \
  nSvcDesc s "Ressource Impactee" \
  nSvcStateID i 2 \
  nSvcOutput s "Message de sortie du trap SNMP"
```

Le connecteur reçoit ce trap, le convertit en JSON et le transmet à Canopsis
dans l'exchange `canopsis.snmp`.

Le JSON produit par le connecteur peut être vérifié dans les logs du service
ou conteneur `snmp2canopsis`. Exemple de log :

```
[2019-06-20 08:27:44.754789] INFO: snmp2canopsis: Read configuration from /etc/snmp2canopsis.conf
[2019-06-20 08:27:44.791038] INFO: snmp: Trap debug enabled
[2019-06-20 08:27:44.791343] DEBUG: amqp: Thread started
[2019-06-20 08:27:44.791476] INFO: amqp: Connecting to cpsrabbit@rabbitmq, on canopsis
[2019-06-20 08:27:44.919838] DEBUG: amqp: Read the snmp queue
[2019-06-20 08:27:44.920357] INFO: snmp: Start SNMP listener on 0.0.0.0:162
[2022-10-11 12:08:17.110603] DEBUG: snmp: {"snmp_version": "2c", "event_type": "trap", "timestamp": 1665490097.110236, "component": "172.20.0.1", "state_type": 1, "source_type": "component", "snmp_trap_oid": "1.3.6.1.4.1.20006.1.7", "snmp_vars": {"1.3.6.1.2.1.1.3.0": "278844995", "1.3.6.1.4.1.20006.1.3.1.17": "Message de sortie du trap SNMP", "1.3.6.1.6.3.1.1.4.1.0": "1.3.6.1.4.1.20006.1.7", "1.3.6.1.4.1.20006.1.1.1.2": "Equipement Impacte", "1.3.6.1.4.1.20006.1.3.1.6": "Ressource Impactee", "1.3.6.1.4.1.20006.1.3.1.7": "2"}, "connector": "snmp", "state": 3, "connector_name": "snmp2canopsis"}
```

## Suite

Afin d'être transformés en évènements Canopsis susceptibles de créer
ou mettre à jour des alarmes, les traps doivent encore faire l'objet
d'un traitement spécifique.

Deux cas sont possibles :

1. Les traps « standard » peuvent faire l'objet de règles de 
traitement grâce à des MIB ;
2. Les traps « non standard » ou traps pour lesquels vous ne 
possédez pas les MIB peuvent être traités avec du code personnalisé.

Se référer à la documentation sur les
[règles SNMP](../../guide-utilisation/menu-exploitation/regles-snmp.md
qui détaille ces parties de la mise en œuvre.
