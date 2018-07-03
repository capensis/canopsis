# InfluxDB

## Sommaire:

1. [Préambule](#préambule)
2. [Prérequis et dépendances](#prérequis-et-dépendances)
3. [Résumé](#résumé)
4. [Usage](#usage)
5. [Détails de fonctionnement](#détails-de-fonctionnement)
6. [Vérifications](#vérifications)
7. [Problèmes rencontrés](#problèmes-rencontrés)
8. [Informations complémentaires](#informations-complémentaires)

## Préambule

### Droits

Par Capensis, tous droits réservés.

### Fork

adresse     | commit
:-----------|-------:
https://github.com/bashrc666/ansible-role-influxdb  | [12626a9](https://github.com/bashrc666/ansible-role-influxdb/tree/12626a971a369db035ae786d9f6fb4d898a4f5cf)

### Auteurs

|Nom         | Prénom |
|:-----------|-------:|
|O'Tools     | Gregory|
|MARCHAND    | Paul   |

## Prérequis et dépendances

### OS Supportés

* CentOS7
* CentOS6


### Paquets

Paquets à installer pour faire fonctionner le rôle (ceux installés par le rôle et ceux requis sans toutes les dépendances).

| Paquet          | Version | Dépôt     |
| :-----          | :------ | :----     |
| influxdb.x86_64 | 1.5.1-1 | @influxdb |

*Dépendances :*

| Paquet             | Version     | Dépôt |
| :-----             | :------     | :---- |
| python2-pip.noarch | 8.1.2-5.el7 | @epel |

#### Ports

| Protocol | Port   | Service                  |
| :------- | :----- | :------                  |
| TCP      | 8086   | InfluxDB(HTTP)           |
| UDP      | 8089   | InfluxDB(UDP)            |
| TCP      | 8088   | InfluxDB(RPC for backup) |
## Résumé

Le rôle déploie InfluxDB, le configure et le lance sur le noeud de destination.

## Usage

### Mode d'utilisation

Dans un scénario simple, vous pouvez appeler le rôle de la façon suivante :
*playbook.yml*
```yaml
---
- hosts: all
  roles:
    - ansible-role-influxdb-maat
```

## Détails de fonctionnement

Le déploiement peut être modifié par le biais de variables, les voici :

`influxdb_version`<br/>
La version d'influxDB à déployer (par défaut: `1.2.2`)

`influxdb_admin_interface`<br/>
L'interface d'écoute de l'interface Administrateur (par défaut : `''`)

`influxdb_admin_port`<br/>
Le port découte de l'interface Web (par défaut : `8083`)

`influxdb_admin_use_ipv6`<br/>
Activez ou non le support de l'ipv6 (par défaut: `false`)

`influxdb_http_interface`<br/>
L'interface d'écoute de l'API (par défaut : `''`)

`influxdb_http_port`<br/>
Le port d'écoute de l'API (par défaut : `8086`)

`influxdb_http_use_ipv6`<br/>
Activez ou non le support de l'ipv6 (par défaut: `false`)

`influxdb_udp_enabled`<br/>
Active l'écoute en mode UDP (par défaut: `true`)

`influxdb_udp_bind_address`<br/>
L'interface d'écoute de l'API en UDP (par défaut: `''`)

`influxdb_udp_database_name`<br/>
Nom de la base de donnée Canopsis dans Influxdb (par défaut: `canopsis`)

### Création d'une base de données

Vous pouvez lancer la création d'une base de données en spécifiant les caractéristiques de cette dernière dans les variables présentés :
`influxdb_database_name`<br>
Le nom de la base de donnée (par défault: non défini)
`influxdb_username`<br>
Le nom d'utilisateur rattaché à la base de données (par défault: non défini)
`influxdb_password`<br>
Le mot de passe de l'utilisateur sus mentionné (par défault: non défini)


## Vérifications

Vérification du bon fonctionnement du déploiement du service :

> systemctl status influxdb

Vérifiez les ports en écoute pour le service :

> lsof -anPp $(pgrep influxd) -i4 -i6
```shell
COMMAND   PID     USER   FD   TYPE  DEVICE SIZE/OFF NODE NAME
influxd 31887 influxdb    3u  IPv6 8163148      0t0  TCP *:8088 (LISTEN)
influxd 31887 influxdb    6u  IPv6 8165637      0t0  TCP *:8083 (LISTEN)
influxd 31887 influxdb    9u  IPv6 8165638      0t0  TCP *:8086 (LISTEN)
```

**/!\ Attention**: vous devrez installer lsof pour cette vérification

## Informations complémentaires

### Accès externes

Dépôt Influx : https://repos.influxdata.com/

### SELinux

* Impactant : Non
* Géré par le rôle : Non

### Dossiers applicatifs

| Contenu | Path                 |
| :----   | :----                |
| Données | `/var/lib/influxdb/` |
