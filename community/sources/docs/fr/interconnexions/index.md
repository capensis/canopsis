# Liste des interconnexions Canopsis

Les interconnexions sont les canaux par lesquels Canopsis communique avec d’autres applications. Elles peuvent être réparties en trois catégories :

1.  Les [connecteurs](#connecteurs) alimentent Canopsis en évènements transmis par des sources extérieures.
2.  Les [drivers](#drivers) amènent le référentiel qui permettra d’enrichir les évènements.
3.  Les [notifications](#notifications) sont émises par Canopsis vers différents outils à partir de jeux de règles et de déclencheurs.

En complément, Canopsis embarque des [API](#exploitation-par-les-api) que l'on peut utiliser pour l'exploitation.

## Connecteurs

Un connecteur permet d’envoyer à Canopsis des évènements à partir de sources d'informations extérieures.

!!! note
    Veuillez noter que les alarmes générées par vos connecteurs ont une limitation de 256 caractères, sur certains champs. Si cette limite est dépassée sur ces champs, les évènements générés par vos connecteurs ne pourront pas être traités par Canopsis.

    Consultez la [documentation des limitations des évènements Canopsis](../guide-utilisation/limitations/index.md#limitations-des-evenements) pour en savoir plus.

### Base de données

| **Nom** | **Source** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:------------:|
| [SQL](Base-de-donnees/Mysql-MariaDB-PostgreSQL-Oracle.md) | Mysql, PostgeSQL, Oracle, DB2 et MSSQL | Community | Oui | Toutes versions |
| [Elasticsearch](Base-de-donnees/Elasticsearch.md) | Elasticsearch | Community | Oui | Toutes versions |

### Transport

| **Nom** | **Source** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| [Logstash](Transport/Logstash.md) | [Liste des sources](https://www.elastic.co/guide/en/logstash/current/input-plugins.html) | Community | Oui | Toutes versions |
| [Email](Transport/Mail.md) | Messages provenant d’une boîte mail **POP3** | Pro | Oui | Toutes versions |
| [Send_event](Transport/send_event.md) | Script Python (version 2.x et 3.x) exécutable en environnement Linux ou Windows | Community | Oui | Toutes versions |

### Supervision

| **Nom** | **Source** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| [SNMP trap](Supervision/SNMPtrap.md) | Tout trap SNMP respectant la [RFC1157](https://www.rfc-editor.org/rfc/pdfrfc/rfc1157.txt.pdf) ou nécéssitant un traitement spécifique | Pro | Oui | Toutes versions |
| [Nagios](Supervision/Nagios-et-Icinga.md) | [Icinga](https://icinga.com/) 1, [Nagios](https://www.nagios.org/)  ≤ 3.x (Nagios 4.x en beta) | Community | Oui | Toutes versions |
| [Centreon Legacy](Supervision/Centreon.md) | [Centreon](https://www.centreon.com/) 2.11.5 à 2.11.7, 3.0.3 à 3.0.11, 3.0.13, 3.0.14, 3.0.16, 18.10 et 19.04 | Community | Oui | Toutes versions |
| [Centreon Stream Connector](Supervision/Centreon-stream-connector.md) | [Centreon Stream Connector](https://docs.centreon.com/current/en/developer/developer-broker-stream-connector.html) | Community | Oui | 19.10.5, >= 20.04.2 |
| [Icinga](Supervision/Icinga2.md) | [Icinga 2](https://icinga.com/) | Community | Oui | Toutes versions |
| [Zabbix](Supervision/Zabbix.md) | [Zabbix](https://www.zabbix.com/) | Community | Oui | Toutes versions |
| [LibreNMS](Supervision/LibreNMS.md) | [LibreNMS](https://www.librenms.org/) | Community | Oui | Toutes versions |
| [Shinken](Supervision/Shinken.md) | [Shinken](http://www.shinken-monitoring.org/) | Community | Oui | Toutes versions |
| Datametrie | [Datametrie](https://web.archive.org/web/20190804083856/https://www.ip-label.fr/produits/datametrie-global-experience/) | Pro | Non | Toutes versions |
| Canopsis | Toute données provenant d’un autre Canopsis | Pro | Oui | Toutes versions |
| [Nokia NSP](Supervision/NokiaNSP.md) | [Version 19.3](https://www.nokia.com/networks/products/network-services-platform/) | Pro | Non  Version ≥ 3  |
| [Prometheus](Supervision/Prometheus.md) | [Prometheus](https://prometheus.io/docs/introduction/overview/) | Community | Oui | Toutes versions |
| [PRTG](Supervision/PRTG.md) | [PRTG](https://www.paessler.com/manuals/prtg/introduction) | Community | Oui | >= 4.3 |

## Drivers

Le driver permet de peupler le référentiel interne Canopsis en vue de l’enrichissement des alarmes.

**NB :** Chaque driver dans ce tableau est à considérer comme un framework de synchronisation qui doit être adapté à chaque contexte client (modèle de données, champs à synchroniser, liens applicatifs…).

### Référentiel

| **Nom** | **Source** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| iTop | Version [Pro 1.3.4-3287](https://www.combodo.com/itop) et [Community 2.4](https://www.combodo.com/itop) | Pro | Non | Version ≥ 3.25 |
| Service Now | Version [Madrid](https://www.servicenow.fr/) | Pro | Non | Version ≥ 3.25 |
| Easyvista | [Easyvista](https://www.easyvista.com/fr) | Pro | Non | Version ≥ 3.25 |
| CSV | Sources de données CSV spécifiques | Pro | Non | Version ≥ 3.25 |
| [API](drivers/driver-api.md) | Source de données via API externe | Pro | Oui | Version ≥ 4.3 |

## Notifications

Canopsis permet d’émettre des notifications vers différents outils à partir d’un jeu de règles et de déclencheurs (créations d’alarmes, acquittements, changements de criticité…). Les possibilités de notifications offertes par Canopsis sont toutes dépendantes du modèle de données de l’outil cible (la création d’un ticket d’incident ServiceNow n’est pas forcément identique d’une instance à l’autre).

### Générique

| **Nom** | **Destination** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| [Webhooks](../guide-utilisation/menu-exploitation/scenarios/) | Tout outil / API HTTP qui peut réceptionner des appels webhooks | Pro | Oui | Version ≥ 3 (moteurs Go) |

### Transport

| **Nom** | **Destination** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| [Logstash](../guide-utilisation/cas-d-usage/notifications.md) | [Liste des destinations](https://www.elastic.co/guide/en/logstash/current/output-plugins.html) | Pro | Oui | Version ≥ 3 (moteurs Go) |
| [IM](../guide-utilisation/cas-d-usage/notifications.md) | Toute messagerie instantanée disposant d’une API qui accepte des requêtes HTTP | Pro | Oui | Version ≥ 3 (moteurs Go) |
| Email | Tout serveur email disposant d’une API qui accepte des requêtes HTTP | Pro | Oui | Version ≥ 3 (moteurs Go) |

### Ordonnanceur de Tâches

| **Nom** | **Version Validée(s)** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| [AWX](https://doc.canopsis.net/guide-administration/remediation/#configuration-awx) | 14.0.0 | Pro | Non | Version ≥ 4 (moteurs Go) |
| [Rundeck](https://doc.canopsis.net/guide-administration/remediation/#configuration-pour-rundeck) | 3.3.7 | Pro | Non | Version ≥ 4 (moteurs Go) |
| [Jenkins](https://doc.canopsis.net/guide-administration/remediation/#configuration-pour-jenkins) | 2.297 | Pro | Non | Version ≥ 4 (moteurs Go) |


### Ticketing

| **Nom** | **Destination** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:--------:|:---------:|:----------:|:-----------:|
| Service Now | Version [Madrid](https://www.servicenow.fr/) | Pro | Oui | Version ≥ 3 (moteurs Go) |
| Observer | Observer | Pro | Non | Version ≥ 3 (moteurs Go) |

## Exploitation par les API

| **Nom** | **Édition** | **Supporté** *(dans le cadre de mise à jour)* | **Compatibilité Canopsis** |
|:-----:|:---------:|:----------:|:-----------:|
| [Publication d’évènement](../guide-developpement/index.md#api) | Community | Oui | Version ≥ 3 |
| [Manipulation de comportements périodiques](../guide-utilisation/interface/pbehaviors/index.md) | Community | Oui | Version ≥ 3 |
