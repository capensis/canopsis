# Zabbix

**TODO (DWU) :** remettre au propre cette doc.

`connector-zabbix2canopsis` est basée sur le toolkit sqlalchemy (http://www.sqlalchemy.org/) pour interroger la base de donée et sur la librairie kombu pour publier des messages dans RabbitMQ.

```
    pip install sqlalchemy
    pip install kombu
```

SQLAlchemy comprend différents dialects (http://docs.sqlalchemy.org/en/latest/dialects/index.html) qui on besoin de dépendances.

Par exemple, si vous souhaitez interroger une base de données MySQL, vous devez installer python-mysql lib

```
    pip install mysql
```

Version stable :

```
    pip install connector-zabbix2canopsis
```

Development version:

```
    pip install https://git.canopsis.net/canopsis-connectors/connector-zabbix2canopsis
```

## Install Zabbix action

Create hostgroup hg_canopsis, all hosts monitored must belong to this hostgroup.  

AlertScriptsPath is configued in zabbix_server.conf.  
Create action ac_send_canopsis.  

Créer le groupe hôte hg_canopsis, tous les hôtes surveillés doivent appartenir à ce groupe.  

AlertScriptsPath est configuré dans zabbix_server.conf. Créez l'action ac_send_canopsis.  

```
    Configure condition :
    Host group = hg_canopsis
    Configure operation :
    Operation type : remote command
    Execute on : zabbix server
    Command(replace AlertScriptsPath by its value) : "AlertScriptsPath"/send_zab_event2canop.py "{EVENT.DATE}" "{EVENT.TIME}" "{STATUS}" "{TRIGGER.NSEVERITY}" "{TRIGGER.ID}" "{TRIGGER.NAME}" "{HOST.NAME1}" "{ITEM.NAME1}" "{ITEM.VALUE1}" "{HOST.NAME2}" "{ITEM.NAME2}" "{ITEM.VALUE2}" "{HOST.NAME3}" "{ITEM.NAME3}" "{ITEM.VALUE3}" "{HOST.NAME4}" "{ITEM.NAME4}" "{ITEM.VALUE4}" "{HOST.NAME5}" "{ITEM.NAME5}" "{ITEM.VALUE5}" "{HOST.NAME6}"
```

Création d'un log file :

```
    touch /var/log/send_zab_event2canop.log && chown zabbix:zabbix /var/log/send_zab_event2canop.log
```

Créer un dossier tampon (remplacez AlertScriptsPath par sa valeur):

```
    mkdir "AlertScriptsPath"/connector_buffer  && chown zabbix:zabbix  "AlertScriptsPath"/connector_buffer""
```

Le jeton doit être identique à celui du fichier connector-zabbix2canopsis.config.
