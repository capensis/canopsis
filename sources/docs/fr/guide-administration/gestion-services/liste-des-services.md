# Listing

| Nom du service | Catégorie du service  | Fichier de configuration du service         |
|----------------|-----------------------|---------------------------------------------|
| Redis          | Cache                 | `/opt/canopsis/etc/common/redis_store.conf` |
| MongoDB        | SGDB                  | `/opt/canopsis/etc/common/mongo_store.conf`, `/opt/canopsis/etc/cstorage.conf`, `/opt/canopsis/etc/mongo/storage.conf`|
| RabbitMQ       | Transfert de messages | `/opt/canopsis/etc/amqp.conf`               |
| InfluxDB       | Métriques             | `/opt/canopsis/etc/infux/storage.conf`      |
| Canopsis       | Hyperviseur           |                                             |

**TODO (DWU) :** `canoctl` livre un Nginx (ou HAproxy) par défaut. Est-ce qu'on l'utilise par défaut dans Core ?
