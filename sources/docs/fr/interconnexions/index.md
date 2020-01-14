# Liste des interconnexions Canopsis

Les interconnexions sont les canaux par lesquels Canopsis communique avec d’autres applications. Elles peuvent être réparties en trois catégories :

1.  Les [connecteurs](#connecteurs) alimentent Canopsis en évènements transmis par des sources extérieures.
2.  Les [drivers](#drivers) amènent le référentiel qui permettra d’enrichir les évènements.
3.  Les [notifications](#notifications) sont émises par Canopsis vers différents outils à partir de jeux de règles et de déclencheurs.

En complément, Canopsis embarque des [API](#exploitation-par-les-api) que l'on peut utiliser pour l'exploitation.

## Connecteurs

Un connecteur permet d’envoyer à Canopsis des évènements à partir de sources d'informations extérieures.

### Base de données

| **Nom** | **Source(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| [SQL](Base-de-donnees/Mysql-MariaDB-PostgreSQL-Oracle.md) | Mysql, PostgeSQL, Oracle, DB2 et MSSQL | OpenCore | Oui | Outil de marché | Toutes versions |

### Transport

| **Nom** | **Source(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| [Logstash](Transport/Logstash.md) | [Liste des sources](https://www.elastic.co/guide/en/logstash/current/input-plugins.html) | OpenCore | Oui | Outil de marché | Toutes versions |
| [Email](Transport/Mail.md) | Messages provenant d’une boîte mail **POP3** | CAT | Oui | Outil de marché | Toutes versions |
| [Send_event](Transport/send_event.md) | Script Python (version 2.x et 3.x) exécutable en environnement Linux ou Windows | OpenCore | Oui | Outil de marché | Toutes versions |

### Supervision

| **Nom** | **Source(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| [Snmp Trap](Supervision/SNMPtrap.md) | Tout trap SNMP respectant la [RFC1157](https://www.rfc-editor.org/rfc/pdfrfc/rfc1157.txt.pdf) | CAT | Oui | Outil de marché | Toutes versions |
| [Snmp Trap Custom](Supervision/SNMPtrap_custom.md) | Tout trap SNMP nécessitant un traitement spécifique | CAT | Oui | Outil de marché | Version ≥ 3 |
| [Nagios](Supervision/Nagios-et-Icinga.md) | [Icinga](https://icinga.com/) 1, [Nagios](https://www.nagios.org/)  ≤ 3.x (Nagios 4.x en beta) | OpenCore | Oui | Outil de marché | Toutes versions |
| [Centreon](Supervision/Centreon.md) | [Centreon](https://www.centreon.com/) 2.11.5 à 2.11.7, 3.0.3 à 3.0.11, 3.0.13, 3.0.14, 3.0.16, 18.10 et 19.04 | OpenCore | Oui | Outil de marché | Toutes versions |
| [Icinga](Supervision/Icinga2.md) | [Icinga 2](https://icinga.com/) | OpenCore | Oui | Outil de marché | Toutes versions |
| [Zabbix](Supervision/Zabbix.md) | [Zabbix](https://www.zabbix.com/) | OpenCore | Oui | Outil de marché | Toutes versions |
| [LibreNMS](Supervision/LibreNMS.md) | [LibreNMS](https://www.librenms.org/) | OpenCore | Oui | Outil de marché | Toutes versions |
| [Shinken](Supervision/Shinken.md) | [Shinken](http://www.shinken-monitoring.org/) | OpenCore | Oui | Outil de marché | Toutes versions |
| [Datametrie](Supervision/datametrie.md) | [Datametrie](https://www.ip-label.fr/produits/datametrie-global-experience/) | CAT | Non | Outil client | Toutes versions |
| Canopsis | Toute données provenant d’un autre Canopsis | CAT | Oui | Outil de marché | Toutes versions |
| [Nokia NSP](Supervision/NokiaNSP.md) | [Version 19.3](https://www.nokia.com/networks/products/network-services-platform/) | CAT | Non | Outil client | Version ≥ 3  |
| [Cisco IoT FND](Supervision/Ciscoiotfnd.md) | [Version North Bound](https://www.cisco.com/c/en/us/products/cloud-systems-management/iot-field-network-director/index.html) API Release 3.0 | CAT | Non | Outil client | Version ≥ 3 |

## Drivers

Le driver permet de peupler le référentiel interne Canopsis en vue de l’enrichissement des évènements.

**NB :** Chaque driver dans ce tableau est à considérer comme un framework de synchronisation qui doit être adapté à chaque contexte client (modèle de données, champs à synchroniser, liens applicatifs…).

### Référentiel

| **Nom** | **Source(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| iTop | Version [Pro 1.3.4-3287](https://www.combodo.com/itop) et [Community 2.4](https://www.combodo.com/itop) | CAT | Non | Outil client | Version ≥ 3.25 |
| Service Now | Version [Madrid](https://www.servicenow.fr/) | CAT | Non | Outil client | Version ≥ 3.25 |
| Easyvista | [Easyvista](https://www.easyvista.com/fr) | CAT | Non | Outil client | Version ≥ 3.25 |
| CSV | Sources de données CSV spécifiques | CAT | Non | Outil client | Version ≥ 3.25 |

## Notifications

Canopsis permet d’émettre des notifications vers différents outils à partir d’un jeu de règles et de déclencheurs (créations d’alarmes, ACK, changements d’état…). Les possibilités de notifications offertes par Canopsis sont toutes dépendantes du modèle de données de l’outil cible (la création d’un ticket d’incident Snow n’est pas forcément identique d’une instance à l’autre).

### Générique

| **Nom** | **Destination(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| [Webhooks](../guide-administration/moteurs/moteur-webhook.md) | Tout outil qui peut réceptionner des webhooks | CAT | Oui | Outil de marché | Version ≥ 3 (moteurs Go) |

### Transport

| **Nom** | **Destination(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| [Logstash](../guide-utilisation/cas-d-usage/notifications.md) | [Liste des destinations](https://www.elastic.co/guide/en/logstash/current/output-plugins.html) | CAT | Oui | Outil de marché | Version ≥ 3 (moteurs Go) |
| [IM](../guide-utilisation/cas-d-usage/notifications.md) | Toute messagerie instantanée disposant d’une API qui accepte des requêtes HTTP | CAT | Oui | Outil de marché | Version ≥ 3 (moteurs Go) |
| Email | Tout serveur email disposant d’une API qui accepte des requêtes HTTP | CAT | Oui | Outil de marché | Version ≥ 3 (moteurs Go) |

### Ticketing

| **Nom** | **Destination(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|:------------:|
| Service Now | Version [Madrid](https://www.servicenow.fr/) | CAT | Oui | Outil de marché | Version ≥ 3 (moteurs Go) |
| Zelos | Zelos | CAT | Non | Outil client | Version ≥ 2 (moteurs Python) |
| Observer | Observer | CAT | Non |  Outil client| Version ≥ 3 (moteurs Go) |

## Exploitation par les API

| **Nom** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Type** | **Compatibilité Canopsis** |
|:-----:|:---------:|:----------:|:-----------:|:------------:|
| [Publication d’évènement](../guide-developpement/index.md#api) | OpenCore | Oui | Outil de marché | Version ≥ 3 |
| [Manipulation de Plage de Maintenance](../guide-utilisation/interface/pbehaviors/index.md) | OpenCore | Oui | Outil de marché | Version ≥ 3 |
| [Récupération de l’état de Santé de Canopsis](../guide-developpement/api/api-v2-healthcheck.md) | OpenCore | Oui |  Outil de marché | Version ≥ 3.4 |
