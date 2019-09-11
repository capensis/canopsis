# Connecteur Zabbix vers Canopsis (connector-zabbix2canopsis)

Convertit des alertes issues de triggers Zabbix en évènements Canopsis.

## Installation

### Installation des dépendances

Installation de SQLAlchemy pour interroger la base de données, et de Kombu pour communiquer avec RabbitMQ :

```sh
pip install sqlalchemy
pip install kombu
```

SQLAlchemy comprend différents dialectes de bases de données. Par exemple, si votre Zabbix est configuré pour utiliser MySQL, vous devez installer la dépendance MySQL pour Python :

```sh
pip install mysql
```

### Installation du connecteur

Installation de la version stable du connecteur :

```sh
pip install connector-zabbix2canopsis
```

Installation de la version de développement du connecteur :

```sh
pip install https://git.canopsis.net/canopsis-connectors/connector-zabbix2canopsis
```

## Configuration dans Zabbix

Créer un hostgroup `hg_canopsis`. Tous les hôtes à surveiller doivent appartenir à ce groupe.

`AlertScriptsPath` est configuré dans `zabbix_server.conf`. Créer une action `ac_send_canopsis`.


```
    Configure condition :
    Host group = hg_canopsis

    Configure operation :
    Operation type : remote command

    Execute on : zabbix server

    Command (replace AlertScriptsPath by its value) : "AlertScriptsPath"/send_zab_event2canop.py "{EVENT.DATE}" "{EVENT.TIME}" "{STATUS}" "{TRIGGER.NSEVERITY}" "{TRIGGER.ID}" "{TRIGGER.NAME}" "{HOST.NAME1}" "{ITEM.NAME1}" "{ITEM.VALUE1}" "{HOST.NAME2}" "{ITEM.NAME2}" "{ITEM.VALUE2}" "{HOST.NAME3}" "{ITEM.NAME3}" "{ITEM.VALUE3}" "{HOST.NAME4}" "{ITEM.NAME4}" "{ITEM.VALUE4}" "{HOST.NAME5}" "{ITEM.NAME5}" "{ITEM.VALUE5}" "{HOST.NAME6}"
```

Création d'un fichier de log :

```sh
touch /var/log/send_zab_event2canop.log && chown zabbix:zabbix /var/log/send_zab_event2canop.log
```

Créer un dossier tampon (**note :** remplacer `AlertScriptsPath` par sa valeur) :

```sh
mkdir -p "AlertScriptsPath/connector_buffer" && chown zabbix:zabbix "AlertScriptsPath/connector_buffer"
```

Le jeton doit être identique à celui du fichier `connector-zabbix2canopsis.config`.
