# Listing

|Nom du service|catégorie du service |fichier de configuration du service|
|--------------|---------------------|-----------------------------------|
|Redis         |Cache                | /opt/canopsis/etc/common/redis_store.conf     |
|MongoDB       |SGDB                 | /opt/canopsis/etc/common/mongo_store.conf|
|RabbitMQ      |Transfert de messages| /opt/canopsis/etc/amqp.conf|
|InfluxDB      |Métriques            | /opt/canopsis/etc/infux/storage.conf|
|Canopsis      |Hyperviseur          | |

**TODO (DWU) :** Pas sûr du contenu :

*  Redis : uniquement une variable d'environnement ?
*  MongoDB : OK.
*  RabbitMQ : OK.
*  InfluxDB : valider si c'est bien `/opt/canopsis/etc/infux/storage.conf`.
*  

**TODO (DWU) :** `canoctl` livre un Nginx (ou HAproxy) par défaut. Est-ce qu'on l'utilise par défaut dans Core ?
