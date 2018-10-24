# Service

La commande service  permet de démarrer les services canopsis.


```bash
service <service> start
service <service> stop
service <service> status
service <service> restart
```

les services disponible dans canopsis sont: 

```
mongodb
influxdb
rabbitmq-server
amqp2engines
webserver
```

## Attention

Dans le cas d'une installation Ansible les services mongodb, influxdb et rabbitmq-server utilisent les paquets du systéme, il faut donc utiliser les commandes des paquets hors de l'environnement canopsis.

```
service mongod (start, stop, restart, status)
service influxdb (start, stop, restart, status)
service rabbitmq-server (start, stop, restart, status)
```
